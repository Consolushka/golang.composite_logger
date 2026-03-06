package composite_logger

import (
	"context"
	"sync"

	"github.com/Consolushka/golang.composite_logger/internal"
	"github.com/Consolushka/golang.composite_logger/pkg/ports"
)

var (
	instance *CompositeLogger
	mu       sync.Mutex
)

// Logger is a re-export of ports.Logger for convenience.
type Logger = ports.Logger

// LoggerSetting is a re-export of ports.LoggerSetting for convenience.
type LoggerSetting = ports.LoggerSetting

type logEntry struct {
	level   Level
	message string
	context map[string]interface{}
	ctx     context.Context
}

// CompositeLogger manages a collection of loggers and handles asynchronous log dispatching.
type CompositeLogger struct {
	loggers []ports.Logger
	ch      chan logEntry
	wg      sync.WaitGroup
}

// Init initializes the global logger instance with the provided settings.
// If an instance already exists, it will be gracefully shut down before re-initialization.
//
// Usage:
//
//	composite_logger.Init(setting.ConsoleSetting{LowerLevel: composite_logger.InfoLevel}, setting.TelegramSetting{LowerLevel: composite_logger.InfoLevel, Enabled: true})
func Init(settings ...LoggerSetting) {
	mu.Lock()
	defer mu.Unlock()

	// In case of re-initializing composite logger
	if instance != nil {
		instance.flushAndClose()
	}

	loggers := make([]ports.Logger, 0, len(settings))
	for _, s := range settings {
		if s.IsEnabled() {
			loggers = append(loggers, s.InitLogger())
		}
	}

	instance = &CompositeLogger{
		loggers: loggers,
		ch:      make(chan logEntry, 1000),
	}

	instance.wg.Add(1)
	go instance.listenAndBroadcast()
}

// listenAndBroadcast is a background worker that processes the log queue
// and sends entries to all registered adapters.
func (cl *CompositeLogger) listenAndBroadcast() {
	defer cl.wg.Done()
	for entry := range cl.ch {
		for _, logger := range cl.loggers {
			if entry.ctx != nil {
				switch entry.level {
				case InfoLevel:
					logger.InfoContext(entry.ctx, entry.message, entry.context)
				case WarningLevel:
					logger.WarnContext(entry.ctx, entry.message, entry.context)
				case ErrorLevel:
					logger.ErrorContext(entry.ctx, entry.message, entry.context)
				case FatalLevel:
					logger.FatalContext(entry.ctx, entry.message, entry.context)
				}
			} else {
				switch entry.level {
				case InfoLevel:
					logger.Info(entry.message, entry.context)
				case WarningLevel:
					logger.Warn(entry.message, entry.context)
				case ErrorLevel:
					logger.Error(entry.message, entry.context)
				case FatalLevel:
					logger.Fatal(entry.message, entry.context)
				}
			}
		}
	}
}

// flushAndClose closes the log queue and waits for the worker to finish processing remaining entries.
func (cl *CompositeLogger) flushAndClose() {
	close(cl.ch)
	cl.wg.Wait()
}

// Stop gracefully shuts down the global logger, ensuring all queued logs are processed
// before the application exits. It is highly recommended to defer this call in your main function.
//
// Usage:
//
//	defer composite_logger.Stop()
func Stop() {
	mu.Lock()
	defer mu.Unlock()
	if instance != nil {
		instance.flushAndClose()
		instance = nil
	}
}

// Info asynchronously logs a message with the INFO level.
//
// Usage:
//
//	composite_logger.Info("app started", map[string]interface{}{"env": "prod"})
func Info(msg string, ctx map[string]interface{}) {
	mu.Lock()
	defer mu.Unlock()
	if instance == nil || instance.ch == nil {
		return
	}
	instance.ch <- logEntry{
		level:   InfoLevel,
		message: "[INFO] " + msg,
		context: ctx,
	}
}

// InfoContext asynchronously logs a message with the INFO level and provides calling context.
// Use this to correlate logs with specific requests or background tasks.
//
// Usage:
//
//	composite_logger.InfoContext(ctx, "processing user request", map[string]interface{}{"userID": 123})
func InfoContext(ctx context.Context, msg string, fields map[string]interface{}) {
	mu.Lock()
	defer mu.Unlock()
	if instance == nil || instance.ch == nil {
		return
	}
	instance.ch <- logEntry{
		level:   InfoLevel,
		message: "[INFO] " + msg,
		context: fields,
		ctx:     ctx,
	}
}

