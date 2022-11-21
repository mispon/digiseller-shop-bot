package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"strings"
	"time"
)

var (
	blacklist = map[string]struct{}{
		"xbox":     {},
		"ps4":      {},
		"fortnite": {},
		"ключ":     {},
		"аккаунт":  {},
		"подписка": {},
	}
)

func (b *bot) SearchCmd(upd tgbotapi.Update) {
	if len(upd.Message.Text) < 3 {
		// ignore too short messages
		return
	}

	text := strings.ToLower(upd.Message.Text)
	if _, ok := blacklist[text]; ok {
		// ingore the must common words
		return
	}

	products, ok := b.cache.Search(text)
	if !ok {
		reply := tgbotapi.NewMessage(upd.Message.Chat.ID, "🤖")
		b.apiRequest(reply)
		return
	}

	for _, product := range products {
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL("Купить", product.PaymentURL(b.opts.sellerId)),
			),
		)

		reply := tgbotapi.NewMessage(upd.Message.Chat.ID, product.String())
		reply.ReplyMarkup = keyboard
		reply.ParseMode = "html"
		reply.DisableWebPagePreview = false

		if err := b.apiRequest(reply); err != nil {
			b.logger.Error("failed to show product", zap.String("product", product.Name), zap.Error(err))
		}

		time.Sleep(time.Second)
	}
}
