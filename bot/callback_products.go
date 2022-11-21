package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func (b *bot) ProductsCallback(upd tgbotapi.Update, subCategoryEntity callbackEntity) {
	subCategoryName, products, ok := b.cache.Products(subCategoryEntity.id)
	if !ok {
		b.logger.Error("sub category is empty", zap.String("sub category", subCategoryName))
		return
	}

	rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(products)+1)
	for _, product := range products {
		data := callbackEntity{
			parentType: Products,
			parentIds:  append(subCategoryEntity.parentIds, subCategoryEntity.id),
			cbType:     Product,
			id:         product.Id,
		}
		button := tgbotapi.NewInlineKeyboardButtonData(product.Name, marshallCb(data))
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(button))
	}
	rows = append(rows, backButton(SubCategory, subCategoryEntity.parentIds))

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

func (b *bot) ProductCallback(upd tgbotapi.Update, productsEntity callbackEntity) {
	subCategoryId := productsEntity.parentIds[len(productsEntity.parentIds)-1]
	product, ok := b.cache.Product(subCategoryId, productsEntity.id)
	if !ok {
		b.logger.Error("product not found", zap.String("sub_category_id", subCategoryId), zap.String("product_id", productsEntity.id))
		return
	}

	rows := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL("Купить", product.PaymentURL(b.opts.sellerId)),
	)
	if product.AddInfo != "" {
		data := callbackEntity{
			parentType: Product,
			parentIds:  append(productsEntity.parentIds, productsEntity.id),
			cbType:     ProductInstruction,
			id:         productsEntity.id,
		}
		button := tgbotapi.NewInlineKeyboardButtonData("Инструкция", marshallCb(data))
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(button)...)
	}
	rows = append(rows, backButton(Products, productsEntity.parentIds)...)

	reply := tgbotapi.NewEditMessageTextAndMarkup(
		upd.CallbackQuery.Message.Chat.ID,
		upd.CallbackQuery.Message.MessageID,
		product.String(),
		tgbotapi.NewInlineKeyboardMarkup(rows),
	)
	reply.ParseMode = "html"
	reply.DisableWebPagePreview = false

	if err := b.apiRequest(reply); err != nil {
		b.logger.Error("failed to show product", zap.String("product", product.Name), zap.Error(err))
	}
}

func (b *bot) ProductInstructionCallback(upd tgbotapi.Update, productEntity callbackEntity) {
	subCategoryId := productEntity.parentIds[len(productEntity.parentIds)-2]
	product, ok := b.cache.Product(subCategoryId, productEntity.id)
	if !ok {
		b.logger.Error("product not found", zap.String("sub_category_id", subCategoryId), zap.String("product_id", productEntity.id))
		return
	}

	rows := backButton(Product, productEntity.parentIds)

	reply := tgbotapi.NewEditMessageTextAndMarkup(
		upd.CallbackQuery.Message.Chat.ID,
		upd.CallbackQuery.Message.MessageID,
		product.Instruction(),
		tgbotapi.NewInlineKeyboardMarkup(rows),
	)
	reply.ParseMode = "html"
	reply.DisableWebPagePreview = true

	if err := b.apiRequest(reply); err != nil {
		b.logger.Error("failed to show product instruction", zap.String("product", product.Name), zap.Error(err))
	}
}
