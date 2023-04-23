package bot

import (
	"encoding/json"
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

	if err := json.Unmarshal([]byte(jsonConfig), &config); err == nil {
		b.userConfig.Lock()
		defer b.userConfig.Unlock()

		if len(config.ConversionRates) != 0 {
			b.userConfig.ConversionRates = config.ConversionRates
		}
		if config.MinARSPrice != 0 {
			b.userConfig.MinARSPrice = config.MinARSPrice
		}
	} else {
		result = "Неверный формат"
	}

	reply := tgbotapi.NewMessage(upd.Message.Chat.ID, result)
	if err := b.apiRequest(reply); err != nil {
		b.logger.Error("failed to send online message", zap.Error(err))
	}
}
