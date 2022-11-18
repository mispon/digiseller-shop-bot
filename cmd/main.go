package main

import (
	xsbot "github.com/mispon/xbox-store-bot/bot"
	"go.uber.org/zap"
	"log"
	"os"
)

func main() {
	token, ok := os.LookupEnv("XSB_TELEGRAM_TOKEN")
	if !ok {
		log.Fatal("bot token is not specified")
	}

	sellerId, ok := os.LookupEnv("XSB_SELLER_ID")
	if !ok {
		log.Fatal("seller id is not specified")
	}

	logger := mustLogger()

	opts := []xsbot.Option{
		xsbot.WithDebug(os.Getenv("XSB_DEBUG")),
	}
	bot, err := xsbot.New(logger, token, sellerId, opts...)
	if err != nil {
		log.Fatal("failed to create bot", err)
	}

	bot.Run()
}

func mustLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	return logger
}
