# Composite Logger (Ports & Adapters)

A flexible Go logging library implementing the composite pattern and hexagonal architecture to broadcast logs across multiple destinations.

## Project Overview

- **Purpose**: Unified logging interface with pluggable adapters.
- **Architecture**:
    - `pkg/ports/`: Core interfaces (`Logger`, `LoggerSetting`).
    - `pkg/adapters/setting/`: Adapter implementations for configuration (Console, File, Telegram).
    - `internal/adapters/logger/`: Concrete logger implementations (hidden from public API).
    - `pkg/composite_logger.go`: The central hub that manages multiple loggers.

## Initialization

The library uses a variadic `Init` function that accepts multiple settings:

```go
import (
    "composite_logger/pkg"
    "composite_logger/pkg/adapters/setting"
)

func main() {
    // Initialize with desired loggers
    composite_logger.Init(
        setting.ConsoleSetting{LowerLevel: composite_logger.InfoLevel},
        setting.FileSetting{
            Path:       "logs/app.log",
            LowerLevel: composite_logger.WarningLevel,
        },
        setting.TelegramSetting{
            BotKey:               "YOUR_BOT_KEY",
            ChatId:               12345678,
            LowerLevel:           composite_logger.ErrorLevel,
            UseLevelTitleWrapper: true, // Use emojis/custom wrappers
        },
    )
    
    // Use the global logger
    composite_logger.Info("App started", nil)
    
    // Capture panics
    defer composite_logger.Recover(map[string]interface{}{"context": "main"})
}
```

## Available Adapters (Settings)

### Console
Writes to standard output using Logrus.
- **Setting**: `setting.ConsoleSetting{LowerLevel: composite_logger.Level}`

### File
Writes to a file and standard output using Logrus.
- **Setting**: `setting.FileSetting{Path: string, LowerLevel: composite_logger.Level}`

### Telegram
Sends formatted MarkdownV2 messages to a Telegram chat.
- **Setting**: 
    - `BotKey`: Telegram bot token.
    - `ChatId`: ID of the chat/user to receive logs.
    - `LowerLevel`: Minimum level to send.
    - `UseLevelTitleWrapper`: (bool) Wrap level name with symbols (e.g., 🚨 ERROR 🚨).
    - `LevelWrappers`: (map) Custom wrappers per level (defaults to emojis like ℹ️ℹ️, 🚨🚨).
    - `LevelTitles`: (map) Custom display names for levels (defaults to uppercase level name).

## Error Handling
The `CompositeLogger` automatically captures stack traces when `Error` or `Fatal` methods are called. Use `composite_logger.Recover(ctx)` in defer statements to safely catch and log panics. Stack traces are cleaned to exclude internal library frames.
