package bot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *bot) StartCmd(upd tgbotapi.Update) {
	message := fmt.Sprintf("Welcome, %s!", upd.Message.From.UserName)
	reply := tgbotapi.NewMessage(upd.Message.Chat.ID, message)
	b.apiRequest(reply)
}
