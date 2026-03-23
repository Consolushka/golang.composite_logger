package composite_logger

import (
	"context"
	"sync"

	"github.com/Consolushka/golang.composite_logger/pkg/ports"
)

// SyncCompositeLogger implements the ports.CompositeLogger interface for immediate, blocking log dispatch.
type SyncCompositeLogger struct {
	loggers     []ports.Logger
	hooks       []ports.Hook
	contextKeys []any
	mu          sync.RWMutex
}

// NewSyncCompositeLogger creates a new synchronous composite logger instance.
func NewSyncCompositeLogger(loggers []ports.Logger) *SyncCompositeLogger {
	return &SyncCompositeLogger{
		loggers: loggers,
	}
}

// Log dispatches a log entry to all registered adapters synchronously.
func (s *SyncCompositeLogger) Log(level int, message string, fields map[string]interface{}) {
	s.broadcast(context.Background(), level, message, fields)
}

// LogContext dispatches a log entry with calling context to all registered adapters synchronously.
func (s *SyncCompositeLogger) LogContext(ctx context.Context, level int, message string, fields map[string]interface{}) {
	s.broadcast(ctx, level, message, fields)
}

// AddHook registers a new hook to the logger instance.
func (s *SyncCompositeLogger) AddHook(hook ports.Hook) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.hooks = append(s.hooks, hook)
}

// SetContextKeys registers context keys for automatic log enrichment.
func (s *SyncCompositeLogger) SetContextKeys(keys ...any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.contextKeys = keys
}

// Stop for SyncCompositeLogger is a no-op as there is no background processing.
func (s *SyncCompositeLogger) Stop() {}

func (s *SyncCompositeLogger) broadcast(ctx context.Context, level int, message string, fields map[string]interface{}) {
	s.mu.RLock()
	hooks := s.hooks
	contextKeys := s.contextKeys
	loggers := s.loggers
	s.mu.RUnlock()

	broadcast(ctx, level, message, fields, loggers, hooks, contextKeys)
}
