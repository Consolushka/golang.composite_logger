package main

import (
	"time"

	"github.com/Consolushka/golang.composite_logger/pkg"
	"github.com/Consolushka/golang.composite_logger/pkg/adapters/setting"
)

func main() {
	// Initialize Composite Logger with multiple destinations
	composite_logger.Init(
		// 1. Console with JSON (default)
		setting.ConsoleSetting{
			Enabled:    true,
			LowerLevel: composite_logger.InfoLevel,
		},
		// 2. File with Text format
		setting.FileSetting{
			Enabled:         true,
			Path:            "logs/combined.log",
			IsJsonFormatter: &[]bool{false}[0],
			LowerLevel:      composite_logger.WarningLevel,
		},
		// 3. Telegram for critical alerts
		setting.TelegramSetting{
			Enabled:    true,
			BotKey:     "YOUR_BOT_TOKEN",
			ChatId:     0,
			Timeout:    10 * time.Second,
			LowerLevel: composite_logger.ErrorLevel,
		},
	)
	defer composite_logger.Stop()

	// This goes only to Console
	composite_logger.Info("Normal operation", nil)

	// This goes to Console and File
	composite_logger.Warn("High memory usage", map[string]interface{}{"usage": "85%"})

	// This goes to Console, File, AND Telegram
	composite_logger.Error("Database connection lost!", nil)
}
