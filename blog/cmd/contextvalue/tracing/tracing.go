package tracing

import (
	"context"
)

type traceIdContextKeyType struct{}

var (
	traceIdContextKey = traceIdContextKeyType{}
)

func ContextWithTraceId(ctx context.Context, traceId string) context.Context {
	return context.WithValue(ctx, traceIdContextKey, traceId)
}

func TraceIdFromContext(ctx context.Context) string {
	if traceId := ctx.Value(traceIdContextKey); traceId != nil {
		value, ok := traceId.(string)
		if !ok {
			return ""
		}
		return value
	}
	return ""
}
