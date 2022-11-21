package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *bot) BackCallback(upd tgbotapi.Update, entity callbackEntity) {
	switch entity.parentType {
	case Categories:
		b.ShopCallback(upd)
	case SubCategory:
		b.SubCategoryCallback(upd, entity)
	case Products:
		b.ProductsCallback(upd, entity)
	case Product:
		b.ProductCallback(upd, entity)
	}
}
