package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func (b *bot) HelpCmd(upd tgbotapi.Update) {
	message := `
💬 <b>Поддержка:</b> @noobmaster111, пишите если у вас возникли какие то проблемы.
⌚ <b>Онлайн:</b> Примерно с 10:00 - 00:00 по мск.
❗Ничего не покупаю и не беру на реализацию, рекламы в боте нет.
	`
	reply := tgbotapi.NewMessage(upd.Message.Chat.ID, message)
	reply.ParseMode = "html"

	if err := b.apiRequest(reply); err != nil {
		b.logger.Error("failed to send help message", zap.Error(err))
	}
}
