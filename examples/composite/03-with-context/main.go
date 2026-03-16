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

	ctx := context.WithValue(context.Background(), "trace_id", "trace-789")

	// Bind context once to create a LoggingContext
	log := composite_logger.WithContext(ctx)

	// Now use short methods - context is preserved automatically
	log.Info("Operation started", nil)
	log.Warn("High latency detected", map[string]interface{}{"ms": 500})
	log.Info("Operation completed", nil)
}
