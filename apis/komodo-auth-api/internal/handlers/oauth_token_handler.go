package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"komodo-auth-api/internal/registry"
	"komodo-forge-sdk-go/crypto/jwt"
	"komodo-forge-sdk-go/crypto/oauth"
	httpErr "komodo-forge-sdk-go/http/errors"
	logger "komodo-forge-sdk-go/logging/runtime"
)

type TokenRequest struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	GrantType    string `json:"grantType"`
	Scope        string `json:"scope,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"` // For refresh_token grant
	Code         string `json:"code,omitempty"`          // For authorization_code grant
	RedirectURI  string `json:"redirectUri,omitempty"`  // For authorization_code grant
	Username     string `json:"username,omitempty"`      // For password grant
	Password     string `json:"password,omitempty"`      // For password grant
}

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	TokenType    string `json:"tokenType"`
	ExpiresIn    int    `json:"expiresIn"`
	RefreshToken string `json:"refreshToken,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

// Unified OAuth 2.0 token endpoint (RFC 6749 Section 3.2).
func OAuthTokenHandler(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")
	wtr.Header().Set("Cache-Control", "no-store")

	var reqBody TokenRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		logger.Error("failed to parse request body", err)
		httpErr.SendError(
			wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail("failed to parse request body"),
		)
		return
	}

	if reqBody.GrantType == "" {
		logger.Error("missing grant type", fmt.Errorf("missing grant type"))
		httpErr.SendError(
			wtr, req, httpErr.Auth.UnsupportedGrantType, httpErr.WithDetail("missing grant type"),
		)
		return
	}
	if !oauth.IsValidGrantType(reqBody.GrantType) {
		logger.Error("unsupported grant type: " + reqBody.GrantType, fmt.Errorf("unsupported grant type"))
		httpErr.SendError(
			wtr, req, httpErr.Auth.UnsupportedGrantType, httpErr.WithDetail("unsupported grant type"),
		)
		return
	}

	// Route to appropriate grant handler
	switch reqBody.GrantType {
		case "client_credentials":
			handleClientCredentials(wtr, req, &reqBody)
		case "refresh_token":
			handleRefreshToken(wtr, req, &reqBody)
		case "authorization_code":
			handleAuthorizationCode(wtr, req, &reqBody)
		default:
			logger.Error("unsupported grant type: " + reqBody.GrantType, fmt.Errorf("unsupported grant type"))
			httpErr.SendError(
				wtr, req, httpErr.Auth.UnsupportedGrantType, httpErr.WithDetail("unsupported grant type"),
			)
	}
}

// Handles M2M service authentication (RFC 6749 Section 4.4)
func handleClientCredentials(wtr http.ResponseWriter, req *http.Request, reqBody *TokenRequest) {
	// Validate client credentials
	if reqBody.ClientID == "" || reqBody.ClientSecret == "" {
		logger.Error("missing client credentials", fmt.Errorf("missing client credentials"))
		httpErr.SendError(
			wtr, req, httpErr.Auth.InvalidClientCredentials, httpErr.WithDetail("missing client credentials"),
		)
		return
	}

	if !registry.Validate(reqBody.ClientID, reqBody.ClientSecret) {
		logger.Error("invalid client credentials for: "+reqBody.ClientID, fmt.Errorf("invalid client credentials"))
		httpErr.SendError(
			wtr, req, httpErr.Auth.InvalidClientCredentials, httpErr.WithDetail("invalid client credentials"),
		)
		return
	}

	if reqBody.Scope != "" {
		if !oauth.IsValidScope(reqBody.Scope) {
			logger.Error("invalid scope: "+reqBody.Scope, fmt.Errorf("invalid scope"))
			httpErr.SendError(wtr, req, httpErr.Auth.InvalidScope, httpErr.WithDetail("invalid scope: "+reqBody.Scope))
			return
		}
		rec, _ := registry.Get(reqBody.ClientID)
		if !rec.HasScope(reqBody.Scope) {
			logger.Error("scope not permitted for client: "+reqBody.ClientID, fmt.Errorf("scope not permitted"))
			httpErr.SendError(wtr, req, httpErr.Auth.InsufficientScope, httpErr.WithDetail("client not permitted to request scope: "+reqBody.Scope))
			return
		}
	}

	// Issue access token (JWT) - no refresh token for client_credentials
	accessExpiresIn := int64(3600) // 1 hour

	// Parse scopes from space-separated string
	var scopes []string
	if reqBody.Scope != "" {
		scopes = strings.Fields(reqBody.Scope)
	}

	accessToken, err := jwt.SignToken(
		"komodo-auth-api",
		reqBody.ClientID,
		"komodo-apis:service",
		accessExpiresIn,
		scopes,
	)

	if err != nil {
		logger.Error("failed to sign access token", err)
		httpErr.SendError(wtr, req, httpErr.Global.Internal, httpErr.WithDetail("failed to sign access token"))
		return
	}

	// TODO: Store token JTI in Elasticache for tracking/revocation

	wtr.WriteHeader(http.StatusOK)
	json.NewEncoder(wtr).Encode(TokenResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int(accessExpiresIn),
		Scope:       reqBody.Scope,
	})

	logger.Info("issued client_credentials token for: " + reqBody.ClientID)
}

// Handles token refresh (RFC 6749 Section 6)
func handleRefreshToken(wtr http.ResponseWriter, req *http.Request, reqBody *TokenRequest) {
	if reqBody.RefreshToken == "" {
		httpErr.SendError(wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail("missing refresh token"))
		return
	}

	// Parse claims from refresh token
	claims, err := jwt.ParseClaims(reqBody.RefreshToken)
	if err != nil {
		logger.Error("failed to parse refresh token", err)
		httpErr.SendError(wtr, req, httpErr.Auth.InvalidToken, httpErr.WithDetail("failed to parse refresh token"))
		return
	}

	// Check if token is expired
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		logger.Error("refresh token is expired", fmt.Errorf("refresh token is expired"))
		httpErr.SendError(wtr, req, httpErr.Auth.InvalidToken, httpErr.WithDetail("refresh token is expired"))
		return
	}

	// TODO: Check if refresh token is revoked in Elasticache

	// Extract subject and scopes from claims
	clientID := claims.Subject
	scope := ""
	if len(claims.Scopes) > 0 {
		scope = strings.Join(claims.Scopes, " ")
	}

	// Issue new access token
	accessExpiresIn := int64(3600) // 1 hour

	accessToken, err := jwt.SignToken(
		"komodo-auth-api",
		clientID,
		"komodo-apis:user",
		accessExpiresIn,
		claims.Scopes,
	)
	if err != nil {
		logger.Error("failed to sign access token", err)
		httpErr.SendError(
			wtr, req, httpErr.Global.Internal, httpErr.WithDetail("failed to sign access token"),
		)
		return
	}

	wtr.WriteHeader(http.StatusOK)
	json.NewEncoder(wtr).Encode(TokenResponse{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(accessExpiresIn),
		Scope:        scope,
		RefreshToken: reqBody.RefreshToken, // Can optionally rotate
	})

	logger.Info("refreshed token for: " + clientID)
}

// Handles authorization code exchange (RFC 6749 Section 4.1)
func handleAuthorizationCode(wtr http.ResponseWriter, req *http.Request, reqBody *TokenRequest) {
	// TODO: Implement authorization code flow
	// 1. Validate code against stored authorization grants
	// 2. Verify redirect_uri matches original request
	// 3. Verify client credentials
	// 4. Issue access + refresh tokens
	// 5. Delete used authorization code

	httpErr.SendError(
		wtr, req, httpErr.Global.NotImplemented, httpErr.WithDetail("authorizationcode grant not yet implemented"),
	)
}
