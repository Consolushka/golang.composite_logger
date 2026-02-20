package setting

import (
	"composite_logger/internal/adapters/logger"
	composite_logger "composite_logger/pkg"
	"composite_logger/pkg/ports"

	"github.com/sirupsen/logrus"
)

type ConsoleSetting struct {
	LowerLevel composite_logger.Level
}

func (s ConsoleSetting) InitLogger() ports.Logger {
	logrusInstance := logrus.New()
	logrusInstance.SetLevel(s.LowerLevel.ToLogrus())

	return logger.NewConsoleLogger(logrusInstance)
}
