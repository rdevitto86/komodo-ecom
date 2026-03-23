package logger

import (
	"context"
	"log/slog"
	"net/http"

	ctxKeys "komodo-forge-sdk-go/http/context"
)

// --- Standard field constructors ---

func Attr(key string, value any) slog.Attr          { return slog.Any(key, value) }
func AttrError(err error) slog.Attr                 { return slog.Any("error", err) }
func AttrRequestID(id string) slog.Attr             { return slog.String("request_id", id) }
func AttrCorrelationID(id string) slog.Attr         { return slog.String("correlation_id", id) }
func AttrUserID(id string) slog.Attr                { return slog.String("user_id", id) }
func AttrSessionID(id string) slog.Attr             { return slog.String("session_id", id) }
func AttrDetails(details map[string]any) slog.Attr  { return slog.Any("details", details) }

// --- Context extraction ---

// FromContext extracts standard log fields from a request context and returns
// them as slog args. Attach to any log call for automatic E2E correlation:
//
//	logger.Info("Processing order", logger.FromContext(ctx)...)
//
// Fields populated (if present in ctx): request_id, correlation_id, user_id, session_id.
func FromContext(ctx context.Context) []any {
	var args []any
	if id, ok := ctx.Value(ctxKeys.REQUEST_ID_KEY).(string); ok && id != "" {
		args = append(args, AttrRequestID(id))
	}
	if id, ok := ctx.Value(ctxKeys.CORRELATION_ID_KEY).(string); ok && id != "" {
		args = append(args, AttrCorrelationID(id))
	}
	if id, ok := ctx.Value(ctxKeys.USER_ID_KEY).(string); ok && id != "" {
		args = append(args, AttrUserID(id))
	}
	if id, ok := ctx.Value(ctxKeys.SESSION_ID_KEY).(string); ok && id != "" {
		args = append(args, AttrSessionID(id))
	}
	return args
}

// --- HTTP helpers ---

// AttrRequest logs method + path only. Headers are intentionally omitted —
// they frequently contain Authorization tokens and other PII.
func AttrRequest(req *http.Request) slog.Attr {
	if req == nil {
		return slog.Any("request", nil)
	}
	return slog.Group("request",
		slog.String("method", req.Method),
		slog.String("path", req.URL.Path),
	)
}

// AttrResponse logs status code only.
func AttrResponse(res *http.Response) slog.Attr {
	if res == nil {
		return slog.Any("response", nil)
	}
	return slog.Group("response",
		slog.Int("status", res.StatusCode),
	)
}
