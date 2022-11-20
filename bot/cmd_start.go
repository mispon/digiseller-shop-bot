package bot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type replyKeyboardValue string

const (
	ReplyCategories = replyKeyboardValue("Категории")
	ReplyReviews    = replyKeyboardValue("Отзывы")
	ReplyHelp       = replyKeyboardValue("Помощь")
)

func (b *bot) StartCmd(upd tgbotapi.Update) {
	name := upd.Message.From.UserName
	if name == "" {
		name = upd.Message.From.FirstName
	}
	message := fmt.Sprintf("Добро пожаловать в xbox-stoge, %s!", name)
	reply := tgbotapi.NewMessage(upd.Message.Chat.ID, message)

	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(string(ReplyCategories)),
			tgbotapi.NewKeyboardButton(string(ReplyReviews)),
			tgbotapi.NewKeyboardButton(string(ReplyHelp)),
		),
	)
	reply.ReplyMarkup = keyboard

	b.apiRequest(reply)
}
