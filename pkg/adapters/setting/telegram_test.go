package setting

import (
	"testing"

	composite_logger "github.com/Consolushka/golang.composite_logger/pkg"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
)

func TestTelegramSetting_InitLogger(t *testing.T) {
	// Mock the bot constructor
	oldConstructor := botAPIConstructor
	defer func() { botAPIConstructor = oldConstructor }()

	botAPIConstructor = func(token string) (*tgbotapi.BotAPI, error) {
		return &tgbotapi.BotAPI{}, nil
	}

	s := TelegramSetting{
		Enabled:              true,
		BotKey:               "fake-token",
		ChatId:               12345,
		LowerLevel:           composite_logger.InfoLevel,
		UseLevelTitleWrapper: true,
		LevelWrappers: map[composite_logger.Level]string{
			composite_logger.ErrorLevel: "🔥",
		},
	}

	assert.NotPanics(t, func() {
		l := s.InitLogger()
		assert.NotNil(t, l)
	})
}
