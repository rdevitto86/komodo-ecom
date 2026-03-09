package context

import (
	"context"
	ctxKeys "komodo-forge-sdk-go/http/context"
	httpReq "komodo-forge-sdk-go/http/request"
	"net/http"
	"time"
)

// Enriches request context with common values
func ContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		var reqID string
		if rid := req.Header.Get("X-Request-ID"); rid != "" {
			reqID = rid
		} else if rid, ok := ctx.Value(ctxKeys.REQUEST_ID_KEY).(string); ok && rid != "" {
			reqID = rid
		} else {
			reqID = httpReq.GenerateRequestId()
		}

		ctx = context.WithValue(ctx, ctxKeys.REQUEST_ID_KEY, reqID)
		req.Header.Set("X-Request-ID", reqID)
		wtr.Header().Set("X-Request-ID", reqID)

		ctx = context.WithValue(ctx, ctxKeys.START_TIME_KEY, time.Now().UTC())
		ctx = context.WithValue(ctx, ctxKeys.VERSION_KEY, httpReq.GetAPIVersion(req))
		ctx = context.WithValue(ctx, ctxKeys.URI_KEY, httpReq.GetAPIRoute(req))
		ctx = context.WithValue(ctx, ctxKeys.METHOD_KEY, req.Method)
		ctx = context.WithValue(ctx, ctxKeys.PATH_PARAMS_KEY, httpReq.GetPathParams(req))
		ctx = context.WithValue(ctx, ctxKeys.QUERY_PARAMS_KEY, httpReq.GetQueryParams(req))
		// ctx = context.WithValue(ctx, ctxKeys.CLIENT_IP_KEY, utils.GetClientIP(req))
		ctx = context.WithValue(ctx, ctxKeys.CLIENT_TYPE_KEY, httpReq.GetClientType(req))
		ctx = context.WithValue(ctx, ctxKeys.USER_AGENT_KEY, req.UserAgent())

		next.ServeHTTP(wtr, req.WithContext(ctx))
	})
}
