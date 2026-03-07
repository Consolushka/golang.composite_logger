package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	composite_logger "github.com/Consolushka/golang.composite_logger/pkg"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramLogger struct {
	BotApi               *tgbotapi.BotAPI
	LogChatId            int64
	Level                composite_logger.Level
	UseLevelTitleWrapper bool
	LevelWrappers        map[composite_logger.Level]string
	LevelTitles          map[composite_logger.Level]string
}

func (t TelegramLogger) Info(message string, fields map[string]interface{}) {
	t.send(context.TODO(), message, fields, composite_logger.InfoLevel)
}

func (t TelegramLogger) InfoContext(ctx context.Context, message string, fields map[string]interface{}) {
	t.send(ctx, message, fields, composite_logger.InfoLevel)
}

func (t TelegramLogger) Warn(message string, fields map[string]interface{}) {
	t.send(context.TODO(), message, fields, composite_logger.WarningLevel)
}

func (t TelegramLogger) WarnContext(ctx context.Context, message string, fields map[string]interface{}) {
	t.send(ctx, message, fields, composite_logger.WarningLevel)
}

func (t TelegramLogger) Error(message string, fields map[string]interface{}) {
	t.send(context.TODO(), message, fields, composite_logger.ErrorLevel)
}

func (t TelegramLogger) ErrorContext(ctx context.Context, message string, fields map[string]interface{}) {
	t.send(ctx, message, fields, composite_logger.ErrorLevel)
}

func (t TelegramLogger) Fatal(message string, fields map[string]interface{}) {
	t.send(context.TODO(), message, fields, composite_logger.FatalLevel)
}

func (t TelegramLogger) FatalContext(ctx context.Context, message string, fields map[string]interface{}) {
	t.send(ctx, message, fields, composite_logger.FatalLevel)
}

func (t TelegramLogger) send(ctx context.Context, message string, fields map[string]interface{}, level composite_logger.Level) {
	if t.Level > level {
		return
	}

	if ctx != nil && ctx.Err() != nil {
		return
	}

	text := formatTelegramMarkdown(message, fields, level, t)

	tgMessage := tgbotapi.NewMessage(t.LogChatId, text)
	tgMessage.ParseMode = "MarkdownV2"

	if _, err := t.BotApi.Send(tgMessage); err != nil {
		fmt.Printf("[TelegramLogger Error] Failed to send detailed log to ChatID %d: %v\n", t.LogChatId, err)

		// Fallback: send simple plain text message without Markdown
		fallbackText := fmt.Sprintf("⚠️ [TelegramLogger Error]\nFailed to send detailed log.\nError: %v\nMessage: %s", err, message)
		fallbackMsg := tgbotapi.NewMessage(t.LogChatId, fallbackText)
		if _, fallbackErr := t.BotApi.Send(fallbackMsg); fallbackErr != nil {
			fmt.Printf("[TelegramLogger Error] Failed to send fallback message to ChatID %d: %v\n", t.LogChatId, fallbackErr)
		}
	}
}

func formatTelegramMarkdown(message string, fields map[string]interface{}, level composite_logger.Level, t TelegramLogger) string {
	escapeMarkdownV2 := func(text string) string {
		var markdownV2Regex = regexp.MustCompile(`([\[\]\-_*~` + "`" + `>#+=|{}.!])`)
		return markdownV2Regex.ReplaceAllString(text, "\\$1")
	}

	now := time.Now().Format("[2006-01-02 15:04:05]")
	jsonFields, _ := json.MarshalIndent(normalizeLogContext(fields), "", "    ")

	title, ok := t.LevelTitles[level]
	if !ok || title == "" {
		title = strings.ToUpper(level.String())
	}

	var decoration string
	if t.UseLevelTitleWrapper {
		wrapper := t.LevelWrappers[level]
		decoration = fmt.Sprintf("%s *%s* %s\n", wrapper, title, wrapper)
	} else {
		decoration = fmt.Sprintf("*%s*\n", title)
	}

	text := fmt.Sprintf("%s%s %s\n\n```json\n%s\n```",
		decoration,
		escapeMarkdownV2(now),
		escapeMarkdownV2(message),
		string(jsonFields))

	return text
}

func normalizeLogContext(fields map[string]interface{}) map[string]interface{} {
	if fields == nil {
		return nil
	}

	normalized := make(map[string]interface{}, len(fields))
	for key, value := range fields {
		switch typed := value.(type) {
		case error:
			normalized[key] = typed.Error()
		case map[string]interface{}:
			normalized[key] = normalizeLogContext(typed)
		default:
			normalized[key] = value
		}
	}

	return normalized
}
