package main

import (
	"context"
	"log/slog"
	"sync"
)

type CopyHandler struct {
	mu  *sync.Mutex
	out []slog.Handler
}

func NewCopyHandler(handlers ...slog.Handler) *CopyHandler {
	return &CopyHandler{
		out: handlers,
		mu:  &sync.Mutex{},
	}
}

func (h *CopyHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}

func (h *CopyHandler) Handle(ctx context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, destHandler := range h.out {
		err := destHandler.Handle(ctx, r)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *CopyHandler) WithGroup(name string) slog.Handler {
	// call WithGroup on the underlying handlers
	// we should not make modification the receiver, we return a copy
	if name == "" {
		return h
	}
	h2 := *h
	h2.out = make([]slog.Handler, len(h.out))
	for i, h := range h.out {
		h2.out[i] = h.WithGroup(name)
	}
	return &h2
}

func (h *CopyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// call WithAttrs on the underlying handlers
	// we should not make modification the receiver, we return a copy
	if len(attrs) == 0 {
		return h
	}
	h2 := *h
	h2.out = make([]slog.Handler, len(h.out))
	for i, h := range h.out {
		h2.out[i] = h.WithAttrs(attrs)
	}
	return &h2
}
