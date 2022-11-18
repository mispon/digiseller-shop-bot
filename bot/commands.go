package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type (
	commandEntity struct {
		key    string
		desc   string
		action func(upd tgbotapi.Update)
	}
)

func (b *bot) initCommands() error {
	commands := []commandEntity{
		{
			key:    "start",
			desc:   "Запустить бота",
			action: b.StartCmd,
		},
		{
			key:    "shop",
			desc:   "Показать товары",
			action: b.ShopCmd,
		},
	}

	tgCommands := make([]tgbotapi.BotCommand, 0, len(commands))
	for _, cmd := range commands {
		b.commands[cmd.key] = cmd
		tgCommands = append(tgCommands, tgbotapi.BotCommand{
			Command:     "/" + cmd.key,
			Description: cmd.desc,
		})
	}

	config := tgbotapi.NewSetMyCommands(tgCommands...)
	return b.apiRequest(config)
}
