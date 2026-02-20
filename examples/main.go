package main

import (
	"github.com/Consolushka/golang.composite_logger/pkg"
	"github.com/Consolushka/golang.composite_logger/pkg/adapters/setting"
	"errors"
)

func main() {
	// 1. Initialize the logger with multiple adapters
	composite_logger.Init(
		setting.ConsoleSetting{
			LowerLevel: composite_logger.InfoLevel,
		},
		setting.FileSetting{
			Path:       "logs/app.log",
			LowerLevel: composite_logger.WarningLevel,
		},
		// You can also add Telegram if you have a bot key
		// setting.TelegramSetting{
		//     BotKey:     "YOUR_BOT_KEY",
		//     ChatId:     12345678,
		//     LowerLevel: composite_logger.ErrorLevel,
		// },
	)

	// 2. Simple logging
	composite_logger.Info("Application started", map[string]interface{}{
		"version": "1.0.0",
	})

	// 3. Error logging with automatic stack trace
	err := someFunctionThatFails()
	if err != nil {
		composite_logger.Error("Operation failed", map[string]interface{}{
			"error": err,
			"tags":  []string{"critical", "database"},
		})
	}

	// 4. Recovering from panics
	defer composite_logger.Recover(map[string]interface{}{
		"component": "main_loop",
	})

	// This will be caught by Recover and logged as Fatal
	// panic("something unexpected happened")
}

func someFunctionThatFails() error {
	return errors.New("database connection refused")
}
