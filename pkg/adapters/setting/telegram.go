package setting

import (
	"net/http"
	"time"

	"github.com/Consolushka/golang.composite_logger/internal/adapters/logger"
	compositelogger "github.com/Consolushka/golang.composite_logger/pkg"
	"github.com/Consolushka/golang.composite_logger/pkg/ports"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// TelegramSetting provides configuration for the Telegram logging adapter.
type TelegramSetting struct {
	// Enabled toggles the telegram logger on or off.
	Enabled              bool
	// BotKey is the Telegram bot API token.
	BotKey               string
	// ChatId is the unique identifier for the target chat or user.
	ChatId               int64
	// Timeout sets the HTTP client timeout for API requests.
	Timeout              time.Duration
	// LowerLevel sets the minimum severity level to log.
	LowerLevel           compositelogger.Level
	// UseLevelTitleWrapper enables emoji/symbol decoration around log levels if true (default: true).
	UseLevelTitleWrapper *bool
	// LevelWrappers allows overriding default emojis for specific levels.
	LevelWrappers        map[compositelogger.Level]string
	// LevelTitles allows overriding default display names for log levels.
	LevelTitles          map[compositelogger.Level]string
}

var botAPIConstructor = tgbotapi.NewBotAPI

// InitLogger initializes a Telegram-based logger with MarkdownV2 support.
func (t TelegramSetting) InitLogger() ports.Logger {
	botApi, err := botAPIConstructor(t.BotKey)
	if err != nil {
		panic("Error creating telegram bot api. Error: " + err.Error())
	}

	if t.Timeout > 0 {
		botApi.Client = &http.Client{
			Timeout: t.Timeout,
		}
	}

	useLevelTitleWrapper := true
	if t.UseLevelTitleWrapper != nil {
		useLevelTitleWrapper = *t.UseLevelTitleWrapper
	}

	finalWrappers := make(map[compositelogger.Level]string)
	if useLevelTitleWrapper {
		// First, fill with default wrappers
		for level, wrapper := range compositelogger.DefaultLevelWrappers {
			finalWrappers[level] = wrapper
		}
		// Then override with user-defined wrappers
		for level, wrapper := range t.LevelWrappers {
			if wrapper != "" {
				finalWrappers[level] = wrapper
			}
		}
	}

	return logger.TelegramLogger{
		BotApi:               botApi,
		LogChatId:            t.ChatId,
		Level:                t.LowerLevel,
		UseLevelTitleWrapper: useLevelTitleWrapper,
		LevelWrappers:        finalWrappers,
		LevelTitles:          t.LevelTitles,
	}
}

// IsEnabled returns the current active status of the adapter.
func (t TelegramSetting) IsEnabled() bool {
	return t.Enabled
}
