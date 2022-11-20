package bot

import (
	"fmt"
	"go.uber.org/zap"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type replyKeyboardValue string

const (
	ReplyCategories = replyKeyboardValue("Категории")
	ReplyReviews    = replyKeyboardValue("Отзывы")
	ReplyHelp       = replyKeyboardValue("Помощь")
)

const (
	inviteUrl = "https://t.me/+1-lMGCQ7zOphODUy"
)

func (b *bot) StartCmd(upd tgbotapi.Update) {
	name := upd.Message.From.UserName
	if name == "" {
		name = upd.Message.From.FirstName
	}
	message := `
Добро пожаловать в <b>Xbox Store | Бот-магазин подписок и игр</b>, %s!

Не забывайте подписаться на <a href='%s'>наш основной канал</a> и следить за новостями и выходе игровых новинок и об акциях.
`
	reply := tgbotapi.NewMessage(upd.Message.Chat.ID, fmt.Sprintf(message, name, inviteUrl))
	reply.ParseMode = "html"

	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(string(ReplyCategories)),
			tgbotapi.NewKeyboardButton(string(ReplyReviews)),
			tgbotapi.NewKeyboardButton(string(ReplyHelp)),
		),
	)
	reply.ReplyMarkup = keyboard
	reply.DisableWebPagePreview = true

	if err := b.apiRequest(reply); err != nil {
		b.logger.Error("failed to send start message", zap.Error(err))
	}
}
