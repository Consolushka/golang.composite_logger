package main

import (
	"github.com/Consolushka/golang.composite_logger/pkg"
	"github.com/Consolushka/golang.composite_logger/pkg/adapters/setting"
)

func main() {
	// 3. Use custom emojis/symbols for different log levels
	composite_logger.Init(
		setting.TelegramSetting{
			Enabled:    true,
			BotKey:     "YOUR_BOT_TOKEN", // REPLACE WITH YOUR BOT TOKEN
			ChatId:     0,                // REPLACE WITH YOUR CHAT ID (int64)
			LowerLevel: composite_logger.WarningLevel,
			LevelWrappers: map[composite_logger.Level]string{
				composite_logger.WarningLevel: "🟡",
				composite_logger.ErrorLevel:   "🔴",
				composite_logger.FatalLevel:   "💀",
			},
		},
	)
	defer composite_logger.Stop()

	composite_logger.Warn("Custom yellow circle warning", nil)
	composite_logger.Error("Custom red circle error", nil)
}
