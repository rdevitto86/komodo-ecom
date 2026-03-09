package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"komodo-forge-sdk-go/crypto/jwt"
	httpErr "komodo-forge-sdk-go/http/errors"
	logger "komodo-forge-sdk-go/logging/runtime"
)

type RevokeRequest struct {
	Token         string `json:"token"`
	TokenTypeHint string `json:"token_type_hint,omitempty"` // "access_token" or "refresh_token"
}

// Handles OAuth 2.0 token revocation (RFC 7009).
// Revokes access or refresh tokens
func OAuthRevokeHandler(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")
	wtr.Header().Set("Cache-Control", "no-store")

	// Parse request body
	var reqBody RevokeRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		logger.Error("failed to parse request body", err)
		httpErr.SendError(
			wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail("failed to parse request body"),
		)
		return
	}

	if reqBody.Token == "" {
		logger.Error("missing token parameter", fmt.Errorf("missing token parameter"))
		httpErr.SendError(
			wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail("missing token parameter"),
		)
		return
	}

	// Parse claims from token
	claims, err := jwt.ParseClaims(reqBody.Token)
	if err != nil {
		// Per RFC 7009, return 200 OK even if token is invalid
		// (prevents information disclosure about token validity)
		logger.Warn("invalid token submitted for revocation")
		wtr.WriteHeader(http.StatusOK)
		json.NewEncoder(wtr).Encode(map[string]bool{"revoked": true})
		return
	}

	// Extract JTI (token ID) from claims
	jti := claims.ID
	if jti == "" {
		// Token without JTI cannot be revoked (shouldn't happen in our system)
		logger.Warn("token missing JTI claim")
		wtr.WriteHeader(http.StatusOK)
		json.NewEncoder(wtr).Encode(map[string]bool{"revoked": true})
		return
	}

	// Calculate TTL from expiration time
	ttl := time.Duration(0)
	if claims.ExpiresAt != nil {
		ttl = time.Until(claims.ExpiresAt.Time)
	}
	if ttl <= 0 {
		// Token already expired, no need to revoke
		logger.Info("token already expired, no revocation needed")
		wtr.WriteHeader(http.StatusOK)
		json.NewEncoder(wtr).Encode(map[string]bool{"revoked": true})
		return
	}

	// TODO: Store revoked token in Elasticache with TTL
	// revokeKey := "revoked:token:" + jti
	// if err := elasticache.SetCacheItem(revokeKey, clientID, ttl); err != nil {
	// 	logger.Error("failed to revoke token in cache", err)
	// 	errors.WriteErrorResponse(
	// 		wtr,
	// 		req,
	// 		http.StatusInternalServerError,
	// 		"server_error",
	// 		errCodes.ERR_INTERNAL_SERVER,
	// 	)
	// 	return
	// }

	logger.Info("token revoked successfully for subject: " + claims.Subject + ", JTI: " + jti)

	// Per RFC 7009, return 200 OK with empty response (or small JSON)
	wtr.WriteHeader(http.StatusOK)
	json.NewEncoder(wtr).Encode(map[string]interface{}{
		"revoked":    true,
		"revoked_at": time.Now().Unix(),
	})
}
