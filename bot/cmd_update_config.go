package bot

import (
	"encoding/json"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func (b *bot) UpdateUserConfigCmd(upd tgbotapi.Update) {
	if _, ok := admins[upd.Message.From.UserName]; !ok {
		return
	}

	jsonConfig := strings.TrimPrefix(upd.Message.Text, userConfigPrefix)
	jsonConfig = strings.Trim(jsonConfig, "\n")

	var (
		config = userConfig{}
		result = "Ok"
	)

	if jsonConfig == "" {
		conf, _ := json.Marshal(&b.userConfig)
		result = fmt.Sprintf("Текущая конфигурация \n%s", conf)
		goto SendMessage
	}

	if err := json.Unmarshal([]byte(jsonConfig), &config); err == nil {
		b.userConfig.Lock()
		defer b.userConfig.Unlock()

		if _, ok := ProductsDisplayTypes[config.ProductsDisplayType]; !ok {
			result = fmt.Sprintf("Неверный формат типа отображения: %s", config.ProductsDisplayType)
			goto SendMessage
		}
		b.userConfig.ProductsDisplayType = config.ProductsDisplayType

		if len(config.ConversionRates) != 0 {
			b.userConfig.ConversionRates = config.ConversionRates
		}
		if len(config.BotProducts) != 0 {
			for _, botProduct := range config.BotProducts {
				if _, ok := purchaseTypes[botProduct.PurchaseType]; !ok {
					result = fmt.Sprintf("Неверный формат: %s", botProduct.PurchaseType)
					goto SendMessage
				}
			}
			b.userConfig.BotProducts = config.BotProducts
		}
	} else {
		result = "Неверный формат"
	}
SendMessage:
	reply := tgbotapi.NewMessage(upd.Message.Chat.ID, result)
	if err := b.apiRequest(reply); err != nil {
		b.logger.Error("failed to send online message", zap.Error(err))
	}
}
