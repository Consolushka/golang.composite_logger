package main

import (
	"github.com/Consolushka/golang.composite_logger/pkg"
	"github.com/Consolushka/golang.composite_logger/pkg/adapters/setting"
)

func main() {
	// Initialize File logger with custom rotation settings
	composite_logger.Init(
		setting.FileSetting{
			Enabled:    true,
			Path:       "logs/rotated.log",
			LowerLevel: composite_logger.InfoLevel,
			// Rotation settings
			MaxSize:    10,   // Rotate when file reaches 10 Megabytes
			MaxBackups: 5,    // Keep up to 5 old log files
			MaxAge:     7,    // Retain old logs for 7 days
			Compress:   true, // Compress (gzip) old log files
		},
	)
	defer composite_logger.Stop()

	composite_logger.Info("Logging with rotation enabled", map[string]interface{}{
		"max_size_mb": 10,
		"backups":     5,
	})
}
