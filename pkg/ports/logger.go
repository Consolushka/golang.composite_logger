package ports

import "context"

// Logger defines the interface for all logging adapters.
type Logger interface {
	// Info logs a message with informational severity.
	Info(message string, context map[string]interface{})

	// InfoContext logs a message with informational severity and provides the calling context.
	// This is useful for tracing, where the context may contain trace and span IDs.
	InfoContext(ctx context.Context, message string, fields map[string]interface{})

	// Warn logs a message with warning severity.
	Warn(message string, context map[string]interface{})

	// WarnContext logs a message with warning severity and provides the calling context.
	WarnContext(ctx context.Context, message string, fields map[string]interface{})

	// Error logs a message with error severity.
	Error(message string, context map[string]interface{})

	// ErrorContext logs a message with error severity and provides the calling context.
	ErrorContext(ctx context.Context, message string, fields map[string]interface{})

	// Fatal logs a message with fatal severity and may lead to process termination.
	Fatal(message string, context map[string]interface{})

	// FatalContext logs a message with fatal severity and provides the calling context.
	FatalContext(ctx context.Context, message string, fields map[string]interface{})
}
