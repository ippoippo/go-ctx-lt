package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func poller(ctx context.Context) {
	slog.Info("poller entry: ")
	start := time.Now()

	for {
		select {
		case <-ctx.Done():
			// Here we could do some graceful cleanup
			slog.Info("poller received cancellation signal, stopping execution: ",
				slog.Int64("since-ms", int64(time.Since(start)/time.Millisecond)),
				slog.Any("ctx-err", ctx.Err()))
			return
		default:
			err := getFromUrl(ctx, 1)
			if err != nil {
				// Note: ctx.Err() will NOT be set in the parent context
				slog.Error("getFromUrl(ctx, 1) has returned error: ", slog.String("err", err.Error()), slog.Any("ctx-err", ctx.Err()))
			} else {
				slog.Info("doing some more work, such as writing to DB: ", slog.Any("ctx.Done()", ctx.Done()), slog.Any("ctx-err", ctx.Err()))
			}
			time.Sleep(1 * time.Second) // Let's slow down the amount of logging
		}
	}
}

func getFromUrl(ctx context.Context, delayValue int) error {
	slog.Info("getFromUrl entry: ", slog.Int("delay-value", delayValue))
	start := time.Now()

	// Conventionally, we can just do: `ctx, cancel := context.WithTimeout(ctx, 3*time.Second)`
	ctxWithTimout, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	url := fmt.Sprintf("http://0.0.0.0:80/delay/%d", delayValue)
	req, err := http.NewRequestWithContext(ctxWithTimout, http.MethodGet, url, nil)
	if err != nil {
		return err // We don't expect this to happen
	}

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		// Note: ctxWithTimout.Err() will be set here for the child context.
		slog.Error("httpClient.Do(req) has returned error: ", slog.String("err", err.Error()), slog.Any("ctx-err", ctxWithTimout.Err()))
		return err
	}
	defer resp.Body.Close()
	defer io.Copy(io.Discard, resp.Body)

	slog.Info("getFromUrl completed: ", slog.Int64("since-ms", int64(time.Since(start)/time.Millisecond)))

	return nil
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go poller(ctx)

	// Wait for interrupt signal
	<-ctx.Done()

	// Create a new context, and wait for 10 seconds while we cleanup
	closingCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	time.Sleep(3 * time.Second) // We are doing clean up for example

	slog.Info("completed", slog.Any("closingCtx.", closingCtx.Err())) // Assuming we cleanup in time, we log this
}
