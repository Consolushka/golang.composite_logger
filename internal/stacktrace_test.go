package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCloneContext(t *testing.T) {
	t.Run("Nil context returns empty map", func(t *testing.T) {
		assert.NotNil(t, CloneContext(nil))
		assert.Empty(t, CloneContext(nil))
	})

	t.Run("Clones existing values", func(t *testing.T) {
		original := map[string]interface{}{"key": "value", "id": 123}
		cloned := CloneContext(original)

		assert.Equal(t, original, cloned)
		assert.NotSame(t, &original, &cloned) // Different map instances

		cloned["new"] = "data"
		assert.NotContains(t, original, "new")
	})
}

func TestShouldIncludeFrame(t *testing.T) {
	tests := []struct {
		function string
		include  bool
	}{
		{"main.main", true},
		{"net/http.HandlerFunc.ServeHTTP", true},
		{"github.com/Consolushka/golang.composite_logger/pkg/Info", false},
		{"github.com/Consolushka/golang.composite_logger/pkg/ports/Setting.InitLogger", false},
		{"github.com/Consolushka/golang.composite_logger/internal.BuildErrorContextWithStackTrace", false},
		{"github.com/Consolushka/golang.composite_logger/internal/adapters/logger.TelegramLogger.send", false},
	}

	for _, tt := range tests {
		t.Run(tt.function, func(t *testing.T) {
			assert.Equal(t, tt.include, ShouldIncludeFrame(tt.function))
		})
	}
}
