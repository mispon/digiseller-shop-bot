package bot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func (b *bot) OnlineCmd(upd tgbotapi.Update) {
	if _, ok := admins[upd.Message.From.UserName]; !ok {
		return
	}

	reply := tgbotapi.NewMessage(upd.Message.Chat.ID, fmt.Sprintf("Активных чатов: %d", len(b.chats)))

	if err := b.apiRequest(reply); err != nil {
		b.logger.Error("failed to send online message", zap.Error(err))
	}
}
