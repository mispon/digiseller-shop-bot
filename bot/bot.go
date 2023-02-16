package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mispon/xbox-store-bot/bot/desc"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type (
	bot struct {
		*tgbotapi.BotAPI

		ctx    context.Context
		logger *zap.Logger
		cache  inMemoryCache
		rdb    *redis.Client
		opts   options

		commands  map[commandKey]commandEntity
		callbacks map[callbackType]callbackFn
	}

	inMemoryCache interface {
		Categories() []desc.Category
		SubCategory(categoryId string) (string, []desc.SubCategory, bool)
		Products(subCategoryId string, page, total int) (string, []desc.Product, bool, bool)
		Product(subCategoryId, productId string) (desc.Product, bool)
		Search(text string) ([]desc.Product, bool)
	}
)

// New creates bot instance
func New(logger *zap.Logger, cache inMemoryCache, rdb *redis.Client, token string, opts ...Option) (*bot, error) {
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
		ctx:       context.Background(),
		BotAPI:    api,
		logger:    logger.Named("bot"),
		cache:     cache,
		rdb:       rdb,
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
