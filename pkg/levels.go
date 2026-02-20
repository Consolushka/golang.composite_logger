package composite_logger

import "github.com/sirupsen/logrus"

type Level int

const (
	InfoLevel    Level = 1
	WarningLevel Level = 2
	ErrorLevel   Level = 3
	FatalLevel   Level = 4
)

var DefaultLevelWrappers = map[Level]string{
	InfoLevel:    "ℹ️ℹ️",
	WarningLevel: "⚠️⚠️",
	ErrorLevel:   "‼️‼️",
	FatalLevel:   "🚨🚨",
}

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
