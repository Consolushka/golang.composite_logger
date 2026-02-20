package setting

import (
	composite_logger "composite_logger/pkg"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConsoleSetting_InitLogger(t *testing.T) {
	s := ConsoleSetting{
		LowerLevel: composite_logger.InfoLevel,
	}

	assert.NotPanics(t, func() {
		l := s.InitLogger()
		assert.NotNil(t, l)
	})
}
