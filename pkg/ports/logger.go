package ports

// Logger defines the interface for all logging adapters.
type Logger interface {
	// Info logs a message with informational severity.
	Info(message string, context map[string]interface{})
	// Warn logs a message with warning severity.
	Warn(message string, context map[string]interface{})
	// Error logs a message with error severity.
	Error(message string, context map[string]interface{})
	// Fatal logs a message with fatal severity and may lead to process termination.
	Fatal(message string, context map[string]interface{})
}
