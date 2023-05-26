package glog

import (
	"context"
	"fmt"
	"github.com/gookit/color"
	"golang.org/x/exp/slog"
	"io"
	"runtime"
	"strconv"
	"strings"
)

type Handler struct {
	slog.Handler
	w       io.Writer
	attrs   []slog.Attr
	group   slog.Attr
	console bool
}

func (h *Handler) Writer() io.Writer {
	return h.w
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.Handler.Enabled(ctx, level)
}

func (h *Handler) Handle(ctx context.Context, r slog.Record) (err error) {
	pc, file, line, _ := runtime.Caller(5)
	if r.PC = pc; !h.console {
		if err = h.Handler.Handle(ctx, r); err != nil {
			return err
		}
	}

	level := Level(r.Level)
	prefix := r.Time.Format("2006-01-02 15:04:05.000")
	prefix = color.HEX("#A9B7C6").Sprint(prefix)
	prefix += " | " + level.ColorString()

	var attrs []slog.Attr
	r.Attrs(func(a slog.Attr) bool {
		attrs = append(attrs, a)
		return true
	})

	if h.group.Key != "" {
		h.addGroupAttr(attrs...)
		attrs = []slog.Attr{h.group}
	}

	s := attrString(append(h.attrs, attrs...)...)
	s = color.Cyan.Sprint(s)

	source := strings.TrimPrefix(file+":"+strconv.Itoa(line), Home)
	fmt.Printf("%s | %s > %s %s\n", prefix, source, r.Message, s)

	return
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandler := h.Handler.WithAttrs(attrs)
	handler := &Handler{Handler: newHandler, attrs: h.attrs, group: h.group}

	if handler.group.Key != "" {
		handler.addGroupAttr(attrs...)
	} else {
		handler.attrs = append(h.attrs, attrs...)
	}

	return handler
}

func (h *Handler) WithGroup(name string) slog.Handler {
	newHandler := h.Handler.WithGroup(name)
	handler := &Handler{Handler: newHandler, attrs: h.attrs, group: h.group}

	if group := slog.Group(name); h.group.Key == "" {
		handler.group = group
	} else {
		handler.addGroupAttr(group)
	}

	return handler
}

func (h *Handler) addGroupAttr(attrs ...slog.Attr) {
	if v := &lastGroup(&h.group).Value; v.Kind() == slog.KindGroup {
		*v = slog.GroupValue(append(v.Group(), attrs...)...)
	}
}
