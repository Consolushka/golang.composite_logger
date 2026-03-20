# Composite Logger

A flexible, asynchronous logging library for Go, built with **Hexagonal Architecture (Ports & Adapters)**. It allows you to broadcast logs to multiple destinations simultaneously using a unified, non-blocking engine.

## Key Features

- 🚀 **Asynchronous Core**: Logging operations are non-blocking, handled by a background worker with a buffered queue.
- 🏗 **Hexagonal Architecture**: Core logic is decoupled from output implementations.
- 🌐 **First-class Context Support**: Native support for `context.Context` to handle tracing, request IDs, and cancellation.
- 🔗 **Composite Pattern**: Broadcast a single log message to N destinations with different filters and formats.
- 🔍 **Diagnostic-ready**: Automatic stack trace capture for errors and built-in panic recovery.

---

## Supported Adapters (Outputs)

The library provides pluggable adapters to send logs where you need them:

- **Console**: Standard output with support for JSON or human-readable text.
- **File**: Persistent storage with built-in log rotation (size, age, backups).
- **Telegram**: Formatted MarkdownV2 alerts for critical errors with customizable decorators.

*Want more? You can easily implement your own adapter by satisfying the `ports.Logger` interface.*

---

## Hooks & Middleware

Hooks allow you to process or enrich log entries before they are broadcast to adapters.

- **OpenTelemetry Hook**: Automatically extracts `trace_id` and `span_id` from the context.
- **Custom Hooks**: Implement the `ports.Hook` interface to add custom metadata, audit logs, or filtering logic.

---

## Installation

```bash
go get github.com/Consolushka/golang.composite_logger
```

---

## Full-featured Example

This example demonstrates a production-ready setup: multiple adapters, OTel correlation, and panic safety.

```go
package main

import (
    "context"
    "time"

    "github.com/Consolushka/golang.composite_logger/pkg"
    "github.com/Consolushka/golang.composite_logger/pkg/adapters/hook"
    "github.com/Consolushka/golang.composite_logger/pkg/adapters/setting"
)

func main() {
    // 1. Initialize with multiple destinations
    composite_logger.Init(
        // Console for development
        setting.ConsoleSetting{
            Enabled:    true,
            LowerLevel: composite_logger.InfoLevel,
        },
        // File for persistence with rotation
        setting.FileSetting{
            Enabled:    true,
            Path:       "logs/app.log",
            MaxSize:    10, // MB
            LowerLevel: composite_logger.InfoLevel,
        },
        // Telegram for critical alerts
        setting.TelegramSetting{
            Enabled:    true,
            BotKey:     "YOUR_BOT_TOKEN",
            ChatId:     12345678,
            LowerLevel: composite_logger.ErrorLevel,
        },
    )

    // CRITICAL: Ensure all logs are flushed before exit
    defer composite_logger.Stop()

    // 2. Add OpenTelemetry support for log-trace correlation
    composite_logger.AddHook(hook.OpenTelemetryHook{})

    // 3. Setup panic recovery for the main goroutine
    defer composite_logger.Recover(map[string]interface{}{"scope": "main"})

    // 4. Use context-aware logging
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    composite_logger.InfoContext(ctx, "System initialized", map[string]interface{}{
        "version": "v1.2.0",
    })

    // 5. Bound context logger (convenience wrapper)
    logger := composite_logger.WithContext(ctx)
    logger.Warn("Resource usage high", map[string]interface{}{"cpu": "85%"})
}
```

---

## Learning by Example

Explore the [examples/](./examples) directory for detailed implementation guides:

- **Basic Usage**: [Console](./examples/console) | [File](./examples/file) | [Telegram](./examples/telegram)
- **Architecture**: [Composite Pattern](./examples/composite) | [Custom Adapters](./examples/custom-adapter)
- **Advanced**: [OpenTelemetry Integration](./examples/otel-hook) | [Custom Hooks](./examples/custom-hook)

## License

MIT
