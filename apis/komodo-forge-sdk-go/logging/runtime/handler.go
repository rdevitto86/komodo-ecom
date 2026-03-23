package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ANSI terminal color codes
const (
	ansiReset   = "\033[0m"
	ansiGray    = "\033[90m"
	ansiCyan    = "\033[36m"
	ansiYellow  = "\033[33m"
	ansiRed     = "\033[31m"
	ansiBoldRed = "\033[1;31m"
)

// skippedBaseFields are emitted in JSON but skipped in string format — they are
// implicit from the log group / process context and clutter every line.
var skippedBaseFields = map[string]bool{
	"service": true,
	"env":     true,
	"version": true,
}

// KomodoTextHandler is a custom slog.Handler that formats log records as:
//
//	2006-01-02T15:04:05Z [LEVEL] requestId | message | key=val key2="val with spaces"
//
// Standard base fields (service, env, version) are omitted from string output —
// they belong in JSON (CloudWatch) and are noise in a terminal line.
// ANSI color is enabled per the color flag set at construction.
type KomodoTextHandler struct {
	mu       sync.Mutex
	w        io.Writer
	color    bool
	level    slog.Leveler
	preAttrs []slog.Attr // attrs set via .With() (e.g., service, env, version)
}

// NewKomodoTextHandler constructs a handler writing to w.
// Set color=true for terminal output, false for CI/pipe environments.
func NewKomodoTextHandler(w io.Writer, color bool, level slog.Leveler) *KomodoTextHandler {
	return &KomodoTextHandler{w: w, color: color, level: level}
}

func (h *KomodoTextHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

// WithAttrs returns a new handler with the provided attrs prepended to all future records.
// A new handler struct is constructed rather than copied to avoid copying the mutex.
func (h *KomodoTextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h2 := &KomodoTextHandler{
		w:     h.w,
		color: h.color,
		level: h.level,
	}
	h2.preAttrs = make([]slog.Attr, len(h.preAttrs)+len(attrs))
	copy(h2.preAttrs, h.preAttrs)
	copy(h2.preAttrs[len(h.preAttrs):], attrs)
	return h2
}

// WithGroup is a no-op for the text handler — groups are flattened with dot notation.
func (h *KomodoTextHandler) WithGroup(name string) slog.Handler {
	return h
}

func (h *KomodoTextHandler) Handle(_ context.Context, rec slog.Record) error {
	var buf bytes.Buffer

	requestID := "-"
	var parts []string

	collectAttr := func(attr slog.Attr) {
		if skippedBaseFields[attr.Key] {
			return
		}
		if attr.Key == "request_id" {
			requestID = attr.Value.String()
			return
		}
		if s := formatAttr(attr); s != "" {
			parts = append(parts, s)
		}
	}

	for _, a := range h.preAttrs {
		collectAttr(a)
	}
	rec.Attrs(func(a slog.Attr) bool {
		collectAttr(a)
		return true
	})

	// timestamp
	if h.color {
		buf.WriteString(ansiGray)
	}
	buf.WriteString(rec.Time.UTC().Format(time.RFC3339))
	if h.color {
		buf.WriteString(ansiReset)
	}
	buf.WriteByte(' ')

	// [LEVEL]
	buf.WriteString(h.coloredLevel(rec.Level))
	buf.WriteByte(' ')

	// requestId
	if h.color {
		buf.WriteString(ansiGray)
	}
	buf.WriteString(requestID)
	if h.color {
		buf.WriteString(ansiReset)
	}
	buf.WriteString(" | ")

	// message
	buf.WriteString(rec.Message)

	// details as logfmt
	if len(parts) > 0 {
		buf.WriteString(" | ")
		if h.color {
			buf.WriteString(ansiGray)
		}
		buf.WriteString(strings.Join(parts, " "))
		if h.color {
			buf.WriteString(ansiReset)
		}
	}

	buf.WriteByte('\n')

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.w.Write(buf.Bytes())
	return err
}

func (h *KomodoTextHandler) coloredLevel(level slog.Level) string {
	var color, label string
	switch {
	case level > slog.LevelError:
		color, label = ansiBoldRed, "FATAL"
	case level >= slog.LevelError:
		color, label = ansiRed, "ERROR"
	case level >= slog.LevelWarn:
		color, label = ansiYellow, "WARN"
	case level >= slog.LevelInfo:
		color, label = ansiCyan, "INFO"
	default:
		color, label = ansiGray, "DEBUG"
	}
	if !h.color {
		return "[" + label + "]"
	}
	return color + "[" + label + "]" + ansiReset
}

// formatAttr renders a slog.Attr as a logfmt key=value pair.
// Complex values are JSON-encoded inline and truncated at 200 chars.
// Groups are flattened with dot notation (parent.child=value).
func formatAttr(attr slog.Attr) string {
	k := attr.Key
	v := attr.Value.Resolve()

	switch v.Kind() {
	case slog.KindString:
		s := v.String()
		if needsQuoting(s) {
			return k + "=" + strconv.Quote(s)
		}
		return k + "=" + s
	case slog.KindInt64:
		return k + "=" + strconv.FormatInt(v.Int64(), 10)
	case slog.KindUint64:
		return k + "=" + strconv.FormatUint(v.Uint64(), 10)
	case slog.KindFloat64:
		return fmt.Sprintf("%s=%g", k, v.Float64())
	case slog.KindBool:
		return k + "=" + strconv.FormatBool(v.Bool())
	case slog.KindGroup:
		sub := v.Group()
		if len(sub) == 0 {
			return ""
		}
		parts := make([]string, 0, len(sub))
		for _, a := range sub {
			if s := formatAttr(slog.Attr{Key: k + "." + a.Key, Value: a.Value}); s != "" {
				parts = append(parts, s)
			}
		}
		return strings.Join(parts, " ")
	case slog.KindAny:
		any := v.Any()
		if err, ok := any.(error); ok {
			return k + "=" + strconv.Quote(err.Error())
		}
		b, err := json.Marshal(any)
		if err != nil {
			return k + "=<error>"
		}
		s := string(b)
		if len(s) > 200 {
			s = s[:197] + "..."
		}
		if needsQuoting(s) {
			return k + "=" + strconv.Quote(s)
		}
		return k + "=" + s
	default:
		s := v.String()
		if needsQuoting(s) {
			return k + "=" + strconv.Quote(s)
		}
		return k + "=" + s
	}
}

func needsQuoting(s string) bool {
	return strings.ContainsAny(s, " \t\n\r\"=")
}