// Warn asynchronously logs a message with the WARNING level.
//
// Usage:
//
//	composite_logger.Warn("high latency", map[string]interface{}{"ms": 500})
func Warn(msg string, ctx map[string]interface{}) {
	mu.Lock()
	defer mu.Unlock()
	if instance == nil || instance.ch == nil {
		return
	}
	instance.ch <- logEntry{
		level:   WarningLevel,
		message: "[WARNING] " + msg,
		context: ctx,
	}
}

// WarnContext asynchronously logs a message with the WARNING level and provides calling context.
func WarnContext(ctx context.Context, msg string, fields map[string]interface{}) {
	mu.Lock()
	defer mu.Unlock()
	if instance == nil || instance.ch == nil {
		return
	}
	instance.ch <- logEntry{
		level:   WarningLevel,
		message: "[WARNING] " + msg,
		context: fields,
		ctx:     ctx,
	}
}

// Error captures a stack trace and asynchronously logs a message with the ERROR level.
//
// Usage:
//
//	composite_logger.Error("db connection failed", map[string]interface{}{"error": err})
func Error(msg string, ctx map[string]interface{}) {
	ctx = internal.BuildErrorContextWithStackTrace(ctx)

	mu.Lock()
	defer mu.Unlock()
	if instance == nil || instance.ch == nil {
		return
	}
	instance.ch <- logEntry{
		level:   ErrorLevel,
		message: "[ERROR] " + msg,
		context: ctx,
	}
}

// ErrorContext captures a stack trace and asynchronously logs a message with the ERROR level and context.
func ErrorContext(ctx context.Context, msg string, fields map[string]interface{}) {
	fields = internal.BuildErrorContextWithStackTrace(fields)

	mu.Lock()
	defer mu.Unlock()
	if instance == nil || instance.ch == nil {
		return
	}
	instance.ch <- logEntry{
		level:   ErrorLevel,
		message: "[ERROR] " + msg,
		context: fields,
		ctx:     ctx,
	}
}

// Fatal captures a stack trace and asynchronously logs a message with the FATAL level.
// This level indicates a critical system failure. Note: This library does not call os.Exit;
// your application logic should decide when to terminate.
//
// Usage:
//
//	composite_logger.Fatal("system crashed", map[string]interface{}{"reason": "out of memory"})
func Fatal(msg string, ctx map[string]interface{}) {
	ctx = internal.BuildErrorContextWithStackTrace(ctx)

	mu.Lock()
	defer mu.Unlock()
	if instance == nil || instance.ch == nil {
		return
	}
	instance.ch <- logEntry{
		level:   FatalLevel,
		message: "[FATAL] " + msg,
		context: ctx,
	}
}

// FatalContext captures a stack trace and asynchronously logs a message with the FATAL level and context.
func FatalContext(ctx context.Context, msg string, fields map[string]interface{}) {
	fields = internal.BuildErrorContextWithStackTrace(fields)

	mu.Lock()
	defer mu.Unlock()
	if instance == nil || instance.ch == nil {
		return
	}
	instance.ch <- logEntry{
		level:   FatalLevel,
		message: "[FATAL] " + msg,
		context: fields,
		ctx:     ctx,
	}
}

// Recover is a helper function to be used in defer statements to catch and log panics as FATAL errors.
// It automatically captures the panic reason and the current stack trace.
//
// Usage:
//
//	defer composite_logger.Recover(map[string]interface{}{"handler": "user_create"})
func Recover(ctx map[string]interface{}) {
	if r := recover(); r != nil {
		Fatal("Panic recovered", map[string]interface{}{
			"panic": r,
			"ctx":   ctx,
		})
	}
}

// RecoverContext is a helper function to be used in defer statements to catch and log panics as FATAL errors with context.
// This is useful for capturing trace and request IDs during a crash.
//
// Usage:
//
//	defer composite_logger.RecoverContext(ctx, map[string]interface{}{"handler": "user_create"})
func RecoverContext(ctx context.Context, fields map[string]interface{}) {
	if r := recover(); r != nil {
		FatalContext(ctx, "Panic recovered", map[string]interface{}{
			"panic":  r,
			"fields": fields,
		})
	}
}
