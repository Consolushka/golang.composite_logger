package composite_logger

import (
	"context"
	"sync"

	"github.com/Consolushka/golang.composite_logger/pkg/ports"
)

type logEntry struct {
	ctx     context.Context
	level   int
	message string
	fields  map[string]interface{}
}

// AsyncCompositeLogger implements the ports.CompositeLogger interface for non-blocking log dispatch using a background worker.
type AsyncCompositeLogger struct {
	loggers     []ports.Logger
	hooks       []ports.Hook
	contextKeys []any
	mu          sync.RWMutex
	ch          chan logEntry
	wg          sync.WaitGroup

	// stateMu is separate from mu: the worker holds mu while draining the
	// channel, so guarding the channel with the same mutex could deadlock
	// senders blocked on a full buffer during Stop.
	stateMu sync.RWMutex
	closed  bool
}

// NewAsyncCompositeLogger creates a new asynchronous composite logger and starts its background worker.
func NewAsyncCompositeLogger(loggers []ports.Logger, bufferSize int) *AsyncCompositeLogger {
	a := &AsyncCompositeLogger{
		loggers: loggers,
		ch:      make(chan logEntry, bufferSize),
	}
	a.wg.Add(1)
	go a.listenAndBroadcast()
	return a
}

// Log queues a log entry for asynchronous dispatch.
func (a *AsyncCompositeLogger) Log(level int, message string, fields map[string]interface{}) {
	a.enqueue(logEntry{ctx: context.Background(), level: level, message: message, fields: fields})
}

// LogContext queues a log entry with calling context for asynchronous dispatch.
func (a *AsyncCompositeLogger) LogContext(ctx context.Context, level int, message string, fields map[string]interface{}) {
	if ctx == nil {
		ctx = context.Background()
	}
	// The worker processes the entry after the caller's request may have
	// finished. Detach cancellation so the log survives the request's
	// lifetime while context values stay available for enrichment.
	a.enqueue(logEntry{ctx: context.WithoutCancel(ctx), level: level, message: message, fields: fields})
}

func (a *AsyncCompositeLogger) enqueue(entry logEntry) {
	a.stateMu.RLock()
	defer a.stateMu.RUnlock()
	if a.closed {
		return
	}
	a.ch <- entry
}

// AddHook registers a new hook to the logger instance.
func (a *AsyncCompositeLogger) AddHook(hook ports.Hook) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.hooks = append(a.hooks, hook)
}

// SetContextKeys registers context keys for automatic log enrichment.
func (a *AsyncCompositeLogger) SetContextKeys(keys ...any) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.contextKeys = keys
}

// Stop closes the log channel and waits for the worker to finish processing queued logs.
func (a *AsyncCompositeLogger) Stop() {
	a.stateMu.Lock()
	if a.closed {
		a.stateMu.Unlock()
		return
	}
	a.closed = true
	close(a.ch)
	a.stateMu.Unlock()

	a.wg.Wait()
}

func (a *AsyncCompositeLogger) listenAndBroadcast() {
	defer a.wg.Done()
	for entry := range a.ch {
		a.mu.RLock()
		hooks := a.hooks
		contextKeys := a.contextKeys
		loggers := a.loggers
		a.mu.RUnlock()

		broadcast(entry.ctx, entry.level, entry.message, entry.fields, loggers, hooks, contextKeys)
	}
}
