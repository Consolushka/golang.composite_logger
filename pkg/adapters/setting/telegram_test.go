package setting

import (
	"errors"
	"testing"

	"github.com/Consolushka/golang.composite_logger/internal/adapters/logger"
	composite_logger "github.com/Consolushka/golang.composite_logger/pkg"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
)

func TestTelegramSetting_IsEnabled(t *testing.T) {
	assert.True(t, TelegramSetting{Enabled: true}.IsEnabled())
	assert.False(t, TelegramSetting{Enabled: false}.IsEnabled())
}

func TestTelegramSetting_InitLogger(t *testing.T) {
	// Mock the bot constructor
	oldConstructor := botAPIConstructor
	defer func() { botAPIConstructor = oldConstructor }()

	botAPIConstructor = func(token string) (*tgbotapi.BotAPI, error) {
		return &tgbotapi.BotAPI{}, nil
	}

	t.Run("success with default wrappers", func(t *testing.T) {
		s := TelegramSetting{
			Enabled:    true,
			BotKey:     "fake-token",
			ChatId:     12345,
			LowerLevel: composite_logger.InfoLevel,
		}

		l := s.InitLogger()
		tgLogger, ok := l.(logger.TelegramLogger)
		assert.True(t, ok)
		assert.True(t, tgLogger.UseLevelTitleWrapper)
		assert.Equal(t, composite_logger.DefaultLevelWrappers[composite_logger.InfoLevel], tgLogger.LevelWrappers[composite_logger.InfoLevel])
	})

	t.Run("success with custom wrappers overrides", func(t *testing.T) {
		trueVal := true
		s := TelegramSetting{
			BotKey:               "fake-token",
			UseLevelTitleWrapper: &trueVal,
			LevelWrappers: map[composite_logger.Level]string{
				composite_logger.InfoLevel: "CUSTOM",
			},
		}

		l := s.InitLogger()
		tgLogger := l.(logger.TelegramLogger)
		assert.Equal(t, "CUSTOM", tgLogger.LevelWrappers[composite_logger.InfoLevel])
		// Error level should still have default
		assert.Equal(t, composite_logger.DefaultLevelWrappers[composite_logger.ErrorLevel], tgLogger.LevelWrappers[composite_logger.ErrorLevel])
	})

	t.Run("disabled wrappers", func(t *testing.T) {
		falseVal := false
		s := TelegramSetting{
			BotKey:               "fake-token",
			UseLevelTitleWrapper: &falseVal,
		}

		l := s.InitLogger()
		tgLogger := l.(logger.TelegramLogger)
		assert.False(t, tgLogger.UseLevelTitleWrapper)
		assert.Empty(t, tgLogger.LevelWrappers)
	})

	t.Run("panic on constructor error", func(t *testing.T) {
		botAPIConstructor = func(token string) (*tgbotapi.BotAPI, error) {
			return nil, errors.New("api error")
		}

		assert.Panics(t, func() {
			TelegramSetting{BotKey: "wrong"}.InitLogger()
		})
	})
}
