package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func (b *bot) ShopCmd(upd tgbotapi.Update) {
	reply := tgbotapi.NewMessage(upd.Message.Chat.ID, "Категории")
	reply.ReplyMarkup = b.getCategoriesKeyboard()

	if err := b.apiRequest(reply); err != nil {
		b.logger.Error("failed to show categories", zap.Error(err))
	}
}

func (b *bot) ShopCallback(upd tgbotapi.Update) {
	reply := tgbotapi.NewEditMessageTextAndMarkup(
		upd.CallbackQuery.Message.Chat.ID,
		upd.CallbackQuery.Message.MessageID,
		"Категории",
		b.getCategoriesKeyboard(),
	)

	if err := b.apiRequest(reply); err != nil {
		b.logger.Error("failed to edit categories", zap.Error(err))
	}
}

func (b *bot) getCategoriesKeyboard() tgbotapi.InlineKeyboardMarkup {
	categories := b.cache.Categories()

	rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(categories)+1)
	for _, category := range categories {
		data := callbackEntity{
			parentType: Categories,
			cbType:     SubCategory,
			id:         category.Id,
		}

		button := tgbotapi.NewInlineKeyboardButtonData(category.Name, marshallCb(data))
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(button))
	}

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
