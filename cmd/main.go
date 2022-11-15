package main

import (
	"errors"
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func main() {
	logger := mustLogger()

	bot, err := createBot(logger)
	if err != nil {
		logger.Fatal("failed to create bot", zap.Error(err))
	}

	updatesCfg := tgbotapi.UpdateConfig{
		Offset:  0,
		Timeout: 10,
	}
	for upd := range bot.GetUpdatesChan(updatesCfg) {
		if upd.Message != nil {
			text := upd.Message.Text
			logger.Info("receive msg", zap.String("msg", text))
		}

		if upd.CallbackQuery != nil {

		}
	}
}

func mustLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	return logger
}

func createBot(logger *zap.Logger) (*tgbotapi.BotAPI, error) {
	token, ok := os.LookupEnv("XSB_TOKEN")
	if !ok {
		return nil, errors.New("bot token not specified")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	if val, err := strconv.ParseBool(os.Getenv("XSB_DEBUG")); err == nil {
		bot.Debug = val
	}

	menuCfg := tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{
			Command:     "/test",
			Description: "test command",
		},
	)
	_, _ = bot.Request(menuCfg)

	logger.Info("bot created", zap.String("username", bot.Self.UserName))
	return bot, nil
}
