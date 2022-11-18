package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mispon/xbox-store-bot/utils/http"
	"go.uber.org/zap"
)

type (
	categoryResp struct {
		Items []categoryItem `json:"category"`
	}

	categoryItem struct {
		Id   string
		Name string
		Sub  []categorySubItem
	}

	categorySubItem struct {
		Id   string
		Name string
	}
)

func (b *bot) ShopCmd(upd tgbotapi.Update) {
	resp, err := http.Get[categoryResp](b.Client, categoryUrl+b.sellerId)
	if err != nil {
		b.logger.Error("failed to get categories", zap.Error(err))
		return
	}

	var list string
	for _, category := range resp.Items {
		list += category.Name + "\n"
		for _, sub := range category.Sub {
			list += "\t" + sub.Name + "\n"
		}
		list += "\n"
	}

	reply := tgbotapi.NewMessage(upd.Message.Chat.ID, list)
	b.apiRequest(reply)
}
