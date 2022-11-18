package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

// Run listens updates
func (b *bot) Run() {
	updatesCfg := tgbotapi.UpdateConfig{
		Offset:  0,
		Timeout: 10,
	}
	for upd := range b.GetUpdatesChan(updatesCfg) {
		if upd.Message == nil {
			continue
		}
		b.logger.Debug("receive msg", zap.String("msg", upd.Message.Text))

		if upd.Message.IsCommand() {
			key := upd.Message.Command()
			if cmd, ok := b.commands[key]; ok {
				go cmd.action(upd)
			} else {
				b.logger.Error("command handler not found", zap.String("cmd", key))
			}
		}
	}
}
