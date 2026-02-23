package composite_logger

import (
	"errors"
	"strings"

	"github.com/sirupsen/logrus"
)

// Level represents the severity of a log message.
type Level int

const (
	// InfoLevel is used for informational messages that highlight the progress of the application.
	InfoLevel Level = 1
	// WarningLevel is used for potentially harmful situations or important notices.
	WarningLevel Level = 2
	// ErrorLevel is used for error events that might still allow the application to continue running.
	ErrorLevel Level = 3
	// FatalLevel is used for very severe error events that will presumably lead the application to abort.
	FatalLevel Level = 4
)

// DefaultLevelWrappers provides a default set of emoji wrappers for each log level,
// commonly used by the Telegram adapter.
var DefaultLevelWrappers = map[Level]string{
	InfoLevel:    "ℹ️ℹ️",
	WarningLevel: "⚠️⚠️",
	ErrorLevel:   "‼️‼️",
	FatalLevel:   "🚨🚨",
}

// String returns the lower-case string representation of the log level.
func (l Level) String() string {
	switch l {
	case InfoLevel:
		return "info"
	case WarningLevel:
		return "warning"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	default:
		return "info"
	}
}

// ToLogrus converts the internal Level to a logrus.Level.
func (l Level) ToLogrus() logrus.Level {
	switch l {
	case InfoLevel:
		return logrus.InfoLevel
	case WarningLevel:
		return logrus.WarnLevel
	case ErrorLevel:
		return logrus.ErrorLevel
	case FatalLevel:
		return logrus.FatalLevel
	default:
		return logrus.InfoLevel
	}
}

// ParseLevel parses a string into a Level.
// It returns an error if the string does not match any known log level.
func ParseLevel(lvl string) (Level, error) {
	switch strings.ToLower(lvl) {
	case "info":
		return InfoLevel, nil
	case "warn", "warning":
		return WarningLevel, nil
	case "error":
		return ErrorLevel, nil
	case "fatal":
		return FatalLevel, nil
	}

	return InfoLevel, errors.New("invalid log level: " + lvl)
}
