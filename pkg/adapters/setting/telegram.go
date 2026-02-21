package setting

import (
	"github.com/Consolushka/golang.composite_logger/internal/adapters/logger"
	compositelogger "github.com/Consolushka/golang.composite_logger/pkg"
	"github.com/Consolushka/golang.composite_logger/pkg/ports"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramSetting struct {
	Enabled              bool
	BotKey               string
	ChatId               int64
	LowerLevel           compositelogger.Level
	UseLevelTitleWrapper *bool
	LevelWrappers        map[compositelogger.Level]string
	LevelTitles          map[compositelogger.Level]string
}

var botAPIConstructor = tgbotapi.NewBotAPI

func (t TelegramSetting) InitLogger() ports.Logger {
	botApi, err := botAPIConstructor(t.BotKey)
	if err != nil {
		panic("Error creating telegram bot api. Error: " + err.Error())
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

func (t TelegramSetting) IsEnabled() bool {
	return t.Enabled
}
