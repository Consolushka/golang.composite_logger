package composite_logger

import (
	"context"
	"sync"
	"testing"

	"github.com/Consolushka/golang.composite_logger/pkg/ports"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type recordedCall struct {
	ctx     context.Context
	message string
	fields  map[string]interface{}
}

type recordingLogger struct {
	mu    sync.Mutex
	calls []recordedCall
}

func (r *recordingLogger) record(ctx context.Context, message string, fields map[string]interface{}) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.calls = append(r.calls, recordedCall{ctx: ctx, message: message, fields: fields})
}

func (r *recordingLogger) snapshot() []recordedCall {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]recordedCall(nil), r.calls...)
}

func (r *recordingLogger) Info(message string, fields map[string]interface{}) {
	r.record(nil, message, fields)
}

func (r *recordingLogger) InfoContext(ctx context.Context, message string, fields map[string]interface{}) {
	r.record(ctx, message, fields)
}

func (r *recordingLogger) Warn(message string, fields map[string]interface{}) {
	r.record(nil, message, fields)
}

func (r *recordingLogger) WarnContext(ctx context.Context, message string, fields map[string]interface{}) {
	r.record(ctx, message, fields)
}

func (r *recordingLogger) Error(message string, fields map[string]interface{}) {
	r.record(nil, message, fields)
}

func (r *recordingLogger) ErrorContext(ctx context.Context, message string, fields map[string]interface{}) {
	r.record(ctx, message, fields)
}

func (r *recordingLogger) Fatal(message string, fields map[string]interface{}) {
	r.record(nil, message, fields)
}

func (r *recordingLogger) FatalContext(ctx context.Context, message string, fields map[string]interface{}) {
	r.record(ctx, message, fields)
}

// Regression: in async mode the worker processes entries after the handler
// returns, when the request context is already cancelled. The entry must
// still reach the adapters, and context values must survive for enrichment.
func TestAsync_DeliversEntryWithCancelledContext(t *testing.T) {
	l := &recordingLogger{}
	a := NewAsyncCompositeLogger([]ports.Logger{l}, 10)

	type contextKey string
	traceKey := contextKey("trace_id")
	a.SetContextKeys(traceKey)

	ctx, cancel := context.WithCancel(context.WithValue(context.Background(), traceKey, "trace-123"))

	a.LogContext(ctx, 3, "timeout error", nil)
	cancel()
	a.Stop()

	calls := l.snapshot()
	require.Len(t, calls, 1)
	assert.Equal(t, "timeout error", calls[0].message)
	assert.Equal(t, "trace-123", calls[0].fields["trace_id"], "context values must survive detachment")
	assert.NoError(t, calls[0].ctx.Err(), "cancellation must not propagate to the queued entry")
}

func TestAsync_LogContextWithNilContext_NoPanic(t *testing.T) {
	l := &recordingLogger{}
	a := NewAsyncCompositeLogger([]ports.Logger{l}, 10)

	var nilCtx context.Context
	assert.NotPanics(t, func() {
		a.LogContext(nilCtx, 1, "nil ctx", nil)
	})
	a.Stop()

	require.Len(t, l.snapshot(), 1)
}

func TestAsync_LogAfterStop_NoPanicAndNoDelivery(t *testing.T) {
	l := &recordingLogger{}
	a := NewAsyncCompositeLogger([]ports.Logger{l}, 10)
	a.Stop()

	assert.NotPanics(t, func() {
		a.Log(1, "after stop", nil)
		a.LogContext(context.Background(), 1, "after stop ctx", nil)
		a.Stop()
	})
	assert.Empty(t, l.snapshot())
}

// Regression: Stop used to close the channel without synchronizing with
// senders — a concurrent Log caused "send on closed channel".
func TestAsync_ConcurrentLogAndStop_NoPanic(t *testing.T) {
	for i := 0; i < 100; i++ {
		l := &recordingLogger{}
		a := NewAsyncCompositeLogger([]ports.Logger{l}, 1)

		var wg sync.WaitGroup
		for g := 0; g < 8; g++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				assert.NotPanics(t, func() {
					for j := 0; j < 50; j++ {
						a.Log(1, "concurrent", nil)
					}
				})
			}()
		}
		a.Stop()
		wg.Wait()
	}
}
