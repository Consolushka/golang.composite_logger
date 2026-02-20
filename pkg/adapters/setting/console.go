package setting

import (
	"github.com/Consolushka/golang.composite_logger/internal/adapters/logger"
	compositelogger "github.com/Consolushka/golang.composite_logger/pkg"
	"github.com/Consolushka/golang.composite_logger/pkg/ports"
	"github.com/sirupsen/logrus"
)

type ConsoleSetting struct {
	LowerLevel compositelogger.Level
}

func (s ConsoleSetting) InitLogger() ports.Logger {
	logrusInstance := logrus.New()
	logrusInstance.SetLevel(s.LowerLevel.ToLogrus())

	return logger.NewConsoleLogger(logrusInstance)
}
