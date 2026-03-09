package logger

import (
	"context"
	"komodo-forge-sdk-go/http/utils/redaction"
	"log/slog"
)

type RedactingLogger struct { slog.Handler }

func (rl *RedactingLogger) Handle(ctx context.Context, rec slog.Record) error {
	clean := rec.Clone()
	rec.Attrs(func(attr slog.Attr) bool {
		clean.AddAttrs(slog.Any(attr.Key, redaction.RedactPair(attr.Key, attr.Value.Any())))
		return true
	})
	return rl.Handler.Handle(ctx, clean)
}
