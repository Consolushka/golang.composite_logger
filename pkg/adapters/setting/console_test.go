package setting

import (
	"testing"

	composite_logger "github.com/Consolushka/golang.composite_logger/pkg"

	"github.com/stretchr/testify/assert"
)

func TestConsoleSetting_IsEnabled(t *testing.T) {
	assert.True(t, ConsoleSetting{Enabled: true}.IsEnabled())
	assert.False(t, ConsoleSetting{Enabled: false}.IsEnabled())
}

func TestConsoleSetting_InitLogger(t *testing.T) {
	s := ConsoleSetting{
		LowerLevel: composite_logger.InfoLevel,
	}

	assert.NotPanics(t, func() {
		l := s.InitLogger()
		assert.NotNil(t, l)
	})
}
