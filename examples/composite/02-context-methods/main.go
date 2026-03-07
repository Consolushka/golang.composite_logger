package main

import (
	"context"

	"github.com/Consolushka/golang.composite_logger/pkg"
	"github.com/Consolushka/golang.composite_logger/pkg/adapters/setting"
)

func main() {
	// TIP: Always call composite_logger.Stop() with defer at the beginning of your application.
	defer composite_logger.Stop()

	composite_logger.Init(
		setting.ConsoleSetting{Enabled: true, LowerLevel: composite_logger.InfoLevel},
	)

	// Create context with metadata (e.g. for tracing)
	ctx := context.WithValue(context.Background(), "request_id", "req-123")

	// Use Context-suffixed methods
	composite_logger.InfoContext(ctx, "Handling request", map[string]interface{}{"path": "/api/user"})
	composite_logger.ErrorContext(ctx, "Request failed", map[string]interface{}{"error": "timeout"})
}
