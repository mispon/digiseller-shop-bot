package main

import (
	"flag"
	"log"

	xsbot "github.com/mispon/xbox-store-bot/bot"
	"github.com/mispon/xbox-store-bot/bot/cache"
	"go.uber.org/zap"
)

var (
	token    = flag.String("token", "", "-token=qwerty")
	amToken  = flag.String("am_token", "", "-am_token=qwerty")
	sellerId = flag.String("seller-id", "", "-seller-id=12345")
	debug    = flag.Bool("debug", false, "-debug=true")
)

func init() {
	flag.Parse()
}

func main() {
	if *token == "" {
		log.Fatal("bot token is not specified")
	}

	if *amToken == "" {
		log.Fatal("app metrika token is not specified")
	}

	if *sellerId == "" {
		log.Fatal("seller id is not specified")
	}

	logger := mustLogger()

	botCache, err := cache.New(logger, *sellerId)
	if err != nil {
		log.Fatal(err)
	}

	opts := []xsbot.Option{
		xsbot.WithSeller(*sellerId),
		xsbot.WithDebug(*debug),
	}
	bot, err := xsbot.New(logger, botCache, *token, opts...)
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
