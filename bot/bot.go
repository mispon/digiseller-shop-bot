package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mispon/xbox-store-bot/bot/desc"
	"go.uber.org/zap"
)

type (
	bot struct {
		*tgbotapi.BotAPI

		logger *zap.Logger
		cache  inMemoryCache
		opts   options

		commands  map[commandKey]commandEntity
		callbacks map[callbackType]callbackFn
	}

	inMemoryCache interface {
		Categories() []desc.Category
		SubCategory(categoryId string) (string, []desc.SubCategory, bool)
		Products(subCategoryId string) (string, []desc.Product, bool)
		Product(subCategoryId, productId string) (desc.Product, bool)
		Search(text string) ([]desc.Product, bool)
	}
)

// New creates bot instance
func New(logger *zap.Logger, cache inMemoryCache, token string, opts ...Option) (*bot, error) {
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
		BotAPI:    api,
		logger:    logger.Named("bot"),
		cache:     cache,
		opts:      bo,
		commands:  make(map[commandKey]commandEntity),
		callbacks: make(map[callbackType]callbackFn),
	}

	if err := b.initCommands(); err != nil {
		return nil, err
	}
	b.initCallbacks()

	b.logger.Info("bot created", zap.String("username", api.Self.UserName))
	return b, nil
}

func (b *bot) apiRequest(c tgbotapi.Chattable) error {
	_, err := b.Request(c)
	return err
}
