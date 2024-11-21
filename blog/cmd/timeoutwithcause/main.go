package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type apiTimeoutError struct {
	timeoutDuration time.Duration
}

func (a *apiTimeoutError) Error() string {
	return fmt.Sprintf("timeout after %v seconds", a.timeoutDuration.Seconds())
}

func newApiTimeout(timeoutAfter time.Duration) *apiTimeoutError {
	return &apiTimeoutError{
		timeoutDuration: timeoutAfter,
	}
}

func main() {
	ctx := context.Background()
	checkError(ctx, getFromUrl(ctx, 1))
	checkError(ctx, getFromUrl(ctx, 5))
	slog.Info("completed")
}

func checkError(ctx context.Context, err error) {
	if err != nil {
		var apiTimeoutErr *apiTimeoutError
		if errors.As(err, &apiTimeoutErr) {
			// Was a timeout
			// Note: ctx.Err() will NOT be set in the parent context
			slog.Error("timeout err: ", slog.String("err", err.Error()), slog.Any("ctx-err", ctx.Err()))
		} else {
			// Not a timeout
			// Note: ctx.Err() will NOT be set in the parent context
			slog.Error("unexpected err: ", slog.String("err", err.Error()), slog.Any("ctx-err", ctx.Err()))
		}
	}
}

func getFromUrl(ctx context.Context, delayValue int) error {
	slog.Info("getFromUrl entry: ", slog.Int("delay-value", delayValue))
	start := time.Now()

	// UPDATED HERE
	timeoutAfter := 3 * time.Second
	// Conventionally, we can just do: `ctx, cancel := context.WithTimeoutCause(ctx, timeoutAfter, newApiTimeout(timeoutAfter))`
	ctxWithTimeout, cancel := context.WithTimeoutCause(ctx, timeoutAfter, newApiTimeout(timeoutAfter))
	defer cancel()
	// END OF UPDATED CODE

	url := fmt.Sprintf("http://0.0.0.0:80/delay/%d", delayValue)
	req, err := http.NewRequestWithContext(ctxWithTimeout, http.MethodGet, url, nil)
	if err != nil {
		return err // We don't expect this to happen
	}

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		// Note: ctxWithTimout.Err() will be set here for the child context.
		slog.Error("httpClient.Do(req) has returned error: ", slog.String("err", err.Error()), slog.Any("ctx-err", ctxWithTimeout.Err()))
		return err
	}
	defer resp.Body.Close()
	defer io.Copy(io.Discard, resp.Body)

	slog.Info("getFromUrl completed: ", slog.Int64("since-ms", int64(time.Since(start)/time.Millisecond)))

	return nil
}
