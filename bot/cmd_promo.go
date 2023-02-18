package bot

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func (b *bot) PromoCmd(upd tgbotapi.Update) {
	if _, ok := admins[upd.Message.From.UserName]; !ok {
		return
	}

	message := strings.TrimPrefix(upd.Message.Text, promoPrefix)

	b.chatsMu.RLock()
	defer b.chatsMu.RUnlock()

	for chatID := range b.chats {
		reply := tgbotapi.NewMessage(chatID, message)
		reply.ParseMode = "html"
		reply.DisableWebPagePreview = false

		if err := b.apiRequest(reply); err != nil {
			b.logger.Error("failed to send promo message", zap.Error(err))
		}
	}
}
