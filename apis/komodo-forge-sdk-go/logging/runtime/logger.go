package logger

import (
	"log/slog"
	"os"
	"strings"
	"sync"
)

var (
	slogger  *slog.Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logLevel               = &slog.LevelVar{}
	initOnce sync.Once
)

// Init configures the global logger. Must be called once at service startup.
//
// Local/dev environments (env == "local"|"dev"|"development") use the
// KomodoTextHandler which formats as:
//
//	2006-01-02T15:04:05Z [LEVEL] requestId | message | key=val ...
//
// All other environments use the JSON handler (stdout → CloudWatch).
//
// The version parameter is optional — pass config.GetConfigValue("VERSION") or
// omit it; defaults to "unknown". Using a variadic keeps callers that already
// use the 3-argument form working without changes.
func Init(name string, lvl string, env string, version ...string) {
	initOnce.Do(func() {
		ver := "unknown"
		if len(version) > 0 && version[0] != "" {
			ver = version[0]
		}

		logLevel.Set(parseLevel(lvl))

		var handler slog.Handler
		if isLocalEnv(env) {
			handler = NewKomodoTextHandler(os.Stdout, true, logLevel)
		} else {
			handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
		}

		slogger = slog.New(&RedactingLogger{Handler: handler}).With(
			slog.String("service", name),
			slog.String("env", env),
			slog.String("version", ver),
		)
		slog.SetDefault(slogger)
	})
}

func Debug(msg string, args ...any) { slogger.Debug(msg, args...) }
func Info(msg string, args ...any)  { slogger.Info(msg, args...) }
func Warn(msg string, args ...any)  { slogger.Warn(msg, args...) }

func Error(msg string, err error, args ...any) {
	if err != nil {
		args = append(args, AttrError(err))
	}
	slogger.Error(msg, args...)
}

func Fatal(msg string, err error, args ...any) {
	if err != nil {
		args = append(args, AttrError(err))
	}
	// slog has no Fatal level — use Error. Caller is responsible for os.Exit.
	slogger.Error(msg, args...)
}

func SetLevel(level string) { logLevel.Set(parseLevel(level)) }

func isLocalEnv(env string) bool {
	e := strings.ToLower(env)
	return e == "local" || e == "dev" || e == "development"
}

func parseLevel(lvl string) slog.Level {
	switch strings.ToLower(lvl) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	default:
		if lvl == "" {
			return slog.LevelInfo
		}
		return slog.LevelError
	}
}
