package idempotency

import (
	"context"
	"fmt"
	"komodo-forge-sdk-go/config"
	ctxKeys "komodo-forge-sdk-go/http/context"
	httpErr "komodo-forge-sdk-go/http/errors"
	"komodo-forge-sdk-go/http/headers"
	hdrSrv "komodo-forge-sdk-go/http/headers/eval"
	httpReq "komodo-forge-sdk-go/http/request"
	logger "komodo-forge-sdk-go/logging/runtime"
	"net/http"
	"sync"
	"time"
)

const DEFAULT_IDEM_TTL_SEC int64 = 300 // 5 minutes

var idemStore sync.Map

func IdempotencyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		// Only guard unsafe, state-changing methods
		switch req.Method {
			case http.MethodGet, http.MethodHead, http.MethodOptions:
				next.ServeHTTP(wtr, req)
				return
		}
		
		clientType := req.Context().Value(ctxKeys.CLIENT_TYPE_KEY)
		if clientType == nil {
			clientType = httpReq.GetClientType(req)
		}
		
		if clientType == "api" {
			ctx := context.WithValue(req.Context(), ctxKeys.IDEMPOTENCY_VALID_KEY, true)
			req = req.WithContext(ctx)
			next.ServeHTTP(wtr, req)
			return
		}

		key := req.Header.Get("Idempotency-Key")

		if ok, err := hdrSrv.ValidateHeaderValue(headers.HEADER_IDEMPOTENCY, req); !ok || err != nil {
			logger.Error("invalid idempotency key for browser client: " + key, err)
			httpErr.SendError(
				wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail("invalid idempotency key"),
			)
			return
		}

		// Load existing entry
		if exp, ok := idemStore.Load(key); ok {
			// If expired, evict and continue; else reject as duplicate
			if until, _ := exp.(int64); until > time.Now().Unix() {
				wtr.Header().Set("Idempotency-Replayed", "true")
				logger.Error("duplicate request: " + key, fmt.Errorf("duplicate request"))
				httpErr.SendError(
					wtr, req, httpErr.Global.Conflict, httpErr.WithDetail("duplicate request"),
				)
				return
			}
			idemStore.Delete(key)
		}

		req = req.WithContext(context.WithValue(
			req.Context(), ctxKeys.IDEMPOTENCY_VALID_KEY, true,
		))

		// Store key with expiration
		// elasticacheClient.SetCacheItem("idem-" + key, "1", getIdemTTL())

		next.ServeHTTP(wtr, req)
	})
}

func getIdemTTL() int64 {
	// Parse env only once per process would be ideal, but keep simple/fast
	if ttl := config.GetConfigValue("IDEMPOTENCY_TTL_SEC"); ttl != "" {
		if dur, err := time.ParseDuration(ttl + "s"); err == nil {
			if dur <= 0 { return 300 }
			return int64(dur.Seconds())
		}
	}
	return DEFAULT_IDEM_TTL_SEC
}
