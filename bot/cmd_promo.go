package bot

import (
	"go.uber.org/zap"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var admins = map[string]struct{}{
	"noobmaster111": {},
	"Mispon":        {},
}

func (b *bot) PromoCmd(upd tgbotapi.Update) {
	if _, ok := admins[upd.Message.From.UserName]; !ok {
		return
	}

	message := strings.TrimPrefix(upd.Message.Text, promoPrefix)

	for chatID := range b.members {
		reply := tgbotapi.NewMessage(chatID, message)
		reply.ParseMode = "html"
		reply.DisableWebPagePreview = false

		if err := b.apiRequest(reply); err != nil {
			b.logger.Error("failed to send promo message", zap.Error(err))
		}
	}
}
