package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Consolushka/golang.composite_logger/pkg"
	"github.com/Consolushka/golang.composite_logger/pkg/adapters/setting"
)

// CustomHook is an example of a simple hook that prints to console.
// This could be replaced with an OpenTelemetry hook that creates span events.
type CustomHook struct{}

func (h CustomHook) Fire(ctx context.Context, level string, message string, fields map[string]interface{}) error {
	fmt.Printf("[HOOK] Level: %s, Message: %s, Fields: %v\n", level, message, fields)
	return nil
}

func main() {
	// Initialize composite logger with console output
	composite_logger.Init(
		setting.ConsoleSetting{
			Enabled:    true,
			LowerLevel: composite_logger.InfoLevel,
		},
	)
	defer composite_logger.Stop()

	// Add our custom hook
	composite_logger.AddHook(CustomHook{})

	// Log some messages
	composite_logger.Info("Hello with hook", map[string]interface{}{"foo": "bar"})

	// Use context to show it's passed to the hook
	ctx := context.WithValue(context.Background(), "request_id", "12345")
	composite_logger.InfoContext(ctx, "Hello with context hook", map[string]interface{}{"user": "admin"})

	// Give some time for async processing
	time.Sleep(100 * time.Millisecond)
}
