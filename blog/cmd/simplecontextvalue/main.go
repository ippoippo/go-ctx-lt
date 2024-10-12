package main

import (
	"context"
	"log/slog"
)

func main() {
	ctx := context.Background()
	ctxWithValue := context.WithValue(ctx, "key", "value")

	slog.Info("ctx",
		slog.Any("context-key", ctx.Value("key")),
		slog.Any("context-not-exist-key", ctx.Value("other")))

	slog.Info("ctxWithValue",
		slog.Any("context-key", ctxWithValue.Value("key")),
		slog.Any("context-not-exist-key", ctxWithValue.Value("other")))
}
