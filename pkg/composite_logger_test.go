package composite_logger

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Consolushka/golang.composite_logger/internal"
	"github.com/Consolushka/golang.composite_logger/pkg/ports"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testSetting struct {
	l ports.Logger
}

func (t testSetting) InitLogger() ports.Logger {
	return t.l
}

func (t testSetting) IsEnabled() bool {
	return true
}

type stackAwareError struct {
	message string
	stack   string
}

func (s stackAwareError) Error() string {
	return s.message
}

func (s stackAwareError) Format(state fmt.State, verb rune) {
	if verb == 'v' {
		_, _ = state.Write([]byte(s.stack))
		return
	}

	_, _ = state.Write([]byte(s.message))
}

type logCall struct {
	message string
	context map[string]interface{}
}

type fakeLogger struct {
	infoCalls  []logCall
	warnCalls  []logCall
	errorCalls []logCall
	fatalCalls []logCall
}

func (f *fakeLogger) Info(message string, context map[string]interface{}) {
	f.infoCalls = append(f.infoCalls, logCall{message: message, context: context})
}

func (f *fakeLogger) Warn(message string, context map[string]interface{}) {
	f.warnCalls = append(f.warnCalls, logCall{message: message, context: context})
}

func (f *fakeLogger) Error(message string, context map[string]interface{}) {
	f.errorCalls = append(f.errorCalls, logCall{message: message, context: context})
}

func (f *fakeLogger) Fatal(message string, context map[string]interface{}) {
	f.fatalCalls = append(f.fatalCalls, logCall{message: message, context: context})
}

func TestInfo_FanOutAndPrefix(t *testing.T) {
	l1 := &fakeLogger{}
	l2 := &fakeLogger{}
	Init(testSetting{l1}, testSetting{l2})

	ctx := map[string]interface{}{"requestId": "abc-123"}
	Info("process started", ctx)
	Stop()

	require.Len(t, l1.infoCalls, 1)
	require.Len(t, l2.infoCalls, 1)
	assert.Equal(t, "[INFO] process started", l1.infoCalls[0].message)
	assert.Equal(t, "[INFO] process started", l2.infoCalls[0].message)
	assert.Equal(t, "abc-123", l1.infoCalls[0].context["requestId"])
}

func TestWarn_FanOutAndPrefix(t *testing.T) {
	l1 := &fakeLogger{}
	l2 := &fakeLogger{}
	Init(testSetting{l1}, testSetting{l2})

	ctx := map[string]interface{}{"task": "poll"}
	Warn("slow response", ctx)
	Stop()

	require.Len(t, l1.warnCalls, 1)
	require.Len(t, l2.warnCalls, 1)
	assert.Equal(t, "[WARNING] slow response", l1.warnCalls[0].message)
	assert.Equal(t, "[WARNING] slow response", l2.warnCalls[0].message)
	assert.Equal(t, "poll", l2.warnCalls[0].context["task"])
}

func TestError_AddsFallbackStackTraceAndDoesNotMutateInputContext(t *testing.T) {
	l1 := &fakeLogger{}
	l2 := &fakeLogger{}
	Init(testSetting{l1}, testSetting{l2})

	inputCtx := map[string]interface{}{"error": errors.New("plain error"), "taskType": "poll"}
	Error("failed", inputCtx)
	Stop()

	require.Len(t, l1.errorCalls, 1)
	require.Len(t, l2.errorCalls, 1)
	assert.Equal(t, "[ERROR] failed", l1.errorCalls[0].message)
	assert.Equal(t, "poll", l1.errorCalls[0].context["taskType"])
	assert.Contains(t, l1.errorCalls[0].context, "stackTrace")
	assert.NotEmpty(t, l1.errorCalls[0].context["stackTrace"])

	_, hasStackInOriginal := inputCtx["stackTrace"]
	assert.False(t, hasStackInOriginal, "original context must not be mutated")
}

func TestFatal_FanOutAndPrefix(t *testing.T) {
	l1 := &fakeLogger{}
	l2 := &fakeLogger{}
	Init(testSetting{l1}, testSetting{l2})

	ctx := map[string]interface{}{"service": "scheduler"}
	Fatal("critical failure", ctx)
	Stop()

	require.Len(t, l1.fatalCalls, 1)
	require.Len(t, l2.fatalCalls, 1)
	assert.Equal(t, "[FATAL] critical failure", l1.fatalCalls[0].message)
	assert.Equal(t, "scheduler", l1.fatalCalls[0].context["service"])
	assert.Contains(t, l1.fatalCalls[0].context, "stackTrace")
}

