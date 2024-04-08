package logger

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

func getTraceID(ctx context.Context) string {
	spanContext := trace.SpanContextFromContext(ctx)
	if !spanContext.HasTraceID() {
		return ""
	}
	return spanContext.TraceID().String()
}

func getSpanID(ctx context.Context) string {
	spanContext := trace.SpanContextFromContext(ctx)
	if !spanContext.HasSpanID() {
		return ""
	}
	return spanContext.SpanID().String()
}
