package bot

import (
	"encoding/json"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func (b *bot) ConversionRates(upd tgbotapi.Update) {
	if _, ok := admins[upd.Message.From.UserName]; !ok {
		return
	}

	jsonMap := strings.TrimPrefix(upd.Message.Text, conversionRatesPrefix)
	var conversionRatesMap map[string]float64

	err := json.Unmarshal([]byte(jsonMap), &conversionRatesMap)
	if err != nil || len(conversionRatesMap) == 0 {
		reply := tgbotapi.NewMessage(upd.Message.Chat.ID, "Неверный формат \nПример: Курс {\"ARS\":0.75, \"TRY\":6}")

		if err := b.apiRequest(reply); err != nil {
			b.logger.Error("failed to send online message", zap.Error(err))
		}
		return
	}

	b.convRatesMu.Lock()
	defer b.convRatesMu.Unlock()
	b.conversionRates = conversionRatesMap
}
