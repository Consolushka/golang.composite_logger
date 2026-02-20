package composite_logger

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestLevel_String(t *testing.T) {
	tests := []struct {
		level    Level
		expected string
	}{
		{InfoLevel, "info"},
		{WarningLevel, "warning"},
		{ErrorLevel, "error"},
		{FatalLevel, "fatal"},
		{Level(99), "info"},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, tt.level.String())
	}
}

func TestLevel_ToLogrus(t *testing.T) {
	tests := []struct {
		level    Level
		expected logrus.Level
	}{
		{InfoLevel, logrus.InfoLevel},
		{WarningLevel, logrus.WarnLevel},
		{ErrorLevel, logrus.ErrorLevel},
		{FatalLevel, logrus.FatalLevel},
		{Level(99), logrus.InfoLevel},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, tt.level.ToLogrus())
	}
}
