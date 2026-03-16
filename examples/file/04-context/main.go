package main

import (
	"context"
	"github.com/Consolushka/golang.composite_logger/pkg"
	"github.com/Consolushka/golang.composite_logger/pkg/adapters/setting"
)

func main() {
	// 1. Initialize file logger
	composite_logger.Init(
		setting.FileSetting{
			Enabled:    true,
			Path:       "logs/context_example.log",
			LowerLevel: composite_logger.InfoLevel,
		},
	)
	defer composite_logger.Stop()

	// 2. Create context with trace ID
	ctx := context.WithValue(context.Background(), "trace_id", "file-trace-123")

	// 3. Log with context to file
	composite_logger.InfoContext(ctx, "Logging with context to file", map[string]interface{}{
		"user": "alice",
	})
}
