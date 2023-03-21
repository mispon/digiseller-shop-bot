package bot

import (
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mispon/xbox-store-bot/bot/desc"
	"go.uber.org/zap"
)

var admins = map[string]struct{}{
	"noobmaster111": {},
	"Mispon":        {},
	"kotovro":       {},
}

type searchParams struct {
	category  string
	query     string
	messageID int
}
type chat struct {
	searchParams searchParams
}

type (
	bot struct {
		*tgbotapi.BotAPI

		logger *zap.Logger
		cache  inMemoryCache
		opts   options

		commands  map[commandKey]commandEntity
		callbacks map[callbackType]callbackFn

		chatsMu   sync.RWMutex
		chats     map[int64]chat
		chatsFile *os.File

		convRatesMu     sync.RWMutex
		conversionRates map[string]float64

		client *http.Client
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
func New(logger *zap.Logger, cache inMemoryCache, chatsFile *os.File, token string, opts ...Option) (*bot, error) {
	api, aErr := tgbotapi.NewBotAPI(token)
	if aErr != nil {
		return nil, aErr
	}

	var bo options
	for _, optFn := range opts {
		optFn(&bo)
	}
	api.Debug = bo.debug

	logger = logger.Named("bot")
	b := &bot{
		BotAPI:    api,
		logger:    logger,
		cache:     cache,
		opts:      bo,
		commands:  make(map[commandKey]commandEntity),
		callbacks: make(map[callbackType]callbackFn),
		chats:     readChats(logger, chatsFile),
		chatsFile: chatsFile,
		client:    http.DefaultClient,
		conversionRates: map[string]float64{
			"ARS": 0.75,
			"TRY": 6,
		},
	}

	if err := b.initCommands(); err != nil {
		return nil, err
	}
	b.initCallbacks()

	go b.chatsSaver()

	b.logger.Info("bot created", zap.String("username", api.Self.UserName))
	return b, nil
}

func (b *bot) apiRequest(c tgbotapi.Chattable) error {
	_, err := b.Request(c)
	return err
}

func readChats(logger *zap.Logger, file *os.File) map[int64]chat {
	chats := make(map[int64]chat)

	bytes, err := io.ReadAll(file)
	if err != nil {
		logger.Error("failed to read chats", zap.Error(err))
		return chats
	}

	chatIDs := strings.Split(string(bytes), ",")
	for _, val := range chatIDs {
		if val == "" {
			continue
		}

		chatID, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			logger.Error("failed to parse chat id", zap.String("val", val), zap.Error(err))
			continue
		}
		chats[chatID] = chat{searchParams{}}
	}

	logger.Info("chats loaded from file", zap.Int("count", len(chats)))
	return chats
}

func (b *bot) writeChats() {
	b.chatsMu.RLock()
	defer b.chatsMu.RUnlock()

	if err := b.chatsFile.Truncate(0); err != nil {
		b.logger.Error("failed to truncate chats file", zap.Error(err))
		return
	}
	if _, err := b.chatsFile.Seek(0, 0); err != nil {
		b.logger.Error("failed to seek chats file", zap.Error(err))
		return
	}

	chatIDs := make([]string, 0, len(b.chats))
	for chatID := range b.chats {
		chatIDs = append(chatIDs, strconv.FormatInt(chatID, 10))
	}
	data := strings.Join(chatIDs, ",")

	if _, err := b.chatsFile.Write([]byte(data)); err != nil {
		b.logger.Error("failed to write chats", zap.Error(err))
	}

	b.logger.Info("chats saved to file", zap.Int("count", len(b.chats)))
}

func (b *bot) addChat(chatID int64) {
	b.chatsMu.Lock()
	defer b.chatsMu.Unlock()

	b.chats[chatID] = chat{searchParams{}}
}

func (b *bot) chatExist(chatID int64) bool {
	b.chatsMu.RLock()
	defer b.chatsMu.RUnlock()

	if _, ok := b.chats[chatID]; ok {
		return true
	}
	return false
}

func (b *bot) chatsSaver() {
	ticker := time.NewTicker(10 * time.Minute)
	for {
		<-ticker.C
		b.writeChats()
	}
}

func (b *bot) getConversionRate(currency string) float64 {
	b.convRatesMu.RLock()
	defer b.convRatesMu.RUnlock()
	return b.conversionRates[currency]
}
