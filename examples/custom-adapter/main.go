package main

import (
	"context"
	"fmt"

	"github.com/Consolushka/golang.composite_logger/pkg"
)

// 1. Define your custom logger that implements ports.Logger interface
type MyCustomLogger struct {
	Prefix string
}

func (m MyCustomLogger) Info(msg string, fields map[string]interface{}) {
	fmt.Printf("%s [INFO] %s | Context: %v\n", m.Prefix, msg, fields)
}

func (m MyCustomLogger) InfoContext(ctx context.Context, msg string, fields map[string]interface{}) {
	fmt.Printf("%s [INFO] %s | Context: %v | Fields: %v\n", m.Prefix, msg, ctx, fields)
}

func (m MyCustomLogger) Warn(msg string, fields map[string]interface{}) {
	fmt.Printf("%s [WARN] %s\n", m.Prefix, msg)
}

func (m MyCustomLogger) WarnContext(ctx context.Context, msg string, fields map[string]interface{}) {
	fmt.Printf("%s [WARN] %s | Fields: %v\n", m.Prefix, msg, fields)
}

func (m MyCustomLogger) Error(msg string, fields map[string]interface{}) {
	fmt.Printf("%s [ERROR] %s\n", m.Prefix, msg)
}

func (m MyCustomLogger) ErrorContext(ctx context.Context, msg string, fields map[string]interface{}) {
	fmt.Printf("%s [ERROR] %s | Fields: %v\n", m.Prefix, msg, fields)
}

func (m MyCustomLogger) Fatal(msg string, fields map[string]interface{}) {
	fmt.Printf("%s [FATAL] %s\n", m.Prefix, msg)
}

func (m MyCustomLogger) FatalContext(ctx context.Context, msg string, fields map[string]interface{}) {
	fmt.Printf("%s [FATAL] %s | Fields: %v\n", m.Prefix, msg, fields)
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
