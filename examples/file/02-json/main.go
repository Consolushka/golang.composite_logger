package main

import (
	"github.com/Consolushka/golang.composite_logger/pkg"
	"github.com/Consolushka/golang.composite_logger/pkg/adapters/setting"
)

func main() {
	// Initialize File logger with JSON format and Warning threshold
	composite_logger.Init(
		setting.FileSetting{
			Enabled:    true,
			Path:       "logs/json.log",
			LowerLevel: composite_logger.WarningLevel,
		},
	)
	defer composite_logger.Stop()

	composite_logger.Info("Hidden info", nil)
	composite_logger.Warn("Warning in file", map[string]interface{}{"file": "demo"})
	composite_logger.Error("Error in file", nil)
}
