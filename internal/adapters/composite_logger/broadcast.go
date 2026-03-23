package composite_logger

import (
	"context"
	"fmt"

	"github.com/Consolushka/golang.composite_logger/pkg/ports"
)

func broadcast(ctx context.Context, level int, message string, fields map[string]interface{}, loggers []ports.Logger, hooks []ports.Hook, contextKeys []any) {
	// 1. Enrich fields from context keys
	enrichedFields := enrichFromContext(ctx, fields, contextKeys)

	// 2. Fire registered hooks
	fireHooks(ctx, level, message, enrichedFields, hooks)

	// 3. Dispatch to each logger adapter
	fanOutToLoggers(ctx, level, message, enrichedFields, loggers)
}

func enrichFromContext(ctx context.Context, fields map[string]interface{}, keys []any) map[string]interface{} {
	if len(keys) == 0 {
		return fields
	}

	if fields == nil {
		fields = make(map[string]interface{})
	}

	for _, key := range keys {
		if val := ctx.Value(key); val != nil {
			keyStr, ok := key.(string)
			if !ok {
				keyStr = fmt.Sprintf("%v", key)
			}
			if _, exists := fields[keyStr]; !exists {
				fields[keyStr] = val
			}
		}
	}

	return fields
}

func fireHooks(ctx context.Context, level int, message string, fields map[string]interface{}, hooks []ports.Hook) {
	if len(hooks) == 0 {
		return
	}

	levelStr := levelToString(level)
	for _, hook := range hooks {
		_ = hook.Fire(ctx, levelStr, message, fields)
	}
}

func fanOutToLoggers(ctx context.Context, level int, message string, fields map[string]interface{}, loggers []ports.Logger) {
	for _, l := range loggers {
		switch level {
		case 1: // Info
			l.InfoContext(ctx, message, fields)
		case 2: // Warning
			l.WarnContext(ctx, message, fields)
		case 3: // Error
			l.ErrorContext(ctx, message, fields)
		case 4: // Fatal
			l.FatalContext(ctx, message, fields)
		}
	}
}

func levelToString(level int) string {
	switch level {
	case 1:
		return "info"
	case 2:
		return "warning"
	case 3: // error
		return "error"
	case 4: // fatal
		return "fatal"
	default:
		return "info"
	}
}
