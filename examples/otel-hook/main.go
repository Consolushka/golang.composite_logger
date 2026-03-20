package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Consolushka/golang.composite_logger/pkg"
	"github.com/Consolushka/golang.composite_logger/pkg/adapters/hook"
	"github.com/Consolushka/golang.composite_logger/pkg/adapters/setting"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

func initTracer() *sdktrace.TracerProvider {
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("example-service"),
		)),
	)
	otel.SetTracerProvider(tp)
	return tp
}

func main() {
	// 1. Setup dummy OTel Tracer
	tp := initTracer()
	defer func() { _ = tp.Shutdown(context.Background()) }()
	tracer := tp.Tracer("example-tracer")

	// 2. Init Logger with Console Adapter
	composite_logger.Init(
		setting.ConsoleSetting{
			Enabled:    true,
			LowerLevel: composite_logger.InfoLevel,
		},
	)
	defer composite_logger.Stop()

	// 3. Add OpenTelemetry Hook (Adapter)
	composite_logger.AddHook(hook.OpenTelemetryHook{})

	// 4. Create a span and log within its context
	ctx, span := tracer.Start(context.Background(), "main-operation")
	defer span.End()

	fmt.Printf("Current TraceID: %s\n", span.SpanContext().TraceID())

	// This log will automatically include trace_id and span_id thanks to the hook
	composite_logger.InfoContext(ctx, "Operation started with OTel tracing", map[string]interface{}{
		"custom_field": "val",
	})

	// Give some time for async processing
	time.Sleep(100 * time.Millisecond)
}
