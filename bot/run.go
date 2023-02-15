package bot

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

const promoPrefix = "Промо"

// Run listens updates
func (b *bot) Run() {
	updatesCfg := tgbotapi.UpdateConfig{
		Offset:  0,
		Timeout: 10,
	}
	for upd := range b.GetUpdatesChan(updatesCfg) {
		if upd.MyChatMember != nil {
			// if user left or kicked bot
			if upd.MyChatMember.NewChatMember.Status == "left" || upd.MyChatMember.NewChatMember.Status == "kicked" {
				// remove chat from active chats
				delete(b.members, upd.MyChatMember.Chat.ID)
			}
		}

		if upd.Message != nil {
			// update chat status on any activity
			b.members[upd.Message.Chat.ID] = struct{}{}

			if upd.Message.IsCommand() {
				key := upd.Message.Command()
				if cmd, ok := b.commands[commandKey(key)]; ok {
					go cmd.action(upd)
				} else {
					b.logger.Error("command handler not found", zap.String("cmd", key))
				}
				continue
			}

			if cmd, ok := b.replyToCommand(upd.Message.Text); ok {
				go cmd.action(upd)
				continue
			}

			if strings.HasPrefix(upd.Message.Text, promoPrefix) {
				go b.PromoCmd(upd)
			} else {
				go b.SearchCmd(upd)
			}
		}

		if upd.CallbackQuery != nil {
			data := upd.CallbackData()
			entity := unmarshallCb(data)

			callback := tgbotapi.NewCallback(upd.CallbackQuery.ID, "")
			b.apiRequest(callback)

			b.callbacks[entity.cbType](upd, entity)
		}
	}
}

func backButton(parentType callbackType, parentIds []string) []tgbotapi.InlineKeyboardButton {
	data := callbackEntity{
		parentType: parentType,
		parentIds:  parentIds[0 : len(parentIds)-1],
		cbType:     Back,
		id:         parentIds[len(parentIds)-1],
	}
	button := tgbotapi.NewInlineKeyboardButtonData("Назад", marshallCb(data))
	return tgbotapi.NewInlineKeyboardRow(button)
}
