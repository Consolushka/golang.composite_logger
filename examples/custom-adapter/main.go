package main

import (
	"fmt"

	"github.com/Consolushka/golang.composite_logger/pkg"
)

// 1. Define your custom logger that implements ports.Logger interface
type MyCustomLogger struct {
	Prefix string
}

func (m MyCustomLogger) Info(msg string, ctx map[string]interface{}) {
	fmt.Printf("%s [INFO] %s | Context: %v
", m.Prefix, msg, ctx)
}

func (m MyCustomLogger) Warn(msg string, ctx map[string]interface{}) {
	fmt.Printf("%s [WARN] %s
", m.Prefix, msg)
}

func (m MyCustomLogger) Error(msg string, ctx map[string]interface{}) {
	fmt.Printf("%s [ERROR] %s
", m.Prefix, msg)
}

func (m MyCustomLogger) Fatal(msg string, ctx map[string]interface{}) {
	fmt.Printf("%s [FATAL] %s
", m.Prefix, msg)
}

// 2. Define your setting that implements ports.LoggerSetting interface
type MyCustomSetting struct {
	Enabled bool
	Prefix  string
}

func (s MyCustomSetting) InitLogger() composite_logger.Logger {
	return MyCustomLogger{Prefix: s.Prefix}
}

func (s MyCustomSetting) IsEnabled() bool {
	return s.Enabled
}

func main() {
	// 3. Plug in your custom adapter along with standard ones
	composite_logger.Init(
		MyCustomSetting{
			Enabled: true,
			Prefix:  "🚀[MyApp]",
		},
	)
	defer composite_logger.Stop()

	composite_logger.Info("Using custom adapter", map[string]interface{}{"status": "ok"})
}
