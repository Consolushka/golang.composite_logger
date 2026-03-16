package ports

import "context"

// Hook represents a component that can be notified of log events.
// It is called for each log entry before it is broadcast to the adapters.
type Hook interface {
	// Fire is called for each log entry.
	// It can modify the fields (passed by reference) or perform other side effects.
	// ctx is the context passed to the logging function.
	// level is the string representation of the log level (e.g., "info", "warning", "error", "fatal").
	// message is the log message.
	// fields is a map of key-value pairs associated with the log entry.
	Fire(ctx context.Context, level string, message string, fields map[string]interface{}) error
}
