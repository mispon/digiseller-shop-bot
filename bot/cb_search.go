package bot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mispon/xbox-store-bot/bot/search"
	"go.uber.org/zap"
)

func (b *bot) SearchCallback(upd tgbotapi.Update, searchEntity callbackEntity) {
	var chatID int64

	if upd.CallbackQuery != nil {
		chatID = upd.CallbackQuery.Message.Chat.ID
	} else if upd.Message != nil {
		chatID = upd.Message.Chat.ID
	}

	sp := b.getSearchParams(chatID)
	if sp.query == "" || sp.category == "" {
		return
	}

	products, total, err := search.Search(b.client, b.opts.search.url, sp.category, sp.query, productsPerPage, searchEntity.skip)
	if err != nil {
		b.logger.Warn("failed to search in xbox com", zap.Error(err))

		reply := tgbotapi.NewEditMessageTextAndMarkup(
			chatID,
			sp.messageID,
			search.EmptyResultMessage,
			tgbotapi.NewInlineKeyboardMarkup(backButton(SearchSubCategory, []string{})),
		)
		b.apiRequest(reply)
		return
	}

	buttons := productsButtons(products, callbackEntity{parentType: Search, cbType: SearchProduct})

	rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(buttons)+3)
	rows = append(rows, buttons...)

	pagesRow := tgbotapi.NewInlineKeyboardRow()
	if searchEntity.skip > 0 {
		data := callbackEntity{
			cbType: Search,
		}
		data.skip = searchEntity.skip - productsPerPage
		button := tgbotapi.NewInlineKeyboardButtonData("<-", marshallCb(data))
		pagesRow = append(pagesRow, button)
	}
	if total > searchEntity.skip+productsPerPage {
		data := callbackEntity{
			cbType: Search,
		}
		data.skip = searchEntity.skip + productsPerPage
		button := tgbotapi.NewInlineKeyboardButtonData("->", marshallCb(data))
		pagesRow = append(pagesRow, button)
	}
	if len(pagesRow) > 0 {
		rows = append(rows, pagesRow)
	}

	rows = append(rows, backButton(SearchSubCategory, []string{}))

	reply := tgbotapi.NewEditMessageTextAndMarkup(
		chatID,
		sp.messageID,
		sp.query,
		tgbotapi.NewInlineKeyboardMarkup(rows...),
	)

	if err := b.apiRequest(reply); err != nil {
		b.logger.Error("failed to search products", zap.Error(err))
	}
}

func (b *bot) SearchParamsCallback(upd tgbotapi.Update, searchParamsEntity callbackEntity) {
	backButton := backButton(SearchSubCategory, []string{})

	reply := tgbotapi.NewEditMessageTextAndMarkup(
		upd.CallbackQuery.Message.Chat.ID,
		upd.CallbackQuery.Message.MessageID,
		search.CategoryDescriptionMessage,
		tgbotapi.NewInlineKeyboardMarkup(backButton),
	)

	if err := b.apiRequest(reply); err != nil {
		b.logger.Error("failed to show search description", zap.Error(err))
	}

	b.updateSearchParams(upd.CallbackQuery.Message.Chat.ID,
		searchParams{
			category:  searchParamsEntity.id,
			messageID: upd.CallbackQuery.Message.MessageID,
		})
}

func (b *bot) SearchInstructionCallback(upd tgbotapi.Update, productEntity callbackEntity) {
	rows := backButton(productEntity.parentType, productEntity.parentIds)

	reply := tgbotapi.NewEditMessageTextAndMarkup(
		upd.CallbackQuery.Message.Chat.ID,
		upd.CallbackQuery.Message.MessageID,
		fmt.Sprintf(
			search.InstructionMessage,
			b.getConversionRate("TRY"),
			b.getConversionRate("ARS")),
		tgbotapi.NewInlineKeyboardMarkup(rows),
	)
	reply.ParseMode = "html"
	reply.DisableWebPagePreview = true

	if err := b.apiRequest(reply); err != nil {
		b.logger.Error("failed to show search instruction", zap.Error(err))
	}
}
