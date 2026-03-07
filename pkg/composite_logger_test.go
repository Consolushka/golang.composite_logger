package composite_logger

import (
	"context"
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
	fields  map[string]interface{}
	ctx     context.Context
}

type fakeLogger struct {
	infoCalls  []logCall
	warnCalls  []logCall
	errorCalls []logCall
	fatalCalls []logCall
}

func (f *fakeLogger) Info(message string, fields map[string]interface{}) {
	f.infoCalls = append(f.infoCalls, logCall{message: message, fields: fields})
}

func (f *fakeLogger) InfoContext(ctx context.Context, message string, fields map[string]interface{}) {
	f.infoCalls = append(f.infoCalls, logCall{message: message, fields: fields, ctx: ctx})
}

func (f *fakeLogger) Warn(message string, fields map[string]interface{}) {
	f.warnCalls = append(f.warnCalls, logCall{message: message, fields: fields})
}

func (f *fakeLogger) WarnContext(ctx context.Context, message string, fields map[string]interface{}) {
	f.warnCalls = append(f.warnCalls, logCall{message: message, fields: fields, ctx: ctx})
}

func (f *fakeLogger) Error(message string, fields map[string]interface{}) {
	f.errorCalls = append(f.errorCalls, logCall{message: message, fields: fields})
}

func (f *fakeLogger) ErrorContext(ctx context.Context, message string, fields map[string]interface{}) {
	f.errorCalls = append(f.errorCalls, logCall{message: message, fields: fields, ctx: ctx})
}

func (f *fakeLogger) Fatal(message string, fields map[string]interface{}) {
	f.fatalCalls = append(f.fatalCalls, logCall{message: message, fields: fields})
}

func (f *fakeLogger) FatalContext(ctx context.Context, message string, fields map[string]interface{}) {
	f.fatalCalls = append(f.fatalCalls, logCall{message: message, fields: fields, ctx: ctx})
}

func TestInfo_FanOutAndPrefix(t *testing.T) {
	l1 := &fakeLogger{}
	l2 := &fakeLogger{}
	Init(testSetting{l1}, testSetting{l2})

	fields := map[string]interface{}{"requestId": "abc-123"}
	Info("process started", fields)
	Stop()

	require.Len(t, l1.infoCalls, 1)
	require.Len(t, l2.infoCalls, 1)
	assert.Equal(t, "[INFO] process started", l1.infoCalls[0].message)
	assert.Equal(t, "[INFO] process started", l2.infoCalls[0].message)
	assert.Equal(t, "abc-123", l1.infoCalls[0].fields["requestId"])
}

func TestInfoContext_FanOutAndPrefix(t *testing.T) {
	l1 := &fakeLogger{}
	l2 := &fakeLogger{}
	Init(testSetting{l1}, testSetting{l2})

	type contextKey string
	ctx := context.WithValue(context.Background(), contextKey("traceId"), "trace-abc-123")
	fields := map[string]interface{}{"requestId": "abc-123"}
	InfoContext(ctx, "process started", fields)
	Stop()

	require.Len(t, l1.infoCalls, 1)
	require.Len(t, l2.infoCalls, 1)
	assert.Equal(t, "[INFO] process started", l1.infoCalls[0].message)
	assert.Equal(t, ctx, l1.infoCalls[0].ctx)
	assert.Equal(t, "abc-123", l1.infoCalls[0].fields["requestId"])
}

func TestWithContext_UsesBoundContext(t *testing.T) {
	l := &fakeLogger{}
	Init(testSetting{l})

	type contextKey string
	ctx := context.WithValue(context.Background(), contextKey("traceId"), "trace-with-context")

	lc := WithContext(ctx)
	lc.Info("bound log", map[string]interface{}{"foo": "bar"})
	Stop()

	require.Len(t, l.infoCalls, 1)
	assert.Equal(t, ctx, l.infoCalls[0].ctx)
	assert.Equal(t, "[INFO] bound log", l.infoCalls[0].message)
	assert.Equal(t, "bar", l.infoCalls[0].fields["foo"])
}

func TestSetContextKeys_EnrichesLogs(t *testing.T) {
	l := &fakeLogger{}
	Init(testSetting{l})
	SetContextKeys("trace_id", "user_id")

	ctx := context.WithValue(context.Background(), "trace_id", "abc-123")
	ctx = context.WithValue(ctx, "user_id", 42)

	InfoContext(ctx, "enriched log", nil)
	Stop()

	require.Len(t, l.infoCalls, 1)
	assert.Equal(t, "abc-123", l.infoCalls[0].fields["trace_id"])
	assert.Equal(t, 42, l.infoCalls[0].fields["user_id"])
}

