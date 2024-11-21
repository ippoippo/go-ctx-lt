package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

func main() {
	ctx := context.Background()
	err := getFromUrl(ctx, 1)
	if err != nil {
		// Note: ctx.Err() will NOT be set in the parent context
		slog.Error("getFromUrl(ctx, 1) has returned error: ", slog.String("err", err.Error()), slog.Any("ctx-err", ctx.Err()))
	}
	err = getFromUrl(ctx, 5)
	if err != nil {
		// Note: ctx.Err() will NOT be set in the parent context
		slog.Error("getFromUrl(ctx, 5) has returned error: ", slog.String("err", err.Error()), slog.Any("ctx-err", ctx.Err()))
	}
	slog.Info("completed")
}

func getFromUrl(ctx context.Context, delayValue int) error {
	slog.Info("getFromUrl entry: ", slog.Int("delay-value", delayValue))
	start := time.Now()

	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(3*time.Second))
	defer cancel()

	url := fmt.Sprintf("http://0.0.0.0:80/delay/%d", delayValue)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err // We don't expect this to happen
	}

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		// Note: ctx.Err() will be set here for the child context.
		slog.Error("httpClient.Do(req) has returned error: ", slog.String("err", err.Error()), slog.Any("ctx-err", ctx.Err()))
		return err
	}
	defer resp.Body.Close()
	defer io.Copy(io.Discard, resp.Body)

	slog.Info("getFromUrl completed: ", slog.Int64("since-ms", int64(time.Since(start)/time.Millisecond)))

	return nil
}
