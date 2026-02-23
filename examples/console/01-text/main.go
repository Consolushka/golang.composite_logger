package main

import (
	"github.com/Consolushka/golang.composite_logger/pkg"
	"github.com/Consolushka/golang.composite_logger/pkg/adapters/setting"
)

func main() {
	// Initialize Console logger with Text format and Warning threshold
	composite_logger.Init(
		setting.ConsoleSetting{
			Enabled:         true,
			IsJsonFormatter: &[]bool{false}[0], // Text format
			LowerLevel:      composite_logger.WarningLevel,
		},
	)
	defer composite_logger.Stop()

	// This WILL NOT be visible (level too low)
	composite_logger.Info("This info message is hidden", nil)

	// These WILL BE visible
	composite_logger.Warn("Warning: resource limit reached", map[string]interface{}{
		"usage": "95%",
	})

	composite_logger.Error("Error: operation failed", nil)
}
