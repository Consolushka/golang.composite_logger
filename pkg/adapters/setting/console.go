package setting

import (
	"github.com/Consolushka/golang.composite_logger/internal/adapters/logger"
	compositelogger "github.com/Consolushka/golang.composite_logger/pkg"
	"github.com/Consolushka/golang.composite_logger/pkg/ports"
	"github.com/sirupsen/logrus"
)

type ConsoleSetting struct {
	Enabled         bool
	IsJsonFormatter *bool
	LowerLevel      compositelogger.Level
}

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

func (s ConsoleSetting) IsEnabled() bool {
	return s.Enabled
}
