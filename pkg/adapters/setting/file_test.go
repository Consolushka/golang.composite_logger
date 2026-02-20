package setting

import (
	"os"
	"path/filepath"
	"testing"

	composite_logger "github.com/Consolushka/golang.composite_logger/pkg"

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
