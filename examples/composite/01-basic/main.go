package main

import (
	"github.com/Consolushka/golang.composite_logger/pkg"
	"github.com/Consolushka/golang.composite_logger/pkg/adapters/setting"
)

func main() {
	// CRITICAL: When using InitAsync, you MUST call Stop() to ensure all queued logs are flushed.
	defer composite_logger.Stop()

	// Initialize in ASYNCHRONOUS mode (non-blocking).
	// This is recommended for high-performance applications.
	composite_logger.InitAsync(
		setting.ConsoleSetting{
			Enabled:    true,
			LowerLevel: composite_logger.InfoLevel,
		},
		setting.FileSetting{
			Enabled:    true,
			Path:       "logs/composite_basic_async.log",
			LowerLevel: composite_logger.InfoLevel,
		},
	)

	// Logging here is non-blocking (logs are sent to a background worker).
	composite_logger.Info("Application started (async)", map[string]interface{}{"version": "1.0.0"})
	composite_logger.Warn("Resource usage high", map[string]interface{}{"cpu": "90%"})
}
