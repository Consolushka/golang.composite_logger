package main

import (
	"context"
	"fmt"

	"github.com/Consolushka/golang.composite_logger/pkg"
	"github.com/Consolushka/golang.composite_logger/pkg/adapters/setting"
)

func main() {
	// TIP: Providing a context with cancellation (or timeout) is highly recommended for network-based adapters like Telegram. It prevents your application from making useless API calls if the operation has already been aborted.
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

	// 2. Demonstrate context cancellation
	// Creating a context that will be canceled before log dispatch
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancelling immediately

	// 3. Try to log with a canceled context
	// In the logs, you should see nothing, and no network request will be sent to Telegram.
	composite_logger.InfoContext(ctx, "This log will be skipped due to context cancellation", nil)

	fmt.Println("Check your logs or console output. The message above was skipped before reaching the network.")
}
