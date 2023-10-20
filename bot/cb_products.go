package bot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mispon/digiseller-shop-bot/bot/countries"
	"github.com/mispon/digiseller-shop-bot/bot/desc"
	"github.com/mispon/digiseller-shop-bot/bot/digi"
	"github.com/mispon/digiseller-shop-bot/bot/search"
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
	product, err := search.GetProduct(b.client, b.opts.search.url, productsEntity.id)
	if err != nil {
		b.logger.Error("product not found", zap.String("product_id", productsEntity.id))
		return
	}
	var (
		userConfig  = b.getUserConfig()
		displayType = userConfig.ProductsDisplayType

		rows               = make([][]tgbotapi.InlineKeyboardButton, 0)
		curentPrice        = 0
		currentKeyboardRow = []tgbotapi.InlineKeyboardButton{}
	)

	for idx, botProduct := range userConfig.BotProducts {
		if botProduct.SkipBackwardCompatibil && product.IsBackwardCompatibil() {
			continue
		}
		country, err := countries.GetCountry(botProduct.Country)
		if err != nil {
			continue
		}

		if price, ok := product.Prices[country.Currency]; ok && price != 0 {
			rubPrice := int(price * userConfig.ConversionRates[country.Currency])
			if rubPrice < botProduct.MinPrice {
				rubPrice = botProduct.MinPrice
			}
			if curentPrice == 0 ||
				displayType == ProductsDisplayTypeAll ||
				(displayType == ProductsDisplayTypeMaxPrice && curentPrice < rubPrice) ||
				(displayType == ProductsDisplayTypeMinPrice && curentPrice > rubPrice) {
				curentPrice = rubPrice

				var buttonText string
				if botProduct.PurchaseType == purchaseTypeAcc {
					buttonText = fmt.Sprintf(`"Покупкой на аккаунт" %s за %d руб.`, country.Flag, rubPrice)
				} else if botProduct.PurchaseType == purchaseTypeKey {
					buttonText = fmt.Sprintf(`"Ключ активации" %s за %d руб.`, country.Flag, rubPrice)
				} else {
					continue
				}

				currentKeyboardRow = tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonURL(
						buttonText,
						digi.CustomProductPaymentURL(b.client,
							b.opts.sellerId, product.Name,
							b.opts.search.universalProductId, b.opts.search.universalProductOptionId, rubPrice)))
			}

			if displayType == ProductsDisplayTypeAll ||
				((displayType == ProductsDisplayTypeMaxPrice || displayType == ProductsDisplayTypeMinPrice) &&
					idx+1 == len(userConfig.BotProducts)) {
				rows = append(rows, currentKeyboardRow)
			}
		}
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
