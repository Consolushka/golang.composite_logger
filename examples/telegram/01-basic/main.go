package main

import (
	"github.com/Consolushka/golang.composite_logger/pkg"
	"github.com/Consolushka/golang.composite_logger/pkg/adapters/setting"
)

func main() {
	// 1. Basic Telegram setup with Warning threshold
	composite_logger.Init(
		setting.TelegramSetting{
			Enabled:    true,
			BotKey:     "YOUR_BOT_TOKEN", // REPLACE WITH YOUR BOT TOKEN
			ChatId:     0,                // REPLACE WITH YOUR CHAT ID (int64)
			LowerLevel: composite_logger.WarningLevel,
		},
	)
	defer composite_logger.Stop()

	// This message won't be sent (level too low)
	composite_logger.Info("Hidden info", nil)

	// This message will be sent
	composite_logger.Warn("Visible warning in Telegram", nil)
}
