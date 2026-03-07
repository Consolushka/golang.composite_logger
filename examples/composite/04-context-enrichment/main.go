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

	// 1. Register context keys that should be automatically extracted
	composite_logger.SetContextKeys("trace_id", "request_id")

	// 2. Prepare context with values
	ctx := context.WithValue(context.Background(), "trace_id", "abc-123")
	ctx = context.WithValue(ctx, "request_id", "req-789")

	// 3. Log using InfoContext - keys will be added to the output AUTOMATICALLY
	composite_logger.InfoContext(ctx, "Handling user request", map[string]interface{}{"user_id": 42})
	
	// You can also use WithContext pattern with enrichment
	log := composite_logger.WithContext(ctx)
	log.Info("Steps completed", nil) // trace_id and request_id are still included!
}
