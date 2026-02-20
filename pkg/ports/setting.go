package ports

type LoggerSetting interface {
	InitLogger() Logger
}
