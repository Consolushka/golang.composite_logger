package main

import (
	"github.com/Consolushka/golang.composite_logger/pkg"
	"github.com/Consolushka/golang.composite_logger/pkg/adapters/setting"
)

func main() {
	// TIP: Always call composite_logger.Stop() with defer at the beginning of your application.
	defer composite_logger.Stop()

	// Initialize with multiple loggers
	composite_logger.Init(
		setting.ConsoleSetting{
			Enabled:    true,
			LowerLevel: composite_logger.InfoLevel,
		},
		setting.FileSetting{
			Enabled:    true,
			Path:       "logs/composite_basic.log",
			LowerLevel: composite_logger.InfoLevel,
		},
	)

	// Traditional logging without context
	composite_logger.Info("Application started", map[string]interface{}{"version": "1.0.0"})
	composite_logger.Warn("Resource usage high", map[string]interface{}{"cpu": "90%"})
}
