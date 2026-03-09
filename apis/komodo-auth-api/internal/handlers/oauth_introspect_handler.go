package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"komodo-forge-sdk-go/crypto/jwt"
	logger "komodo-forge-sdk-go/logging/runtime"
)

type IntrospectResponse struct {
	Active    bool   `json:"active"`
	Scope     string `json:"scope,omitempty"`
	ClientID  string `json:"clientId,omitempty"`
	TokenType string `json:"tokenType,omitempty"`
	Exp       int64  `json:"exp,omitempty"`
	Iat       int64  `json:"iat,omitempty"`
	Sub       string `json:"sub,omitempty"`
	Aud       string `json:"aud,omitempty"`
}

// Handles OAuth 2.0 token introspection (RFC 7662).
// Returns token metadata if active, or {"active": false} if invalid/expired/revoked
func OAuthIntrospectHandler(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")
	wtr.Header().Set("Cache-Control", "no-store")

	// Extract token from Authorization header or request body
	tokenString, err := jwt.ExtractTokenFromRequest(req)
	if err != nil {
		logger.Error("no token found in request", err)
		wtr.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(wtr).Encode(IntrospectResponse{Active: false})
		return
	}

	// Parse claims from token
	claims, err := jwt.ParseClaims(tokenString)
	if err != nil {
		logger.Error("failed to parse claims", err)
		wtr.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(wtr).Encode(IntrospectResponse{Active: false})
		return
	}

	// Check if token is expired
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		logger.Info("token is expired")
		wtr.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(wtr).Encode(IntrospectResponse{Active: false})
		return
	}

	// Extract claims from CustomClaims struct
	scope := ""
	if len(claims.Scopes) > 0 { scope = strings.Join(claims.Scopes, " ") }

	aud := ""
	if len(claims.Audience) > 0 { aud = claims.Audience[0] }

	exp := int64(0)
	if claims.ExpiresAt != nil { exp = claims.ExpiresAt.Unix() }

	iat := int64(0)
	if claims.IssuedAt != nil { iat = claims.IssuedAt.Unix() }

	// TODO: Check if token is revoked in Elasticache
	// if claims.ID != "" && elasticache.Exists("revoked:token:" + claims.ID) {
	//     logger.Info("token has been revoked: " + claims.ID)
	//     wtr.WriteHeader(http.StatusOK)
	//     json.NewEncoder(wtr).Encode(IntrospectResponse{Active: false})
	//     return
	// }

	logger.Info("token introspection successful for subject: " + claims.Subject)

	// Return token metadata per RFC 7662
	wtr.WriteHeader(http.StatusOK)
	json.NewEncoder(wtr).Encode(IntrospectResponse{
		Active:    true,
		Scope:     scope,
		ClientID:  claims.Subject,
		TokenType: "Bearer",
		Exp:       exp,
		Iat:       iat,
		Sub:       claims.Subject,
		Aud:       aud,
	})
}
