package main

import (
	"github.com/Consolushka/golang.composite_logger/pkg"
	"github.com/Consolushka/golang.composite_logger/pkg/adapters/setting"
)

func main() {
	// Initialize Console logger with JSON format (default) and Warning threshold
	composite_logger.Init(
		setting.ConsoleSetting{
			Enabled:    true,
			LowerLevel: composite_logger.WarningLevel,
		},
	)
	defer composite_logger.Stop()

	// This WILL NOT be visible (level too low)
	composite_logger.Info("This info message is hidden", nil)

	// These WILL BE visible in JSON format
	composite_logger.Warn("Warning: resource limit reached", map[string]interface{}{
		"usage": "95%",
	})

	composite_logger.Error("Error: operation failed", nil)
}
