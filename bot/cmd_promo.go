package bot

import (
	"go.uber.org/zap"
	"strconv"
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

	for _, key := range b.getAllKeysFromRedis() {
		chatID, err := strconv.ParseInt(key, 10, 64)
		if err != nil {
			b.logger.Error("failed to convert chatID key to int64", zap.Error(err))
			continue
		}

		reply := tgbotapi.NewMessage(chatID, message)
		reply.ParseMode = "html"
		reply.DisableWebPagePreview = false

		if err = b.apiRequest(reply); err != nil {
			b.logger.Error("failed to send promo message", zap.Error(err))
		}
	}
}

func (b *bot) getAllKeysFromRedis() []string {
	var keys []string
	iter := b.rdb.Scan(b.ctx, 0, "*", 0).Iterator()
	for iter.Next(b.ctx) {
		keys = append(keys, iter.Val())
	}
	return keys
}
