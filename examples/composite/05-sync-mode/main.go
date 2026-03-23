package main

import (
	"github.com/Consolushka/golang.composite_logger/pkg"
	"github.com/Consolushka/golang.composite_logger/pkg/adapters/setting"
)

func main() {
	// Stop is still good practice to call, although it's a no-op in synchronous mode.
	defer composite_logger.Stop()

	// Initialize in SYNCHRONOUS mode (blocking).
	// This ensures each log message is fully processed by all adapters
	// before the logging function returns. Ideal for CLI tools.
	composite_logger.Init(
		setting.ConsoleSetting{
			Enabled:    true,
			LowerLevel: composite_logger.InfoLevel,
		},
		setting.FileSetting{
			Enabled:    true,
			Path:       "logs/composite_sync.log",
			LowerLevel: composite_logger.InfoLevel,
		},
	)

	// Logging here is blocking (waits for console and file write to complete).
	composite_logger.Info("Synchronous log entry", map[string]interface{}{"mode": "sync"})
	composite_logger.Error("Critical failure occurred", map[string]interface{}{"disk": "full"})
}
