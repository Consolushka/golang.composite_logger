package main

import (
	"github.com/Consolushka/golang.composite_logger/pkg"
	"github.com/Consolushka/golang.composite_logger/pkg/adapters/setting"
)

func main() {
	// 4. Override default level names (e.g., change ERROR to ALARM)
	composite_logger.Init(
		setting.TelegramSetting{
			Enabled:    true,
			BotKey:     "YOUR_BOT_TOKEN", // REPLACE WITH YOUR BOT TOKEN
			ChatId:     0,                // REPLACE WITH YOUR CHAT ID (int64)
			LowerLevel: composite_logger.WarningLevel,
			LevelTitles: map[composite_logger.Level]string{
				composite_logger.WarningLevel: "ATTENTION",
				composite_logger.ErrorLevel:   "ALARM",
			},
		},
	)
	defer composite_logger.Stop()

	composite_logger.Error("This will show as [ALARM] in Telegram", nil)
}
