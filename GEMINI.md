# Composite Logger (Ports & Adapters)

A flexible Go logging library implementing the composite pattern and hexagonal architecture to broadcast logs across multiple destinations.

## Project Overview

- **Purpose**: Unified logging interface with pluggable adapters.
- **Architecture**:
    - `pkg/ports/`: Core interfaces (`Logger`, `LoggerSetting`).
    - `pkg/adapters/setting/`: Adapter implementations for configuration (Console, File, Telegram).
    - `internal/adapters/logger/`: Concrete logger implementations (hidden from public API).
    - `pkg/composite_logger.go`: The central hub that manages multiple loggers and an asynchronous worker.

## Initialization

The library uses a variadic `Init` function that accepts multiple settings. It starts a background worker to process logs asynchronously.

```go
import (
    "github.com/Consolushka/golang.composite_logger/pkg"
    "github.com/Consolushka/golang.composite_logger/pkg/adapters/setting"
)

func main() {
    // Initialize with desired loggers
    composite_logger.Init(
        setting.ConsoleSetting{
            Enabled:    true, 
            LowerLevel: composite_logger.InfoLevel,
        },
        setting.FileSetting{
            Enabled:    true,
            Path:       "logs/app.log",
            LowerLevel: composite_logger.WarningLevel,
        },
        setting.TelegramSetting{
            Enabled:              true,
            BotKey:               "YOUR_BOT_KEY",
            ChatId:               12345678,
            LowerLevel:           composite_logger.ErrorLevel,
            UseLevelTitleWrapper: &[]bool{true}[0], 
        },
    )
    
    // CRITICAL: Always call Stop() to flush the async queue
    defer composite_logger.Stop()
    
    // Use the global logger
    composite_logger.Info("App started", nil)
    
    // Capture panics
    defer composite_logger.Recover(map[string]interface{}{"context": "main"})
}
```

## Available Adapters (Settings)

### Console
Writes to standard output using Logrus.
- **Settings**: 
    - `Enabled`: (bool)
    - `IsJsonFormatter`: (*bool) Default is true.
    - `LowerLevel`: `composite_logger.Level`

### File
Writes to a file and standard output using Logrus.
- **Settings**: 
    - `Enabled`: (bool)
    - `IsJsonFormatter`: (*bool) Default is true.
    - `Path`: string
    - `MaxSize`: (int) Maximum size in megabytes before rotation (default: 5).
    - `MaxBackups`: (int) Maximum number of old log files to retain (default: 3).
    - `MaxAge`: (int) Maximum number of days to retain old log files (default: 28).
    - `Compress`: (bool) Whether to compress old log files (default: true).
    - `LowerLevel`: `composite_logger.Level`

### Telegram
Sends formatted MarkdownV2 messages to a Telegram chat. Handles errors with console fallback.
- **Settings**: 
    - `Enabled`: (bool)
    - `BotKey`: Telegram bot token.
    - `ChatId`: ID of the chat/user to receive logs.
    - `Timeout`: `time.Duration` for API requests.
    - `LowerLevel`: Minimum level to send.
    - `UseLevelTitleWrapper`: (*bool) Wrap level name with symbols (e.g., 🚨 ERROR 🚨).
    - `LevelWrappers`: (map) Custom wrappers per level.
    - `LevelTitles`: (map) Custom display names for levels.

## Error Handling
The `CompositeLogger` automatically captures stack traces when `Error` or `Fatal` methods are called. Use `composite_logger.Recover(ctx)` in defer statements to safely catch and log panics. Stack traces are cleaned to exclude internal library frames.
