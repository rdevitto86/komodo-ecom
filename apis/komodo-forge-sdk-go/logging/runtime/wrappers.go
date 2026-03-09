package logger

import (
	"context"
	"log/slog"
	"net/http"
)

func Attr(key string, value any) slog.Attr { return slog.Any(key, value) }
func AttrError(err error) slog.Attr { return slog.Any("error", err) }
func AttrContext(ctx context.Context) slog.Attr { return slog.Any("context", ctx) }

func AttrRequest(req *http.Request) slog.Attr {
	if req == nil { return slog.Any("request", nil) }
	return slog.Group("request",
		slog.String("method", req.Method),
		slog.String("url", req.URL.String()),
		slog.Any("headers", req.Header),
	)
}

func AttrResponse(res *http.Response) slog.Attr {
	if res == nil { return slog.Any("response", nil) }
	return slog.Group("response",
		slog.Int("status", res.StatusCode),
		slog.Any("headers", res.Header),
	)
}
