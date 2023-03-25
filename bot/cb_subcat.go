package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mispon/digiseller-shop-bot/bot/desc"
	"github.com/mispon/digiseller-shop-bot/bot/search"
	"go.uber.org/zap"
)

func subCategoriesButtons(subs []desc.SubCategory, entity callbackEntity) [][]tgbotapi.InlineKeyboardButton {
	rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(subs)+1)
	for _, sc := range subs {
		data := entity
		data.id = sc.Id

		button := tgbotapi.NewInlineKeyboardButtonData(sc.Name, marshallCb(data))
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(button))
	}
	return rows
}

func (b *bot) SubCategoryCallback(upd tgbotapi.Update, category callbackEntity) {
	categoryName, subs, ok := b.cache.SubCategory(category.id)
	if !ok {
		b.logger.Error("sub categories not found", zap.String("category", categoryName))
		return
	}

	rows := subCategoriesButtons(subs, callbackEntity{
		parentType: SubCategory,
		parentIds:  append(category.parentIds, category.id),
		cbType:     Products,
	})
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

func (b *bot) SearchSubCategoryCallback(upd tgbotapi.Update, subCategoryEntity callbackEntity) {
	subs, err := search.Categories(b.client, b.opts.search.url)
	if err != nil {
		b.logger.Error("failed to get search categories", zap.Error(err))
		return
	}

	rows := subCategoriesButtons(subs, callbackEntity{
		parentType: SearchSubCategory,
		parentIds:  append(subCategoryEntity.parentIds, subCategoryEntity.id),
		cbType:     SearchParams,
	})

	data := callbackEntity{
		parentType: SearchSubCategory,
		parentIds:  append(subCategoryEntity.parentIds, subCategoryEntity.id),
		cbType:     SearchInstruction,
		id:         subCategoryEntity.id,
	}
	instructionButton := tgbotapi.NewInlineKeyboardButtonData("Инструкция", marshallCb(data))
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(instructionButton))
	rows = append(rows, backButton(Categories, subCategoryEntity.parentIds))

	reply := tgbotapi.NewEditMessageTextAndMarkup(
		upd.CallbackQuery.Message.Chat.ID,
		upd.CallbackQuery.Message.MessageID,
		search.CategoriesDescriptionMessage,
		tgbotapi.NewInlineKeyboardMarkup(rows...),
	)
	reply.ParseMode = "html"
	reply.DisableWebPagePreview = true

	if err := b.apiRequest(reply); err != nil {
		b.logger.Error("failed to show sub categories", zap.Error(err))
	}
}
