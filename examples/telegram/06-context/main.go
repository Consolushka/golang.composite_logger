package main

import (
	"context"
	"fmt"

	"github.com/Consolushka/golang.composite_logger/pkg"
	"github.com/Consolushka/golang.composite_logger/pkg/adapters/setting"
)

func main() {
	// The context passed to *Context methods is a source of values (trace id,
	// request id, ...) — cancellation never suppresses a log. A log written for
	// a timed-out or cancelled request is exactly the one you want to keep.
	defer composite_logger.Stop()

	// 1. Initialize telegram logger (use real keys in your project)
	composite_logger.Init(
		setting.TelegramSetting{
			Enabled:    true,
			BotKey:     "YOUR_BOT_TOKEN",
			ChatId:     12345678,
			LowerLevel: composite_logger.InfoLevel,
		},
	)

	// 2. Even a cancelled context does not stop the log from being delivered
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancelling immediately

	// 3. The message is still sent to Telegram
	composite_logger.InfoContext(ctx, "This log is delivered even though the context is cancelled", nil)

	fmt.Println("Check your Telegram chat. The message above was delivered despite the cancelled context.")
}
