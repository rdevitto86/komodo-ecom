package auth

import (
	"context"
	"komodo-forge-sdk-go/crypto/jwt"
	ctxKeys "komodo-forge-sdk-go/http/context"
	httpErr "komodo-forge-sdk-go/http/errors"
	logger "komodo-forge-sdk-go/logging/runtime"
	"net/http"
)

// Handles OAuth2 + JWT Bearer token authentication
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		tokenString, err := jwt.ExtractTokenFromRequest(req)
		if err != nil {
			logger.Error("failed to extract token", err)
			httpErr.SendError(wtr, req, httpErr.Auth.InvalidToken, httpErr.WithDetail(err.Error()))
			return
		}

		valid, err := jwt.ValidateToken(tokenString)
		if !valid || err != nil {
			logger.Error("token validation failed", err)
			httpErr.SendError(wtr, req, httpErr.Auth.InvalidToken, httpErr.WithDetail(err.Error()))
			return
		}

		claims, err := jwt.ParseClaims(tokenString)
		if err != nil {
			logger.Error("failed to parse claims", err)
			httpErr.SendError(wtr, req, httpErr.Auth.InvalidToken, httpErr.WithDetail(err.Error()))
			return
		}

		ctx := context.WithValue(req.Context(), ctxKeys.AUTH_VALID_KEY, true)
		
		if claims.Subject != "" {
			ctx = context.WithValue(ctx, ctxKeys.USER_ID_KEY, claims.Subject)
		}
		if claims.ID != "" {
			ctx = context.WithValue(ctx, ctxKeys.SESSION_ID_KEY, claims.ID)
		}

		isUIRequest := len(claims.Scopes) == 0 || claims.IsAdmin
		isAPIRequest := len(claims.Scopes) > 0

		if isUIRequest {
			ctx = context.WithValue(ctx, ctxKeys.REQUEST_TYPE_KEY, "ui")
		}
		if isAPIRequest {
			ctx = context.WithValue(ctx, ctxKeys.REQUEST_TYPE_KEY, "api")
			ctx = context.WithValue(ctx, ctxKeys.SCOPES_KEY, claims.Scopes)
		}
		if claims.IsAdmin {
			ctx = context.WithValue(ctx, ctxKeys.IS_ADMIN_KEY, true)
		}

		next.ServeHTTP(wtr, req.WithContext(ctx))
	})
}
