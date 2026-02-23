package ports

// LoggerSetting defines the interface for configuring and initializing a Logger adapter.
type LoggerSetting interface {
	// InitLogger creates and configures a concrete Logger implementation.
	InitLogger() Logger
	// IsEnabled returns true if the adapter should be active.
	IsEnabled() bool
}
