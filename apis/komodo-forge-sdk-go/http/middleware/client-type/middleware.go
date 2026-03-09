package clienttype

import (
	"context"
	ctxKeys "komodo-forge-sdk-go/http/context"
	"net/http"
)

const (
	ClientTypeAPI string = "api"
	ClientTypeBrowser string = "browser"
)

// Detects whether request is from API client or browser
// and stores the result in context for downstream middleware to use
func ClientTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		authHeader := req.Header.Get("Authorization")
		hasReferer := req.Header.Get("Referer") != ""
		hasCookie := req.Header.Get("Cookie") != ""
		
		clientType := ClientTypeBrowser
		if authHeader != "" && !hasReferer && !hasCookie {
			clientType = ClientTypeAPI
		}
	
		ctx := context.WithValue(req.Context(), ctxKeys.CLIENT_TYPE_KEY, clientType)
		next.ServeHTTP(wtr, req.WithContext(ctx))
	})
}
