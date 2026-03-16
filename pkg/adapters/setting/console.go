package setting

import (
	"github.com/Consolushka/golang.composite_logger/internal/adapters/logger"
	compositelogger "github.com/Consolushka/golang.composite_logger/pkg"
	"github.com/Consolushka/golang.composite_logger/pkg/ports"
	"github.com/sirupsen/logrus"
)

// ConsoleSetting provides configuration for the standard output logging adapter.
// It uses Logrus internally to format and write logs to the console.
type ConsoleSetting struct {
	// Enabled toggles the console logger on or off.
	Enabled bool

	// IsJsonFormatter enables JSON output format if true.
	// If false, it uses the default Logrus text formatter.
	// Defaults to true if nil.
	IsJsonFormatter *bool

	// LowerLevel sets the minimum severity level that this adapter will process.
	LowerLevel compositelogger.Level
}

// InitLogger initializes a logrus-based console logger using the provided settings.
// It implements the ports.LoggerSetting interface.
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

// IsEnabled returns the current active status of the console adapter.
func (s ConsoleSetting) IsEnabled() bool {
	return s.Enabled
}
