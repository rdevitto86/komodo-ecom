package requestid

import (
	"context"
	"net/http"

	ctxKeys "komodo-forge-sdk-go/http/context"
	httpReq "komodo-forge-sdk-go/http/request"
)

// RequestIDMiddleware ensures each request has a unique X-Request-ID in both
// header and context. Priority: header (client-supplied) > context > generated.
//
// Also propagates X-Correlation-ID from the client — this is the browser's
// session fingerprint (correlationId) used to tie frontend and backend logs
// to the same browsing session.
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Request ID
		var reqID string
		if rid := req.Header.Get("X-Request-ID"); rid != "" {
			reqID = rid
		} else if rid, ok := req.Context().Value(ctxKeys.REQUEST_ID_KEY).(string); ok && rid != "" {
			reqID = rid
		} else {
			reqID = httpReq.GenerateRequestId()
		}

		req.Header.Set("X-Request-ID", reqID)
		w.Header().Set("X-Request-ID", reqID)
		ctx := context.WithValue(req.Context(), ctxKeys.REQUEST_ID_KEY, reqID)

		// Correlation ID (browser fingerprint — client-generated, server echoes back)
		if cid := req.Header.Get("X-Correlation-ID"); cid != "" {
			ctx = context.WithValue(ctx, ctxKeys.CORRELATION_ID_KEY, cid)
			w.Header().Set("X-Correlation-ID", cid)
		}

		next.ServeHTTP(w, req.WithContext(ctx))
	})
}
