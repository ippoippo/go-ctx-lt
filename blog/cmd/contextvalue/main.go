package main

import (
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/ippoippo/go-ctx-lt/blog/cmd/contextvalue/tracing"
)

func traceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := tracing.ContextWithTraceId(r.Context(), uuid.New().String())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func messagesHandler(w http.ResponseWriter, r *http.Request) {
	traceID := tracing.TraceIdFromContext(r.Context())

	slog.Info("messages handler", slog.String("trace-id", traceID))

	// Do other business logic
	// result, err := DoThing(r.Context(), ...)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Trace-ID logged"))
}

func main() {
	mux := http.NewServeMux()

	mux.Handle("/messages", traceMiddleware(http.HandlerFunc(messagesHandler)))

	slog.Info("Server is running on http://0.0.0.0:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		slog.Error("error from server", slog.Any("err", err))
	}
}
