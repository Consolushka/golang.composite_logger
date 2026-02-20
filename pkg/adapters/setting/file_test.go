package setting

import (
	composite_logger "composite_logger/pkg"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileSetting_InitLogger(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")

	s := FileSetting{
		Path:       logPath,
		LowerLevel: composite_logger.ErrorLevel,
	}

	assert.NotPanics(t, func() {
		l := s.InitLogger()
		assert.NotNil(t, l)
	})

	// Check if file was created
	_, err := os.Stat(logPath)
	assert.NoError(t, err)
}

func TestFileSetting_InitLogger_PanicOnEmptyPath(t *testing.T) {
	s := FileSetting{
		Path: "",
	}

	assert.Panics(t, func() {
		s.InitLogger()
	})
}
