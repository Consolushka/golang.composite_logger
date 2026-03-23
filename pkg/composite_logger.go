package composite_logger

import (
	"context"
	"sync"

	"github.com/Consolushka/golang.composite_logger/internal"
	hub "github.com/Consolushka/golang.composite_logger/internal/adapters/composite_logger"
	"github.com/Consolushka/golang.composite_logger/pkg/ports"
)

var (
	instance ports.CompositeLogger
	mu       sync.Mutex
)

// Logger is a re-export of ports.Logger for convenience.
type Logger = ports.Logger

// LoggerSetting is a re-export of ports.LoggerSetting for convenience.
type LoggerSetting = ports.LoggerSetting

// Hook is a re-export of ports.Hook for convenience.
type Hook = ports.Hook

// CompositeLogger is a re-export of ports.CompositeLogger for convenience.
type CompositeLogger = ports.CompositeLogger

// LoggingContext represents a logger instance bound to a specific context.
// It provides a convenient way to log multiple messages without passing the context every time.
type LoggingContext struct {
	ctx context.Context
}

// WithContext creates a new LoggingContext bound to the provided context.
// All logging methods on the returned instance will use this context.
func WithContext(ctx context.Context) *LoggingContext {
	return &LoggingContext{ctx: ctx}
}

// Info logs a message with the INFO level using the bound context.
func (l *LoggingContext) Info(msg string, fields map[string]interface{}) {
	InfoContext(l.ctx, msg, fields)
}

// Warn logs a message with the WARNING level using the bound context.
func (l *LoggingContext) Warn(msg string, fields map[string]interface{}) {
	WarnContext(l.ctx, msg, fields)
}

// Error logs a message with the ERROR level using the bound context.
func (l *LoggingContext) Error(msg string, fields map[string]interface{}) {
	ErrorContext(l.ctx, msg, fields)
}

// Fatal logs a message with the FATAL level using the bound context.
func (l *LoggingContext) Fatal(msg string, fields map[string]interface{}) {
	FatalContext(l.ctx, msg, fields)
}

// Init initializes the global logger instance with the provided settings in SYNCHRONOUS mode.
// This is the default mode where log calls block until all adapters have processed the message.
// If an instance already exists, it will be gracefully shut down before the new instance is started.
func Init(settings ...LoggerSetting) {
	mu.Lock()
	defer mu.Unlock()

	if instance != nil {
		instance.Stop()
	}

	loggers := initLoggers(settings)
	instance = hub.NewSyncCompositeLogger(loggers)
}

// InitAsync initializes the global logger instance in ASYNCHRONOUS mode.
// Log calls are non-blocking, as entries are processed by a background worker goroutine.
// If an instance already exists, it will be gracefully shut down before the new instance is started.
//
// CRITICAL: Always call Stop() when using async mode to ensure all queued logs are flushed before exit.
func InitAsync(settings ...LoggerSetting) {
	mu.Lock()
	defer mu.Unlock()

	if instance != nil {
		instance.Stop()
	}

	loggers := initLoggers(settings)
	instance = hub.NewAsyncCompositeLogger(loggers, 1000)
}

func initLoggers(settings []LoggerSetting) []ports.Logger {
	loggers := make([]ports.Logger, 0, len(settings))
	for _, s := range settings {
		if s.IsEnabled() {
			loggers = append(loggers, s.InitLogger())
		}
	}
	return loggers
}

// AddHook adds a hook to the global logger instance.
// Hooks are executed for each log entry before it is broadcast to adapters.
func AddHook(hook ports.Hook) {
	mu.Lock()
	defer mu.Unlock()
	if instance != nil {
		instance.AddHook(hook)
	}
}

// SetContextKeys registers a list of context keys that the logger should automatically
// extract from the provided context and add to the log fields.
func SetContextKeys(keys ...any) {
	mu.Lock()
	defer mu.Unlock()
	if instance != nil {
		instance.SetContextKeys(keys...)
	}
}

// Stop gracefully shuts down the global logger, ensuring all queued logs are processed
// before the application exits. It is highly recommended to defer this call in your main function.
func Stop() {
	mu.Lock()
	defer mu.Unlock()
	if instance != nil {
		instance.Stop()
		instance = nil
	}
}

// Info logs a message with the INFO level.
func Info(msg string, fields map[string]interface{}) {
	mu.Lock()
	defer mu.Unlock()
	if instance == nil {
		return
	}
	instance.Log(int(InfoLevel), "[INFO] "+msg, fields)
}

// InfoContext logs a message with the INFO level and provides calling context.
func InfoContext(ctx context.Context, msg string, fields map[string]interface{}) {
	mu.Lock()
	defer mu.Unlock()
	if instance == nil {
		return
	}
	instance.LogContext(ctx, int(InfoLevel), "[INFO] "+msg, fields)
}

// Warn logs a message with the WARNING level.
func Warn(msg string, fields map[string]interface{}) {
	mu.Lock()
	defer mu.Unlock()
	if instance == nil {
		return
	}
	instance.Log(int(WarningLevel), "[WARNING] "+msg, fields)
}

// WarnContext logs a message with the WARNING level and provides calling context.
func WarnContext(ctx context.Context, msg string, fields map[string]interface{}) {
	mu.Lock()
	defer mu.Unlock()
	if instance == nil {
		return
	}
	instance.LogContext(ctx, int(WarningLevel), "[WARNING] "+msg, fields)
}

// Error captures a stack trace and logs a message with the ERROR level.
func Error(msg string, fields map[string]interface{}) {
	fields = internal.BuildErrorContextWithStackTrace(fields)

	mu.Lock()
	defer mu.Unlock()
	if instance == nil {
		return
	}
	instance.Log(int(ErrorLevel), "[ERROR] "+msg, fields)
}

// ErrorContext captures a stack trace and logs a message with the ERROR level and context.
func ErrorContext(ctx context.Context, msg string, fields map[string]interface{}) {
	fields = internal.BuildErrorContextWithStackTrace(fields)

	mu.Lock()
	defer mu.Unlock()
	if instance == nil {
		return
	}
	instance.LogContext(ctx, int(ErrorLevel), "[ERROR] "+msg, fields)
}

// Fatal captures a stack trace and logs a message with the FATAL level.
func Fatal(msg string, fields map[string]interface{}) {
	fields = internal.BuildErrorContextWithStackTrace(fields)

	mu.Lock()
	defer mu.Unlock()
	if instance == nil {
		return
	}
	instance.Log(int(FatalLevel), "[FATAL] "+msg, fields)
}

// FatalContext captures a stack trace and logs a message with the FATAL level and context.
func FatalContext(ctx context.Context, msg string, fields map[string]interface{}) {
	fields = internal.BuildErrorContextWithStackTrace(fields)

	mu.Lock()
	defer mu.Unlock()
	if instance == nil {
		return
	}
	instance.LogContext(ctx, int(FatalLevel), "[FATAL] "+msg, fields)
}

// Recover is a helper function to be used in defer statements to catch and log panics as FATAL errors.
func Recover(fields map[string]interface{}) {
	if r := recover(); r != nil {
		Fatal("Panic recovered", map[string]interface{}{
			"panic":  r,
			"fields": fields,
		})
	}
}

// RecoverContext is a helper function to be used in defer statements to catch and log panics as FATAL errors with context.
func RecoverContext(ctx context.Context, fields map[string]interface{}) {
	if r := recover(); r != nil {
		FatalContext(ctx, "Panic recovered", map[string]interface{}{
			"panic":  r,
			"fields": fields,
		})
	}
}