func TestWarn_FanOutAndPrefix(t *testing.T) {
	l1 := &fakeLogger{}
	l2 := &fakeLogger{}
	Init(testSetting{l1}, testSetting{l2})

	fields := map[string]interface{}{"task": "poll"}
	Warn("slow response", fields)
	Stop()

	require.Len(t, l1.warnCalls, 1)
	require.Len(t, l2.warnCalls, 1)
	assert.Equal(t, "[WARNING] slow response", l1.warnCalls[0].message)
	assert.Equal(t, "[WARNING] slow response", l2.warnCalls[0].message)
	assert.Equal(t, "poll", l2.warnCalls[0].fields["task"])
}

func TestError_AddsFallbackStackTraceAndDoesNotMutateInputContext(t *testing.T) {
	l1 := &fakeLogger{}
	l2 := &fakeLogger{}
	Init(testSetting{l1}, testSetting{l2})

	inputFields := map[string]interface{}{"error": errors.New("plain error"), "taskType": "poll"}
	Error("failed", inputFields)
	Stop()

	require.Len(t, l1.errorCalls, 1)
	require.Len(t, l2.errorCalls, 1)
	assert.Equal(t, "[ERROR] failed", l1.errorCalls[0].message)
	assert.Equal(t, "poll", l1.errorCalls[0].fields["taskType"])
	assert.Contains(t, l1.errorCalls[0].fields, "stackTrace")
	assert.NotEmpty(t, l1.errorCalls[0].fields["stackTrace"])

	_, hasStackInOriginal := inputFields["stackTrace"]
	assert.False(t, hasStackInOriginal, "original fields must not be mutated")
}

func TestFatal_FanOutAndPrefix(t *testing.T) {
	l1 := &fakeLogger{}
	l2 := &fakeLogger{}
	Init(testSetting{l1}, testSetting{l2})

	fields := map[string]interface{}{"service": "scheduler"}
	Fatal("critical failure", fields)
	Stop()

	require.Len(t, l1.fatalCalls, 1)
	require.Len(t, l2.fatalCalls, 1)
	assert.Equal(t, "[FATAL] critical failure", l1.fatalCalls[0].message)
	assert.Equal(t, "scheduler", l1.fatalCalls[0].fields["service"])
	assert.Contains(t, l1.fatalCalls[0].fields, "stackTrace")
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
	assert.Equal(t, "something went wrong", l.fatalCalls[0].fields["panic"])
	assert.Equal(t, "test", l.fatalCalls[0].fields["fields"].(map[string]interface{})["component"])
	assert.Contains(t, l.fatalCalls[0].fields, "stackTrace")
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
	assert.Equal(t, "precomputed-stack", l.errorCalls[0].fields["stackTrace"])
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
	assert.Equal(t, "embedded-stack-line-1\nembedded-stack-line-2", l.errorCalls[0].fields["stackTrace"])
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
		panic("panic with nil fields")
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
	assert.Contains(t, l.errorCalls[0].fields, "stackTrace")
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
	fields := map[string]interface{}{
		"error": errors.New("plain error"),
	}

	result := internal.BuildErrorContextWithStackTrace(fields)

	assert.NotEmpty(t, result["stackTrace"])
	assert.Contains(t, result["stackTrace"], ".go:")
}

func TestBuildErrorContextWithStackTrace_UsesEmbeddedStackWhenPresent(t *testing.T) {
	fields := map[string]interface{}{
		"error": stackAwareError{
			message: "wrapped error",
			stack:   "embedded-stack-line-1\nembedded-stack-line-2",
		},
	}

	result := internal.BuildErrorContextWithStackTrace(fields)

	assert.Equal(t, "embedded-stack-line-1\nembedded-stack-line-2", result["stackTrace"])
}

func TestBuildErrorContextWithStackTrace_DoesNotOverrideExistingStack(t *testing.T) {
	fields := map[string]interface{}{
		"error":      errors.New("plain"),
		"stackTrace": "precomputed-stack",
	}

	result := internal.BuildErrorContextWithStackTrace(fields)

	assert.Equal(t, "precomputed-stack", result["stackTrace"])
}
