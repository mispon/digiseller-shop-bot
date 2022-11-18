package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type (
	bot struct {
		*tgbotapi.BotAPI

		logger   *zap.Logger
		sellerId string

		commands map[string]commandEntity

		opts options
	}
)

// New creates bot instance
func New(logger *zap.Logger, token, sellerId string, opts ...Option) (*bot, error) {
	api, aErr := tgbotapi.NewBotAPI(token)
	if aErr != nil {
		return nil, aErr
	}

	var bo options
	for _, optFn := range opts {
		optFn(&bo)
	}
	api.Debug = bo.debug

	b := &bot{
		BotAPI:   api,
		logger:   logger,
		sellerId: sellerId,
		commands: make(map[string]commandEntity),
		opts:     bo,
	}

	if err := b.initCommands(); err != nil {
		return nil, err
	}

	logger.Info("bot created", zap.String("username", api.Self.UserName))
	return b, nil
}

func (b *bot) apiRequest(c tgbotapi.Chattable) error {
	_, err := b.Request(c)
	return err
}
