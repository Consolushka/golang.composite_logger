package logger

import (
	"context"

	"github.com/sirupsen/logrus"
)

type ConsoleLogger struct {
	logrus *logrus.Logger
}

func NewConsoleLogger(logger *logrus.Logger) ConsoleLogger {
	return ConsoleLogger{
		logrus: logger,
	}
}

func (c ConsoleLogger) Info(message string, fields map[string]interface{}) {
	c.logrus.WithFields(fields).Info(message)
}

func (c ConsoleLogger) InfoContext(ctx context.Context, message string, fields map[string]interface{}) {
	if ctx != nil && ctx.Err() != nil {
		return
	}
	c.logrus.WithContext(ctx).WithFields(fields).Info(message)
}

func (c ConsoleLogger) Warn(message string, fields map[string]interface{}) {
	c.logrus.WithFields(fields).Warn(message)
}

func (c ConsoleLogger) WarnContext(ctx context.Context, message string, fields map[string]interface{}) {
	if ctx != nil && ctx.Err() != nil {
		return
	}
	c.logrus.WithContext(ctx).WithFields(fields).Warn(message)
}

func (c ConsoleLogger) Error(message string, fields map[string]interface{}) {
	c.logrus.WithFields(fields).Error(message)
}

func (c ConsoleLogger) ErrorContext(ctx context.Context, message string, fields map[string]interface{}) {
	if ctx != nil && ctx.Err() != nil {
		return
	}
	c.logrus.WithContext(ctx).WithFields(fields).Error(message)
}

func (c ConsoleLogger) Fatal(message string, fields map[string]interface{}) {
	c.logrus.WithFields(fields).Log(logrus.FatalLevel, message)
}

func (c ConsoleLogger) FatalContext(ctx context.Context, message string, fields map[string]interface{}) {
	if ctx != nil && ctx.Err() != nil {
		return
	}
	c.logrus.WithContext(ctx).WithFields(fields).Log(logrus.FatalLevel, message)
}
