package setting

import (
	"composite_logger/internal/adapters/logger"
	composite_logger "composite_logger/pkg"
	"composite_logger/pkg/ports"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramSetting struct {
	Enabled              bool
	BotKey               string
	ChatId               int64
	LowerLevel           composite_logger.Level
	UseLevelTitleWrapper bool
	LevelWrappers        map[composite_logger.Level]string
	LevelTitles          map[composite_logger.Level]string
}

var botAPIConstructor = tgbotapi.NewBotAPI

func (t TelegramSetting) InitLogger() ports.Logger {
	botApi, err := botAPIConstructor(t.BotKey)
	if err != nil {
		panic("Error creating telegram bot api. Error: " + err.Error())
	}

	finalWrappers := make(map[composite_logger.Level]string)
	if t.UseLevelTitleWrapper {
		// Сначала заполняем дефолтными
		for level, wrapper := range composite_logger.DefaultLevelWrappers {
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
