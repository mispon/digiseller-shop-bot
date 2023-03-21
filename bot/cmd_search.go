package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func (b *bot) getSearchParams(chatID int64) searchParams {
	b.chatsMu.RLock()
	defer b.chatsMu.RUnlock()

	if chat, ok := b.chats[chatID]; ok {
		return chat.searchParams
	} else {
		return searchParams{}
	}
}

func (b *bot) updateSearchParams(chatID int64, params searchParams) {
	b.chatsMu.Lock()
	defer b.chatsMu.Unlock()

	if chat, ok := b.chats[chatID]; ok {
		if params.category != "" {
			chat.searchParams.category = params.category
		}
		if params.query != "" {
			chat.searchParams.query = params.query
		}
		if params.messageID != 0 {
			chat.searchParams.messageID = params.messageID
		}

		b.chats[chatID] = chat
	}
}

func (b *bot) clearSearchParams(chatID int64) {
	b.chatsMu.RLock()

	if chatDesc, ok := b.chats[chatID]; ok {
		b.chatsMu.RUnlock()
		if chatDesc.searchParams.query != "" {
			b.chatsMu.Lock()
			defer b.chatsMu.Unlock()
			b.chats[chatID] = chat{searchParams{}}
		}
	} else {
		b.chatsMu.RUnlock()
	}
}

func (b *bot) SearchCmd(upd tgbotapi.Update) {
	if len(upd.Message.Text) < 3 {
		// ignore too short messages
		return
	}

	reply := tgbotapi.NewDeleteMessage(upd.Message.Chat.ID, upd.Message.MessageID)
	err := b.apiRequest(reply)
	if err != nil {
		b.logger.Error("failed to delete message", zap.Error(err))
	}

	b.updateSearchParams(upd.Message.Chat.ID, searchParams{query: upd.Message.Text})

	b.SearchCallback(upd, callbackEntity{cbType: Search})
}
