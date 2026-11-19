package lutil

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-colorable"
)

// Minimum log level can be set through MINIMUM_LOG_LEVEL (e.g. "debug", "info", "warn", "error").
var logLevelMap = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
}

type prefixHandler struct {
	h      slog.Handler
	prefix string
}

func (p *prefixHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return p.h.Enabled(ctx, level)
}

const blue = "\x1b[38;5;67m" // slightly brighter muted blue
const reset = "\x1b[0m"

func (p *prefixHandler) Handle(ctx context.Context, r slog.Record) error {
	r.Message = blue + p.prefix + reset + r.Message
	return p.h.Handle(ctx, r)
}

func (p *prefixHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &prefixHandler{h: p.h.WithAttrs(attrs), prefix: p.prefix}
}

func (p *prefixHandler) WithGroup(name string) slog.Handler {
	return &prefixHandler{h: p.h.WithGroup(name), prefix: p.prefix}
}

// Magic applies its environment after package-level loggers have been created.
// Resolve the level when logging so those loggers also honor the local setting.
type environmentLogLevel struct{}

func (environmentLogLevel) Level() slog.Level {
	if level, ok := logLevelMap[os.Getenv("MINIMUM_LOG_LEVEL")]; ok {
		return level
	} else {
		fmt.Println("WARNING: Invalid log level set in MINIMUM_LOG_LEVEL. Correct values are: debug, info, warn or error!")
	}
	return slog.LevelInfo
}

// NewLogger returns a colored logger.
func NewLogger(part string) *slog.Logger {
	return slog.New(&prefixHandler{
		h: tint.NewTextHandler(colorable.NewColorable(os.Stderr), &tint.Options{
			Level:      environmentLogLevel{},
			TimeFormat: time.Kitchen,
		}),
		prefix: part + " > ",
	})
}

// Fatal logs a message at the error level and then exits the program with a status code of 1.
func Fatal(l *slog.Logger, msg string, args ...any) {
	l.Error(msg, args...)
	os.Exit(1)
}

// Prints a warning with the error and message, if it is not nil.
func PotentialWarn(l *slog.Logger, msg string, err error) {
	if err != nil {
		l.Warn(msg, "error", err)
	}
}
