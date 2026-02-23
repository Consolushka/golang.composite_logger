package main

import (
	"time"

	"github.com/Consolushka/golang.composite_logger/pkg"
	"github.com/Consolushka/golang.composite_logger/pkg/adapters/setting"
)

func main() {
	// 5. Configure network timeout for Telegram API requests
	composite_logger.Init(
		setting.TelegramSetting{
			Enabled:    true,
			BotKey:     "YOUR_BOT_TOKEN", // REPLACE WITH YOUR BOT TOKEN
			ChatId:     0,                // REPLACE WITH YOUR CHAT ID (int64)
			Timeout:    5 * time.Second,  // Wait maximum 5 seconds
			LowerLevel: composite_logger.WarningLevel,
		},
	)
	defer composite_logger.Stop()

	composite_logger.Warn("Warning with 5s API timeout setting", nil)
}
