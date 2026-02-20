# Composite Logger

A flexible, multi-destination logging library for Go, built with **Hexagonal Architecture (Ports & Adapters)**. It allows you to broadcast logs to the Console, Files, and Telegram simultaneously with a unified API.

## Features

- 🚀 **Multi-destination**: Log to multiple sinks at once.
- 🏗 **Clean Architecture**: Decoupled core logic from specific implementations.
- 📄 **Structured Logging**: Powered by [Logrus](https://github.com/sirupsen/logrus).
- 🤖 **Telegram Integration**: Send formatted alerts to Telegram with custom emojis and titles.
- 🔍 **Auto Stack Traces**: Automatically captures and cleans stack traces for `Error` and `Fatal` levels.
- 🛡 **Panic Recovery**: Catch and log panics as Fatal errors.

## Installation

```bash
go get github.com/Consolushka/golang.composite_logger
```

## Quick Start

```go
package main

import (
	"github.com/Consolushka/golang.composite_logger/pkg"
	"github.com/Consolushka/golang.composite_logger/pkg/adapters/setting"
)

func main() {
	// Initialize with Console and File loggers
	composite_logger.Init(
		setting.ConsoleSetting{
			LowerLevel: composite_logger.InfoLevel,
		},
		setting.FileSetting{
			Path:       "logs/app.log",
			LowerLevel: composite_logger.WarningLevel,
		},
	)

	composite_logger.Info("Hello, World!", map[string]interface{}{
		"user_id": 42,
	})
}
```

## Advanced Configuration

### Telegram Adapter

The Telegram adapter is highly customizable, allowing you to wrap level names in emojis or symbols.

```go
composite_logger.Init(
    setting.TelegramSetting{
        BotKey:               "YOUR_BOT_TOKEN",
        ChatId:               12345678,
        LowerLevel:           composite_logger.ErrorLevel,
        UseLevelTitleWrapper: true,
        // Optional: override default emojis (ℹ️ℹ️, 🚨🚨, etc.)
        LevelWrappers: map[composite_logger.Level]string{
            composite_logger.FatalLevel: "💀",
        },
        // Optional: override level names
        LevelTitles: map[composite_logger.Level]string{
            composite_logger.ErrorLevel: "ALARM",
        },
    },
)
```

### Panic Recovery

Use `Recover` in your `defer` blocks to ensure panics are captured and logged with full stack traces.

```go
func someDangerousOperation() {
    defer composite_logger.Recover(map[string]interface{}{
        "op": "database_sync",
    })

    panic("unexpected failure") // This will be logged as a FATAL error
}
```

## Project Structure

- `pkg/`: Public API and core types.
- `pkg/ports/`: Interfaces for loggers and settings.
- `pkg/adapters/setting/`: Concrete settings used for initialization.
- `internal/`: Private implementations and internal logic.

## Log Levels

1. `InfoLevel`
2. `WarningLevel`
3. `ErrorLevel`
4. `FatalLevel`

## License

MIT
