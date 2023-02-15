package main

import (
	"flag"
	"go.uber.org/zap/zapcore"
	"log"

	xsbot "github.com/mispon/xbox-store-bot/bot"
	"github.com/mispon/xbox-store-bot/bot/cache"
	"go.uber.org/zap"
)

var (
	token     = flag.String("token", "", "-token=qwerty")
	sellerId  = flag.String("seller-id", "", "-seller-id=12345")
	debug     = flag.Bool("debug", false, "-debug=true")
	loadCache = flag.Bool("load-cache", true, "-load-cache=false")
)

func init() {
	flag.Parse()
}

func main() {
	if *token == "" {
		log.Fatal("bot token is not specified")
	}

	if *sellerId == "" {
		log.Fatal("seller id is not specified")
	}

	logger := mustLogger(*debug)

	botCache, err := cache.New(logger, *sellerId, *loadCache)
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

func mustLogger(debug bool) *zap.Logger {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.DisableStacktrace = true

	logLevel := zapcore.InfoLevel
	if debug {
		logLevel = zapcore.DebugLevel
	}
	cfg.Level = zap.NewAtomicLevelAt(logLevel)

	logger, err := cfg.Build()
	if err != nil {
		log.Fatal(err)
	}

	return logger
}
