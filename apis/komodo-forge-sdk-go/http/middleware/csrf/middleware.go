package csrf

import (
	"context"
	ctxKeys "komodo-forge-sdk-go/http/context"
	httpErr "komodo-forge-sdk-go/http/errors"
	"komodo-forge-sdk-go/http/headers"
	hdrSrv "komodo-forge-sdk-go/http/headers/eval"
	httpReq "komodo-forge-sdk-go/http/request"
	logger "komodo-forge-sdk-go/logging/runtime"
	"net/http"
)

func CSRFMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		switch req.Method {
			case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
				clientType := req.Context().Value(ctxKeys.CLIENT_TYPE_KEY)
				if clientType == nil {
					clientType = httpReq.GetClientType(req)
				}

				if clientType == "api" {
					ctx := context.WithValue(req.Context(), ctxKeys.CSRF_TOKEN_KEY, "api-client-exempt")
					ctx = context.WithValue(ctx, ctxKeys.CSRF_VALID_KEY, true)
					req = req.WithContext(ctx)
					next.ServeHTTP(wtr, req)
					return
				}
				
				// Browser client - require CSRF token
				if ok, err := hdrSrv.ValidateHeaderValue(headers.HEADER_X_CSRF_TOKEN, req); !ok || err != nil {
					logger.Error("invalid or missing CSRF token for browser client", err)
					httpErr.SendError(wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail("invalid CSRF token"))
					return
				}
		}

		ctx := context.WithValue(req.Context(), ctxKeys.CSRF_TOKEN_KEY, "")
		ctx = context.WithValue(ctx, ctxKeys.CSRF_VALID_KEY, true)
		req = req.WithContext(ctx)

		next.ServeHTTP(wtr, req)
	})
}
