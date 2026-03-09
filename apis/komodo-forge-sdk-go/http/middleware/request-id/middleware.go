package requestid

import (
	"context"
	ctxKeys "komodo-forge-sdk-go/http/context"
	httpReq "komodo-forge-sdk-go/http/request"
	"net/http"
)

// Ensures each request has a unique X-Request-ID in both header and context
// Priority: Header (external) > Context (middleware) > Generated (new)
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		var reqID string
		if rid := req.Header.Get("X-Request-ID"); rid != "" {
			reqID = rid
		} else if rid, ok := req.Context().Value(ctxKeys.REQUEST_ID_KEY).(string); ok && rid != "" {
			reqID = rid
		} else {
			reqID = httpReq.GenerateRequestId()
		}

		req.Header.Set("X-Request-ID", reqID)
		ctx := context.WithValue(req.Context(), ctxKeys.REQUEST_ID_KEY, reqID)
		wtr.Header().Set("X-Request-ID", reqID)

		next.ServeHTTP(wtr, req.WithContext(ctx))
	})
}