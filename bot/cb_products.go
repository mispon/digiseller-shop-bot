package bot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mispon/xbox-store-bot/bot/desc"
	"github.com/mispon/xbox-store-bot/bot/digi"
	"github.com/mispon/xbox-store-bot/bot/search"
	"go.uber.org/zap"
)

const productsPerPage = 10

func productsButtons(products []desc.Product, entity callbackEntity) [][]tgbotapi.InlineKeyboardButton {
	rows := make([][]tgbotapi.InlineKeyboardButton, 0, productsPerPage)
	for _, product := range products {
		data := entity
		data.id = product.Id

		button := tgbotapi.NewInlineKeyboardButtonData(product.Name, marshallCb(data))
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(button))
	}
	return rows
}

func (b *bot) ProductsCallback(upd tgbotapi.Update, subCategoryEntity callbackEntity) {
	subCategoryName, products, hasMore, ok := b.cache.Products(subCategoryEntity.id, subCategoryEntity.page, productsPerPage)
	if !ok {
		b.logger.Error("sub category is empty", zap.String("sub category", subCategoryName))
		return
	}

	buttons := productsButtons(products, callbackEntity{
		parentType: Products,
		parentIds:  append(subCategoryEntity.parentIds, subCategoryEntity.id),
		cbType:     Product,
	})

	rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(buttons)+3)
	rows = append(rows, buttons...)

	pagesRow := tgbotapi.NewInlineKeyboardRow()
	if subCategoryEntity.page > 0 {
		data := subCategoryEntity.Clone()
		data.page--
		button := tgbotapi.NewInlineKeyboardButtonData("<-", marshallCb(data))
		pagesRow = append(pagesRow, button)
	}
	if hasMore {
		data := subCategoryEntity.Clone()
		data.page++
		button := tgbotapi.NewInlineKeyboardButtonData("->", marshallCb(data))
		pagesRow = append(pagesRow, button)
	}
	if len(pagesRow) > 0 {
		rows = append(rows, pagesRow)
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
	rows = append(rows, backButton(productsEntity.parentType, productsEntity.parentIds)...)

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

func (b *bot) SearchProductCallback(upd tgbotapi.Update, productsEntity callbackEntity) {
	const minArsPrice = 600

	product, err := search.GetProduct(b.client, b.opts.search.url, productsEntity.id)
	if err != nil {
		b.logger.Error("product not found", zap.String("product_id", productsEntity.id))
		return
	}

	rows := make([][]tgbotapi.InlineKeyboardButton, 0)
	if price, ok := product.Prices["ARS"]; ok && price != 0 {
		rubPrice := int(price * b.getConversionRate("ARS"))
		if rubPrice < minArsPrice {
			rubPrice = minArsPrice
		}
		buttonARText := fmt.Sprintf(`Купить "Покупкой на аккаунт" 🇦🇷 за %d руб.`, rubPrice)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL(
				buttonARText,
				digi.CustomProductPaymentURL(b.client,
					b.opts.sellerId, product.Name,
					b.opts.search.universalProductId, b.opts.search.universalProductOptionId, rubPrice))))
	}
	if price, ok := product.Prices["TRY"]; ok && price != 0 {
		rubPrice := int(price * b.getConversionRate("TRY"))

		var buttonTRText string
		if product.IsBackwardCompatibil() {
			buttonTRText = fmt.Sprintf(`Купить "Покупкой на аккаунт" 🇹🇷 за %d руб.`, rubPrice)
		} else {
			buttonTRText = fmt.Sprintf(`Купить "Ключ активации" 🇹🇷 за %d руб.`, rubPrice)
		}

		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL(
				buttonTRText,
				digi.CustomProductPaymentURL(b.client,
					b.opts.sellerId, product.Name,
					b.opts.search.universalProductId, b.opts.search.universalProductOptionId, rubPrice))))
	}
	rows = append(rows, backButton(Search, []string{}))

	reply := tgbotapi.NewEditMessageTextAndMarkup(
		upd.CallbackQuery.Message.Chat.ID,
		upd.CallbackQuery.Message.MessageID,
		product.String(),
		tgbotapi.NewInlineKeyboardMarkup(rows...),
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
