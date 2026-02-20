package composite_logger

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestParseLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected Level
		wantErr  bool
	}{
		{"info", InfoLevel, false},
		{"INFO", InfoLevel, false},
		{"warn", WarningLevel, false},
		{"warning", WarningLevel, false},
		{"error", ErrorLevel, false},
		{"fatal", FatalLevel, false},
		{"invalid", InfoLevel, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseLevel(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}
