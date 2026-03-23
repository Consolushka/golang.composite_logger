package ports

import "context"

// CompositeLogger defines the interface for the central logging hub.
// It manages a collection of loggers and handles log dispatching according to the implementation (sync/async).
type CompositeLogger interface {
	// Log dispatches a log entry to all registered adapters.
	// level is the integer representation of the log level (1: Info, 2: Warning, 3: Error, 4: Fatal).
	Log(level int, message string, fields map[string]interface{})
	// LogContext dispatches a log entry with calling context to all registered adapters.
	LogContext(ctx context.Context, level int, message string, fields map[string]interface{})
	// AddHook registers a new hook to the logger instance.
	AddHook(hook Hook)
	// SetContextKeys registers context keys for automatic log enrichment.
	SetContextKeys(keys ...any)
	// Stop gracefully shuts down the logger, flushing any pending logs.
	Stop()
}
