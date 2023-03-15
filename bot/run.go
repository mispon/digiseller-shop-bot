package bot

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

const (
	onlinePrefix          = "Онлайн"
	promoPrefix           = "Промо"
	conversionRatesPrefix = "Курс"
)

// Run listens updates
func (b *bot) Run() {
	updatesCfg := tgbotapi.UpdateConfig{
		Offset:  0,
		Timeout: 10,
	}
	for upd := range b.GetUpdatesChan(updatesCfg) {
		go b.processUpdate(upd)
	}
}

func (b *bot) processUpdate(upd tgbotapi.Update) {
	if upd.MyChatMember != nil {
		// if user left or kicked bot
		if upd.MyChatMember.NewChatMember.Status == "left" || upd.MyChatMember.NewChatMember.Status == "kicked" {
			b.chatsMu.Lock()
			delete(b.chats, upd.MyChatMember.Chat.ID)
			b.chatsMu.Unlock()
		}
	}

	if upd.Message != nil {
		if !b.chatExist(upd.Message.Chat.ID) {
			b.addChat(upd.Message.Chat.ID)
		}

		if upd.Message.IsCommand() {
			key := upd.Message.Command()
			if cmd, ok := b.commands[commandKey(key)]; ok {
				go cmd.action(upd)
			} else {
				b.logger.Error("command handler not found", zap.String("cmd", key))
			}
			return
		}

		if cmd, ok := b.replyToCommand(upd.Message.Text); ok {
			go cmd.action(upd)
			return
		}

		if strings.HasPrefix(upd.Message.Text, onlinePrefix) {
			go b.OnlineCmd(upd)
		} else if strings.HasPrefix(upd.Message.Text, promoPrefix) {
			go b.PromoCmd(upd)
		} else if strings.HasPrefix(upd.Message.Text, conversionRatesPrefix) {
			b.ConversionRates(upd)
		} else {
			go b.SearchCmd(upd)
		}
	}

	if upd.CallbackQuery != nil {
		data := upd.CallbackData()
		entity := unmarshallCb(data)

		if entity.cbType != Search && entity.parentType != Search {
			b.clearSearchParams(upd.CallbackQuery.Message.Chat.ID)
		}

		callback := tgbotapi.NewCallback(upd.CallbackQuery.ID, "")
		b.apiRequest(callback)

		b.callbacks[entity.cbType](upd, entity)
	}
}

func backButton(parentType callbackType, parentIds []string) []tgbotapi.InlineKeyboardButton {
	data := callbackEntity{
		parentType: parentType,
		cbType:     Back,
	}
	if len(parentIds) > 0 {
		data.id = parentIds[len(parentIds)-1]
		data.parentIds = parentIds[0 : len(parentIds)-1]
	}
	button := tgbotapi.NewInlineKeyboardButtonData("Назад", marshallCb(data))
	return tgbotapi.NewInlineKeyboardRow(button)
}

func (b *bot) Stop() {
	b.writeChats()
}
