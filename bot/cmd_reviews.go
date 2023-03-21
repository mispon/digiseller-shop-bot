package bot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mispon/xbox-store-bot/bot/digi"
	"go.uber.org/zap"
)

const (
	reviewsURL = "https://x-box-store.ru/reviews"
)

func (b *bot) ReviewsCmd(upd tgbotapi.Update) {
	message := fmt.Sprintf("Отзывы покупателей\n<a href='%s'>&#8205;</a>", digi.ReviewslogoUrl)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Смотреть", reviewsURL),
		),
	)

	reply := tgbotapi.NewMessage(upd.Message.Chat.ID, message)
	reply.ReplyMarkup = keyboard
	reply.ParseMode = "html"

	if err := b.apiRequest(reply); err != nil {
		b.logger.Error("failed to send reviews message", zap.Error(err))
	}
}
