package hook

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

// OpenTelemetryHook is an adapter that extracts trace and span IDs from the context
// and adds them to the log fields. This enables correlation between logs and traces.
type OpenTelemetryHook struct{}

// Fire is called for each log entry. It extracts the span context from the provided
// context and adds "trace_id" and "span_id" to the fields if they are present.
func (h OpenTelemetryHook) Fire(ctx context.Context, level string, message string, fields map[string]interface{}) error {
	span := trace.SpanContextFromContext(ctx)
	if span.HasTraceID() {
		fields["trace_id"] = span.TraceID().String()
	}
	if span.HasSpanID() {
		fields["span_id"] = span.SpanID().String()
	}
	return nil
}
