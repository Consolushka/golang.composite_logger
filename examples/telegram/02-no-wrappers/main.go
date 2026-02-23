package main

import (
	"github.com/Consolushka/golang.composite_logger/pkg"
	"github.com/Consolushka/golang.composite_logger/pkg/adapters/setting"
)

func main() {
	// 2. Disable level decoration (no emojis like ⚠️⚠️ around titles)
	composite_logger.Init(
		setting.TelegramSetting{
			Enabled:              true,
			BotKey:               "YOUR_BOT_TOKEN", // REPLACE WITH YOUR BOT TOKEN
			ChatId:               0,                // REPLACE WITH YOUR CHAT ID (int64)
			LowerLevel:           composite_logger.WarningLevel,
			UseLevelTitleWrapper: &[]bool{false}[0], // Explicitly disable wrappers
		},
	)
	defer composite_logger.Stop()

	composite_logger.Warn("Warning message without emojis in title", nil)
}
