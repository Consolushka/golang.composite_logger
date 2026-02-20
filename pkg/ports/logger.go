package ports

type Logger interface {
	Info(message string, context map[string]interface{})
	Warn(message string, context map[string]interface{})
	Error(message string, context map[string]interface{})
	Fatal(message string, context map[string]interface{})
}
