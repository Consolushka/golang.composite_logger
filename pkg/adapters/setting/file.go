package setting

import (
	"io"
	"os"
	"path/filepath"

	"github.com/Consolushka/golang.composite_logger/internal/adapters/logger"
	compositelogger "github.com/Consolushka/golang.composite_logger/pkg"
	"github.com/Consolushka/golang.composite_logger/pkg/ports"
	"github.com/sirupsen/logrus"
)

type FileSetting struct {
	Enabled         bool
	IsJsonFormatter *bool
	Path            string
	LowerLevel      compositelogger.Level
}

func (f FileSetting) InitLogger() ports.Logger {
	if f.Path == "" {
		panic("File path is not set")
	}

	logrusInstance := logrus.New()
	logrusInstance.SetLevel(f.LowerLevel.ToLogrus())

	isJsonFormatter := true
	if f.IsJsonFormatter != nil {
		isJsonFormatter = *f.IsJsonFormatter
	}

	if isJsonFormatter {
		logrusInstance.SetFormatter(&logrus.JSONFormatter{})
	}

	logDir := filepath.Dir(f.Path)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		logrusInstance.Fatalf("Failed to create logrusInstance directory: %v", err)
	}

	// Keep opened to write logs into it
	logFile, err := os.OpenFile(f.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrusInstance.Fatalf("Failed to open logrusInstance file: %v", err)
	}

	mw := io.MultiWriter(os.Stdout, logFile)
	logrusInstance.SetOutput(mw)

	return logger.NewFileLogger(logrusInstance)
}

func (f FileSetting) IsEnabled() bool {
	return f.Enabled
}
