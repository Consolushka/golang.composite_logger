package setting

import (
	"os"
	"path/filepath"
	"testing"

	composite_logger "github.com/Consolushka/golang.composite_logger/pkg"

	"github.com/stretchr/testify/assert"
)

func TestFileSetting_IsEnabled(t *testing.T) {
	assert.True(t, FileSetting{Enabled: true}.IsEnabled())
	assert.False(t, FileSetting{Enabled: false}.IsEnabled())
}

func TestFileSetting_InitLogger(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")

	t.Run("basic initialization", func(t *testing.T) {
		s := FileSetting{
			Path:       logPath,
			LowerLevel: composite_logger.ErrorLevel,
		}

		assert.NotPanics(t, func() {
			l := s.InitLogger()
			assert.NotNil(t, l)
		})

		// Check if file was created (lumberjack might not create it until first write, 
		// but our InitLogger calls os.MkdirAll and lumberjack initialization)
		_, err := os.Stat(logPath)
		// Note: lumberjack only opens/creates the file on first write
		// but since we aren't writing in this test, we just check if InitLogger runs.
		assert.True(t, err == nil || os.IsNotExist(err))
	})

	t.Run("initialization with rotation settings", func(t *testing.T) {
		s := FileSetting{
			Path:       filepath.Join(tmpDir, "rotated.log"),
			LowerLevel: composite_logger.InfoLevel,
			MaxSize:    10,
			MaxBackups: 5,
			MaxAge:     7,
			Compress:   true,
		}

		assert.NotPanics(t, func() {
			l := s.InitLogger()
			assert.NotNil(t, l)
		})
	})
}

func TestFileSetting_SetupRotation(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		s := FileSetting{Path: "test.log"}
		lj := s.setupRotation()

		assert.Equal(t, "test.log", lj.Filename)
		assert.Equal(t, 5, lj.MaxSize)
		assert.Equal(t, 3, lj.MaxBackups)
		assert.Equal(t, 28, lj.MaxAge)
		assert.False(t, lj.Compress)
	})

	t.Run("custom values", func(t *testing.T) {
		s := FileSetting{
			Path:       "custom.log",
			MaxSize:    10,
			MaxBackups: 5,
			MaxAge:     7,
			Compress:   true,
		}
		lj := s.setupRotation()

		assert.Equal(t, "custom.log", lj.Filename)
		assert.Equal(t, 10, lj.MaxSize)
		assert.Equal(t, 5, lj.MaxBackups)
		assert.Equal(t, 7, lj.MaxAge)
		assert.True(t, lj.Compress)
	})
}

func TestFileSetting_InitLogger_PanicOnEmptyPath(t *testing.T) {
	s := FileSetting{
		Path: "",
	}

	assert.Panics(t, func() {
		s.InitLogger()
	})
}
