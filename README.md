# Composite Logger

A flexible, multi-destination logging library for Go, built with **Hexagonal Architecture (Ports & Adapters)**. It allows you to broadcast logs to the Console, Files, and Telegram simultaneously with a unified API.

## Features

- 🚀 **Asynchronous Engine**: Non-blocking logging using a background worker and buffered channels.
- 🏗 **Clean Architecture**: Decoupled core logic from specific implementations using Ports & Adapters.
- 🌐 **Context Support**: Fully supports `context.Context` for integration with tracing systems (OpenTelemetry, etc.) and cancellation.
- 📄 **Structured Logging**: Powered by [Logrus](https://github.com/sirupsen/logrus) with JSON and Text support.
- 🤖 **Telegram Integration**: Send formatted alerts to Telegram with custom emojis, titles, and configurable timeouts.
- 📦 **Log Rotation**: Built-in log rotation for file adapter using [Lumberjack](https://github.com/natefinch/lumberjack).
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
	"context"
	"github.com/Consolushka/golang.composite_logger/pkg"
	"github.com/Consolushka/golang.composite_logger/pkg/adapters/setting"
)

func main() {
	// Initialize with Console and File loggers
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
	)

	// ALWAYS call Stop() at the end to flush the async queue
	defer composite_logger.Stop()

	// Simple logging
	composite_logger.Info("Hello, World!", map[string]interface{}{
		"user_id": 42,
	})

	// Context-aware logging (useful for tracing/request IDs)
	ctx := context.WithValue(context.Background(), "trace_id", "abc-123")
	composite_logger.InfoContext(ctx, "Operation started", map[string]interface{}{
		"action": "sync",
	})
}
```

## Advanced Configuration

### JSON Formatting
Both Console and File adapters support JSON formatting (enabled by default).

```go
setting.ConsoleSetting{
    Enabled:         true,
    IsJsonFormatter: &[]bool{false}[0], // Explicitly disable JSON to use Text
}
```

### Telegram Adapter
Highly customizable with timeouts and level decorators.

```go
composite_logger.Init(
    setting.TelegramSetting{
        Enabled:              true,
        BotKey:               "YOUR_BOT_TOKEN",
        ChatId:               12345678,
        Timeout:              10 * time.Second, // Network timeout
        LowerLevel:           composite_logger.ErrorLevel,
        UseLevelTitleWrapper: &[]bool{true}[0],
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

Use `Recover` or `RecoverContext` in your `defer` blocks to ensure panics are captured and logged with full stack traces.

```go
func someDangerousOperation(ctx context.Context) {
    // Log panic with trace information from context
    defer composite_logger.RecoverContext(ctx, map[string]interface{}{
        "op": "database_sync",
    })

    panic("unexpected failure") // This will be logged as a FATAL error
}
```

## Examples

The [examples/](./examples) directory contains a structured set of lessons to help you get started:

- **Console**: [Text format](./examples/console/01-text), [JSON format](./examples/console/02-json), [Context support](./examples/console/03-context)
- **File**: [Text format](./examples/file/01-text), [JSON format](./examples/file/02-json), [Rotation](./examples/file/03-rotation), [Context support](./examples/file/04-context)
- **Telegram**: [Basic](./examples/telegram/01-basic), [Decorations](./examples/telegram/02-no-wrappers), [Custom Emojis](./examples/telegram/03-custom-wrappers), [Custom Titles](./examples/telegram/04-custom-titles), [Timeouts](./examples/telegram/05-timeout), [Context Cancellation](./examples/telegram/06-context)
- **Advanced**: [Composite (Basic)](./examples/composite/01-basic), [Composite (Context methods)](./examples/composite/02-context-methods), [Composite (WithContext pattern)](./examples/composite/03-with-context), [Composite (Context enrichment)](./examples/composite/04-context-enrichment), [Custom Adapter implementation](./examples/custom-adapter)

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
