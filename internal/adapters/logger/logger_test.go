package logger

import (
	"bytes"
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func newBufferedLogrus() (*logrus.Logger, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	l := logrus.New()
	l.SetOutput(buf)
	l.SetLevel(logrus.DebugLevel)
	return l, buf
}

// Regression: a cancelled context must not suppress log output. The context
// here is a source of values, not a permission to log — by the time async
// workers process an entry the request context is always cancelled.
func TestConsoleLogger_WritesWithCancelledContext(t *testing.T) {
	instance, buf := newBufferedLogrus()
	consoleLogger := NewConsoleLogger(instance)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	consoleLogger.InfoContext(ctx, "info after cancel", nil)
	consoleLogger.WarnContext(ctx, "warn after cancel", nil)
	consoleLogger.ErrorContext(ctx, "error after cancel", nil)
	consoleLogger.FatalContext(ctx, "fatal after cancel", nil)

	output := buf.String()
	assert.Contains(t, output, "info after cancel")
	assert.Contains(t, output, "warn after cancel")
	assert.Contains(t, output, "error after cancel")
	assert.Contains(t, output, "fatal after cancel")
}

func TestFileLogger_WritesWithCancelledContext(t *testing.T) {
	instance, buf := newBufferedLogrus()
	fileLogger := NewFileLogger(instance)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	fileLogger.InfoContext(ctx, "info after cancel", nil)
	fileLogger.WarnContext(ctx, "warn after cancel", nil)
	fileLogger.ErrorContext(ctx, "error after cancel", nil)
	fileLogger.FatalContext(ctx, "fatal after cancel", nil)

	output := buf.String()
	assert.Contains(t, output, "info after cancel")
	assert.Contains(t, output, "warn after cancel")
	assert.Contains(t, output, "error after cancel")
	assert.Contains(t, output, "fatal after cancel")
}
