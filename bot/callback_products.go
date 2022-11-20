package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func (b *bot) ProductsCallback(upd tgbotapi.Update, subCategory callbackEntity) {
	subCategoryName, products, ok := b.cache.Products(subCategory.id)
	if !ok {
		b.logger.Error("products not found", zap.String("sub category", subCategoryName))
	}

	rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(products)+1)
	for _, product := range products {
		data := callbackEntity{
			parentType: Products,
			parentIds:  append(subCategory.parentIds, subCategory.id),
			cbType:     Product,
			id:         product.Id,
		}
		button := tgbotapi.NewInlineKeyboardButtonData(product.Name, marshallCb(data))
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(button))
	}
	rows = append(rows, backButton(SubCategory, subCategory.parentIds))

	reply := tgbotapi.NewEditMessageTextAndMarkup(
		upd.CallbackQuery.Message.Chat.ID,
		upd.CallbackQuery.Message.MessageID,
		subCategoryName,
		tgbotapi.NewInlineKeyboardMarkup(rows...),
	)

	if err := b.apiRequest(reply); err != nil {
		b.logger.Error("failed to show products", zap.Error(err))
	}
}

func (b *bot) ProductCallback(upd tgbotapi.Update, producs callbackEntity) {
	parentId := producs.parentIds[len(producs.parentIds)-1]
	product, ok := b.cache.Product(parentId, producs.id)
	if !ok {
		b.logger.Error("product not found", zap.String("product_id", producs.id))
	}

	rows := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL("Купить", product.PaymentURL(b.opts.sellerId)),
	)
	rows = append(rows, backButton(Products, producs.parentIds)...)

	reply := tgbotapi.NewEditMessageTextAndMarkup(
		upd.CallbackQuery.Message.Chat.ID,
		upd.CallbackQuery.Message.MessageID,
		product.String(),
		tgbotapi.NewInlineKeyboardMarkup(rows),
	)
	reply.ParseMode = "html"
	reply.DisableWebPagePreview = false

	if err := b.apiRequest(reply); err != nil {
		b.logger.Error("failed to show product", zap.Error(err))
	}
}
