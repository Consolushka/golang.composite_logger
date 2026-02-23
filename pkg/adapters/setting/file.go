package setting

import (
	"io"
	"os"
	"path/filepath"

	"github.com/Consolushka/golang.composite_logger/internal/adapters/logger"
	compositelogger "github.com/Consolushka/golang.composite_logger/pkg"
	"github.com/Consolushka/golang.composite_logger/pkg/ports"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

type FileSetting struct {
	Enabled         bool
	IsJsonFormatter *bool
	Path            string
	LowerLevel      compositelogger.Level
	MaxSize         int  // Maximum size in megabytes before rotation (default: 5)
	MaxBackups      int  // Maximum number of old log files to retain (default: 3)
	MaxAge          int  // Maximum number of days to retain old log files (default: 28)
	Compress        bool // Whether to compress old log files (default: true)
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

	lumberjackLogger := f.setupRotation()

	mw := io.MultiWriter(os.Stdout, lumberjackLogger)
	logrusInstance.SetOutput(mw)

	return logger.NewFileLogger(logrusInstance)
}

func (f FileSetting) setupRotation() *lumberjack.Logger {
	// Set sensible defaults for rotation if not specified
	maxSize := f.MaxSize
	if maxSize == 0 {
		maxSize = 5
	}
	maxBackups := f.MaxBackups
	if maxBackups == 0 {
		maxBackups = 3
	}
	maxAge := f.MaxAge
	if maxAge == 0 {
		maxAge = 28
	}

	// Use lumberjack for log rotation
	lumberjackLogger := &lumberjack.Logger{
		Filename:   f.Path,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   f.Compress,
	}

	return lumberjackLogger
}

func (f FileSetting) IsEnabled() bool {
	return f.Enabled
}
