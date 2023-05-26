package glog

import (
	"context"
	"fmt"
	"golang.org/x/exp/slog"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
)

var (
	defaultLogger atomic.Value
	Home          = func() string {
		home, _ := filepath.Abs("")
		return strings.ReplaceAll(home, "\\", "/") + "/"
	}()
)

func init() {
	defaultLogger.Store(New(os.Stdout, false))
}

type Logger struct {
	*slog.Logger
}

func New(w io.Writer, dev bool) *Logger {
	opts := slog.HandlerOptions{
		AddSource: true,
		Level:     LevelTrace,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				source, _ := a.Value.Any().(*slog.Source)
				file := strings.TrimPrefix(source.File, Home)
				a.Value = slog.StringValue(file + ":" + strconv.Itoa(source.Line))
			} else if a.Key == slog.TimeKey {
				value := a.Value.Time().Format("2006-01-02 15:04:05.000")
				a.Value = slog.StringValue(value)
			} else if a.Key == slog.LevelKey {
				level := Level(a.Value.Any().(slog.Level))
				a.Value = slog.StringValue(level.String())
			}
			return a
		},
	}

	return &Logger{
		Logger: slog.New(&Handler{
			w:       w,
			console: w == os.Stdout || w == os.Stderr,
			Handler: slog.NewJSONHandler(w, &opts),
		}),
	}
}

func (l *Logger) Log(ctx context.Context, level Level, msg any, args ...any) {
	msg = fmt.Sprint(msg)
	l.Logger.Log(ctx, slog.Level(level), msg.(string), args...)

	if level == LevelPanic {
		panic(msg)
	} else if level == LevelFatal {
		os.Exit(1)
	}
}

func (l *Logger) Trace(msg any, args ...any) {
	l.Log(nil, LevelTrace, fmt.Sprint(msg), args...)
}

func (l *Logger) Debug(msg any, args ...any) {
	l.Log(nil, LevelDebug, fmt.Sprint(msg), args...)
}

func (l *Logger) Info(msg any, args ...any) {
	l.Log(nil, LevelInfo, fmt.Sprint(msg), args...)
}

func (l *Logger) Warn(msg any, args ...any) {
	l.Log(nil, LevelWarn, fmt.Sprint(msg), args...)
}

func (l *Logger) Error(msg any, args ...any) {
	l.Log(nil, LevelError, fmt.Sprint(msg), args...)
}

func (l *Logger) Panic(msg any, args ...any) {
	l.Log(nil, LevelPanic, fmt.Sprint(msg), args...)
}

func (l *Logger) Fatal(msg any, args ...any) {
	l.Log(nil, LevelFatal, fmt.Sprint(msg), args...)
}

func (l *Logger) Tracef(msg string, a ...any) {
	l.Log(nil, LevelTrace, fmt.Sprintf(msg, a...))
}

func (l *Logger) Debugf(msg string, a ...any) {
	l.Log(nil, LevelDebug, fmt.Sprintf(msg, a...))
}

func (l *Logger) Infof(msg string, a ...any) {
	l.Log(nil, LevelInfo, fmt.Sprintf(msg, a...))
}

func (l *Logger) Warnf(msg string, a ...any) {
	l.Log(nil, LevelWarn, fmt.Sprintf(msg, a...))
}

func (l *Logger) Errorf(msg string, a ...any) {
	l.Log(nil, LevelError, fmt.Sprintf(msg, a...))
}

func (l *Logger) Panicf(msg string, a ...any) {
	l.Log(nil, LevelPanic, fmt.Sprintf(msg, a...))
}

func (l *Logger) Fatalf(msg string, a ...any) {
	l.Log(nil, LevelFatal, fmt.Sprintf(msg, a...))
}

func (l *Logger) Handler() *Handler {
	return l.Logger.Handler().(*Handler)
}

// top-level package

func Default() *Logger {
	return defaultLogger.Load().(*Logger)
}

func SetDefault(l *Logger) {
	defaultLogger.Store(l)
}

func Trace(msg any, args ...any) {
	Default().Log(nil, LevelTrace, fmt.Sprint(msg), args...)
}

func Debug(msg any, args ...any) {
	Default().Log(nil, LevelDebug, fmt.Sprint(msg), args...)
}

func Info(msg any, args ...any) {
	Default().Log(nil, LevelInfo, fmt.Sprint(msg), args...)
}

func Warn(msg any, args ...any) {
	Default().Log(nil, LevelWarn, fmt.Sprint(msg), args...)
}

func Error(msg any, args ...any) {
	Default().Log(nil, LevelError, fmt.Sprint(msg), args...)
}

func Panic(msg any, args ...any) {
	Default().Log(nil, LevelPanic, fmt.Sprint(msg), args...)
}

func Fatal(msg any, args ...any) {
	Default().Log(nil, LevelFatal, fmt.Sprint(msg), args...)
}

func Tracef(msg string, a ...any) {
	Default().Log(nil, LevelTrace, fmt.Sprintf(msg, a...))
}

func Debugf(msg string, a ...any) {
	Default().Log(nil, LevelDebug, fmt.Sprintf(msg, a...))
}

func Infof(msg string, a ...any) {
	Default().Log(nil, LevelInfo, fmt.Sprintf(msg, a...))
}

func Warnf(msg string, a ...any) {
	Default().Log(nil, LevelWarn, fmt.Sprintf(msg, a...))
}

func Errorf(msg string, a ...any) {
	Default().Log(nil, LevelError, fmt.Sprintf(msg, a...))
}

func Panicf(msg string, a ...any) {
	Default().Log(nil, LevelPanic, fmt.Sprintf(msg, a...))
}

func Fatalf(msg string, a ...any) {
	Default().Log(nil, LevelFatal, fmt.Sprintf(msg, a...))
}
