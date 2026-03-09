package logger

import (
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/lmittmann/tint"
)

var (
	slogger *slog.Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logLevel = &slog.LevelVar{}
	initOnce sync.Once
)

func Init(name string, lvl string, env string) {
	initOnce.Do(func() {
		isLocal := func() bool {
			env = strings.ToLower(env)
			return env == "local" || env == "dev" || env == "development"
		}

		logLevel.Set(parseLevel(lvl))

		var handler slog.Handler
		if isLocal() {
			handler = tint.NewHandler(os.Stdout, &tint.Options{Level: logLevel, TimeFormat: "15:04:05"})
		} else {
			handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
		}

		slogger = slog.New(&RedactingLogger{Handler: handler}).With(slog.String("app", name))
		slog.SetDefault(slogger)
	})
}

func Debug(msg string, args ...any) { slogger.Debug(msg, args...) }
func Info(msg string, args ...any) { slogger.Info(msg, args...) }
func Warn(msg string, args ...any) { slogger.Warn(msg, args...) }
func Error(msg string, err error, args ...any) {
	if err != nil { args = append(args, AttrError(err)) }
	slogger.Error(msg, args...)
}
func Fatal(msg string, err error, args ...any) {
	if err != nil { args = append(args, AttrError(err)) }
	slogger.Error(msg, args...)
}

func SetLevel(level string) { logLevel.Set(parseLevel(level)) }

func parseLevel(lvl string) slog.Level {
	switch strings.ToLower(lvl) {
		case "debug": return slog.LevelDebug
		case "info": 	return slog.LevelInfo
		case "warn": 	return slog.LevelWarn
		default:
			if lvl == "" { return slog.LevelInfo }
			return slog.LevelError
	}
}
