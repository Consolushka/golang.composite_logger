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

// FileSetting provides configuration for the file-based logging adapter with automatic rotation.
type FileSetting struct {
	// Enabled toggles the file logger on or off.
	Enabled         bool
	// IsJsonFormatter enables JSON output format if true (default: true).
	IsJsonFormatter *bool
	// Path is the filesystem path to the log file.
	Path            string
	// LowerLevel sets the minimum severity level to log.
	LowerLevel      compositelogger.Level
	// MaxSize is the maximum size in megabytes before rotation (default: 5).
	MaxSize         int
	// MaxBackups is the maximum number of old log files to retain (default: 3).
	MaxBackups      int
	// MaxAge is the maximum number of days to retain old log files (default: 28).
	MaxAge          int
	// Compress determines if old log files should be compressed (default: true).
	Compress        bool
}

// InitLogger initializes a logrus-based file logger with lumberjack for rotation.
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

// setupRotation configures the lumberjack logger with defaults and user-provided values.
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

// IsEnabled returns the current active status of the adapter.
func (f FileSetting) IsEnabled() bool {
	return f.Enabled
}
