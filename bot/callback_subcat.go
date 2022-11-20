package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func (b *bot) SubCategoryCallback(upd tgbotapi.Update, category callbackEntity) {
	categoryName, subs, ok := b.cache.SubCategory(category.id)
	if !ok {
		b.logger.Error("sub categories not found", zap.String("category", categoryName))
	}

	rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(subs)+1)
	for _, sc := range subs {
		data := callbackEntity{
			parentType: SubCategory,
			parentIds:  append(category.parentIds, category.id),
			cbType:     Products,
			id:         sc.Id,
		}
		button := tgbotapi.NewInlineKeyboardButtonData(sc.Name, marshallCb(data))
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(button))
	}
	rows = append(rows, backButton(Categories, category.parentIds))

	reply := tgbotapi.NewEditMessageTextAndMarkup(
		upd.CallbackQuery.Message.Chat.ID,
		upd.CallbackQuery.Message.MessageID,
		categoryName,
		tgbotapi.NewInlineKeyboardMarkup(rows...),
	)

	if err := b.apiRequest(reply); err != nil {
		b.logger.Error("failed to show sub categories", zap.Error(err))
	}
}
