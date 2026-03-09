package ratelimiting

import (
	"fmt"
	httpErr "komodo-forge-sdk-go/http/errors"
	httpReq "komodo-forge-sdk-go/http/request"
	rl "komodo-forge-sdk-go/http/utils/rate_limiter"
	logger "komodo-forge-sdk-go/logging/runtime"
	"net/http"
	"strconv"
)

// RateLimiterMiddleware delegates core logic to services/rate_limiter.
func RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		key := httpReq.GetClientKey(req)
		allowed, wait, err := rl.Allow(req.Context(), key)

		if err != nil {
			if rl.ShouldFailOpen() {
				logger.Error("rate limiter failing open for client: " + key, err)
			} else {
				logger.Error("rate limiter failed for client: " + key, err)
				httpErr.SendError(
					wtr, req, httpErr.Global.Internal, httpErr.WithDetail("internal rate limiter error"),
				)
				return
			}
		} else if !allowed {
			if wait > 0 {
				wtr.Header().Set("Retry-After", strconv.Itoa(int(wait.Seconds() + 0.5)))
			}
			logger.Error("rate limit exceeded for client: " + key, fmt.Errorf("rate limit exceeded"))
			httpErr.SendError(
				wtr, req, httpErr.Global.TooManyRequests, httpErr.WithDetail("rate limit exceeded"),
			)
			return
		}

		next.ServeHTTP(wtr, req)
	})
}
