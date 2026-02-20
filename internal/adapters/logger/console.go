package logger

import "github.com/sirupsen/logrus"

type ConsoleLogger struct {
	logrus *logrus.Logger
}

func NewConsoleLogger(logger *logrus.Logger) ConsoleLogger {
	return ConsoleLogger{
		logrus: logger,
	}
}

func (c ConsoleLogger) Info(message string, context map[string]interface{}) {
	c.logrus.WithFields(context).Info(message)
}

func (c ConsoleLogger) Warn(message string, context map[string]interface{}) {
	c.logrus.WithFields(context).Warn(message)
}

func (c ConsoleLogger) Error(message string, context map[string]interface{}) {
	c.logrus.WithFields(context).Error(message)
}

func (c ConsoleLogger) Fatal(message string, context map[string]interface{}) {
	c.logrus.WithFields(context).Log(logrus.FatalLevel, message)
}
