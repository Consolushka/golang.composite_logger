package composite_logger

import (
	"composite_logger/internal"
	"composite_logger/pkg/ports"
)

var instance CompositeLogger

type CompositeLogger struct {
	loggers []ports.Logger
}

func Init(settings ...ports.LoggerSetting) {
	loggers := make([]ports.Logger, 0, len(settings))
	for _, s := range settings {
		loggers = append(loggers, s.InitLogger())
	}

	instance = CompositeLogger{
		loggers: loggers,
	}
}

func Info(msg string, ctx map[string]interface{}) {
	for _, logger := range instance.loggers {
		logger.Info("[INFO] "+msg, ctx)
	}
}

func Warn(msg string, ctx map[string]interface{}) {
	for _, logger := range instance.loggers {
		logger.Warn("[WARNING] "+msg, ctx)
	}
}

func Error(msg string, ctx map[string]interface{}) {
	ctx = internal.BuildErrorContextWithStackTrace(ctx)

	for _, logger := range instance.loggers {
		logger.Error("[ERROR] "+msg, ctx)
	}
}

func Fatal(msg string, ctx map[string]interface{}) {
	ctx = internal.BuildErrorContextWithStackTrace(ctx)

	for _, logger := range instance.loggers {
		logger.Fatal("[FATAL] "+msg, ctx)
	}
}

func Recover(ctx map[string]interface{}) {
	if r := recover(); r != nil {
		Fatal("Panic recovered", map[string]interface{}{
			"panic": r,
			"ctx":   ctx,
		})
	}
}
