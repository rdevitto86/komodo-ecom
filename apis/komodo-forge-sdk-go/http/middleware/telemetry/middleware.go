package telemetry

import (
	"fmt"
	httpErr "komodo-forge-sdk-go/http/errors"
	httpReq "komodo-forge-sdk-go/http/request"
	httpRes "komodo-forge-sdk-go/http/response"
	logger "komodo-forge-sdk-go/logging/runtime"
	"net/http"
	"time"
)

func TelemetryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		resWtr := &httpRes.ResponseWriter{ResponseWriter: wtr}
		start := time.Now()

		defer func() {
			ms := time.Since(start).Milliseconds()

			// Recover from panics and ensure a 500 is sent if nothing written.
			if rec := recover(); rec != nil {
				reqID := httpReq.GetRequestID(req)
				_ = reqID
				
				// Safely check status
				status := resWtr.Status
				if status == 0 {
					httpErr.SendError(
						wtr, req, httpErr.Global.Internal, httpErr.WithDetail("error occured while logging telemetry"),
					)
				}
				
				logger.Error("telemetry panicked!", fmt.Errorf("telemetry panicked: %v", rec))
				return
			}

			status := resWtr.Status
			bytesWritten := resWtr.BytesWritten
			if status == 0 {
				status = http.StatusOK
			}

			// Get request ID safely
			reqID := httpReq.GetRequestID(req)

			payload := map[string]any{
				"request_id": reqID,
				"method":     req.Method,
				"path":       req.URL.Path,
				"query":      req.URL.RawQuery,
				"status":     status,
				"bytes":     	bytesWritten,
				"latency_ms": ms,
				"ip":         req.RemoteAddr,
				"user_agent": req.UserAgent(),
				"referer":    req.Referer(),
				"proto":     	req.Proto,
				"host":      	req.Host,
				"start_time": start.UTC().Format(time.RFC3339Nano),
				"finish_time": time.Now().UTC().Format(time.RFC3339Nano),
			}

			if status >= 400 {
				logger.Error("telemetry request failed", fmt.Errorf("telemetry request failed: %v", payload))
			} else {
				logger.Info("telemetry request completed")
			}
		}()

		next.ServeHTTP(resWtr, req)
	})
}
