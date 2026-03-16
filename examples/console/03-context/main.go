package main

import (
	"context"

	"github.com/Consolushka/golang.composite_logger/pkg"
	"github.com/Consolushka/golang.composite_logger/pkg/adapters/setting"
)

func main() {
	// TIP: Use InfoContext, WarnContext, etc., to automatically include trace IDs or request IDs in your logs if your adapters (like Logrus) are configured with appropriate hooks.
	defer composite_logger.Stop()

	// 1. Initialize console logger

	composite_logger.Init(
		setting.ConsoleSetting{
			Enabled:    true,
			LowerLevel: composite_logger.InfoLevel,
		},
	)

	// 2. Create context with trace ID
	ctx := context.WithValue(context.Background(), "trace_id", "console-trace-789")

	// 3. Log with context
	// Note: Logrus hooks can be used to extract "trace_id" and print it automatically
	composite_logger.InfoContext(ctx, "Logging with context to console", map[string]interface{}{
		"component": "auth",
	})
}
