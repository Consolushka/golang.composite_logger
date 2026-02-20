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
	UseLevelTitleWrapper bool
	LevelWrappers        map[compositelogger.Level]string
	LevelTitles          map[compositelogger.Level]string
}

var botAPIConstructor = tgbotapi.NewBotAPI

func (t TelegramSetting) InitLogger() ports.Logger {
	botApi, err := botAPIConstructor(t.BotKey)
	if err != nil {
		panic("Error creating telegram bot api. Error: " + err.Error())
	}

	finalWrappers := make(map[compositelogger.Level]string)
	if t.UseLevelTitleWrapper {
		// Сначала заполняем дефолтными
		for level, wrapper := range compositelogger.DefaultLevelWrappers {
			finalWrappers[level] = wrapper
		}
		// Перекрываем пользовательскими
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
		UseLevelTitleWrapper: t.UseLevelTitleWrapper,
		LevelWrappers:        finalWrappers,
		LevelTitles:          t.LevelTitles,
	}
}
