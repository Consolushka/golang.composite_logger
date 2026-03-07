package logger

import (
	"context"

	"github.com/sirupsen/logrus"
)

type FileLogger struct {
	logrus *logrus.Logger
}

func NewFileLogger(logrusInstance *logrus.Logger) FileLogger {
	return FileLogger{logrusInstance}
}

func (f FileLogger) Info(message string, fields map[string]interface{}) {
	f.logrus.WithFields(fields).Info(message)
}

func (f FileLogger) InfoContext(ctx context.Context, message string, fields map[string]interface{}) {
	if ctx != nil && ctx.Err() != nil {
		return
	}
	f.logrus.WithContext(ctx).WithFields(fields).Info(message)
}

func (f FileLogger) Warn(message string, fields map[string]interface{}) {
	f.logrus.WithFields(fields).Warn(message)
}

func (f FileLogger) WarnContext(ctx context.Context, message string, fields map[string]interface{}) {
	if ctx != nil && ctx.Err() != nil {
		return
	}
	f.logrus.WithContext(ctx).WithFields(fields).Warn(message)
}

func (f FileLogger) Error(message string, fields map[string]interface{}) {
	f.logrus.WithFields(fields).Error(message)
}

func (f FileLogger) ErrorContext(ctx context.Context, message string, fields map[string]interface{}) {
	if ctx != nil && ctx.Err() != nil {
		return
	}
	f.logrus.WithContext(ctx).WithFields(fields).Error(message)
}

func (f FileLogger) Fatal(message string, fields map[string]interface{}) {
	f.logrus.WithFields(fields).Log(logrus.FatalLevel, message)
}

func (f FileLogger) FatalContext(ctx context.Context, message string, fields map[string]interface{}) {
	if ctx != nil && ctx.Err() != nil {
		return
	}
	f.logrus.WithContext(ctx).WithFields(fields).Log(logrus.FatalLevel, message)
}