func TestRecover_LogsPanicAsFatal(t *testing.T) {
	l := &fakeLogger{}
	Init(testSetting{l})

	func() {
		defer Recover(map[string]interface{}{"component": "test"})
		panic("something went wrong")
	}()
	Stop()

	require.Len(t, l.fatalCalls, 1)
	assert.Equal(t, "[FATAL] Panic recovered", l.fatalCalls[0].message)
	assert.Equal(t, "something went wrong", l.fatalCalls[0].context["panic"])
	assert.Equal(t, "test", l.fatalCalls[0].context["ctx"].(map[string]interface{})["component"])
	assert.Contains(t, l.fatalCalls[0].context, "stackTrace")
}

func TestError_UsesExistingStackTraceFromContext(t *testing.T) {
	l := &fakeLogger{}
	Init(testSetting{l})

	Error("failed", map[string]interface{}{
		"error":      errors.New("plain"),
		"stackTrace": "precomputed-stack",
	})
	Stop()

	require.Len(t, l.errorCalls, 1)
	assert.Equal(t, "precomputed-stack", l.errorCalls[0].context["stackTrace"])
}

func TestError_UsesEmbeddedStackTraceFromError(t *testing.T) {
	l := &fakeLogger{}
	Init(testSetting{l})

	Error("failed", map[string]interface{}{
		"error": stackAwareError{
			message: "wrapped error",
			stack:   "embedded-stack-line-1\nembedded-stack-line-2",
		},
	})
	Stop()

	require.Len(t, l.errorCalls, 1)
	assert.Equal(t, "embedded-stack-line-1\nembedded-stack-line-2", l.errorCalls[0].context["stackTrace"])
}

func TestInit_Empty(t *testing.T) {
	// Should not panic even if no loggers are initialized
	Init()
	assert.NotPanics(t, func() {
		Info("test", nil)
		Error("test error", nil)
	})
	Stop()
}

func TestRecover_NilContext(t *testing.T) {
	l := &fakeLogger{}
	Init(testSetting{l})

	assert.NotPanics(t, func() {
		defer Recover(nil)
		panic("panic with nil ctx")
	})
	Stop()

	require.Len(t, l.fatalCalls, 1)
	assert.Equal(t, "[FATAL] Panic recovered", l.fatalCalls[0].message)
}

func TestError_WorksWithNilContext(t *testing.T) {
	l := &fakeLogger{}
	Init(testSetting{l})

	Error("failed", nil)
	Stop()

	require.Len(t, l.errorCalls, 1)
	assert.Equal(t, "[ERROR] failed", l.errorCalls[0].message)
	assert.Contains(t, l.errorCalls[0].context, "stackTrace")
}

func TestInit_Reinitialization(t *testing.T) {
	l1 := &fakeLogger{}
	Init(testSetting{l1})
	Info("first", nil)

	l2 := &fakeLogger{}
	// Re-initializing should flush and close l1's worker before starting instance with l2
	Init(testSetting{l2})
	Info("second", nil)
	Stop()

	assert.Len(t, l1.infoCalls, 1)
	assert.Equal(t, "[INFO] first", l1.infoCalls[0].message)
	assert.Len(t, l2.infoCalls, 1)
	assert.Equal(t, "[INFO] second", l2.infoCalls[0].message)
}

func TestLogging_BeforeInitOrAfterStop(t *testing.T) {
	// Ensure instance is nil
	Stop()

	assert.NotPanics(t, func() {
		Info("should not panic", nil)
		Warn("should not panic", nil)
		Error("should not panic", nil)
		Fatal("should not panic", nil)
	})
}

func TestBuildErrorContextWithStackTrace_UsesFallbackWhenNoEmbeddedStack(t *testing.T) {
	ctx := map[string]interface{}{
		"error": errors.New("plain error"),
	}

	result := internal.BuildErrorContextWithStackTrace(ctx)

	assert.NotEmpty(t, result["stackTrace"])
	assert.Contains(t, result["stackTrace"], ".go:")
}

func TestBuildErrorContextWithStackTrace_UsesEmbeddedStackWhenPresent(t *testing.T) {
	ctx := map[string]interface{}{
		"error": stackAwareError{
			message: "wrapped error",
			stack:   "embedded-stack-line-1\nembedded-stack-line-2",
		},
	}

	result := internal.BuildErrorContextWithStackTrace(ctx)

	assert.Equal(t, "embedded-stack-line-1\nembedded-stack-line-2", result["stackTrace"])
}

func TestBuildErrorContextWithStackTrace_DoesNotOverrideExistingStack(t *testing.T) {
	ctx := map[string]interface{}{
		"error":      errors.New("plain"),
		"stackTrace": "precomputed-stack",
	}

	result := internal.BuildErrorContextWithStackTrace(ctx)

	assert.Equal(t, "precomputed-stack", result["stackTrace"])
}
