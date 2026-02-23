package setting

import (
	"github.com/Consolushka/golang.composite_logger/internal/adapters/logger"
	compositelogger "github.com/Consolushka/golang.composite_logger/pkg"
	"github.com/Consolushka/golang.composite_logger/pkg/ports"
	"github.com/sirupsen/logrus"
)

// ConsoleSetting provides configuration for the standard output logging adapter.
type ConsoleSetting struct {
	// Enabled toggles the console logger on or off.
	Enabled         bool
	// IsJsonFormatter enables JSON output format if true (default: true).
	IsJsonFormatter *bool
	// LowerLevel sets the minimum severity level to log.
	LowerLevel      compositelogger.Level
}

// InitLogger initializes a logrus-based console logger.
func (s ConsoleSetting) InitLogger() ports.Logger {
	logrusInstance := logrus.New()
	logrusInstance.SetLevel(s.LowerLevel.ToLogrus())

	isJsonFormatter := true
	if s.IsJsonFormatter != nil {
		isJsonFormatter = *s.IsJsonFormatter
	}

	if isJsonFormatter {
		logrusInstance.SetFormatter(&logrus.JSONFormatter{})
	}

	return logger.NewConsoleLogger(logrusInstance)
}

// IsEnabled returns the current active status of the adapter.
func (s ConsoleSetting) IsEnabled() bool {
	return s.Enabled
}
