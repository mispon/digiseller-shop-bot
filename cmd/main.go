package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	xsbot "github.com/mispon/digiseller-shop-bot/bot"
	"github.com/mispon/digiseller-shop-bot/bot/cache"
	uhttp "github.com/mispon/digiseller-shop-bot/utils/http"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	token         = flag.String("token", "", "-token=qwerty")
	sellerId      = flag.String("seller-id", "", "-seller-id=12345")
	debug         = flag.Bool("debug", false, "-debug=true")
	loadCache     = flag.Bool("load-cache", true, "-load-cache=false")
	chatsFilePath = flag.String("chats-file", "chats.txt", "-chats=bot/chats.txt")
	searchUrl     = flag.String("search-url", "", "-search-url=http://localhost:8080")
	uProduct      = flag.Int("uproduct", 0, "-uproduct=12345")
	uProductOpt   = flag.Int("uproductopt", 0, "-uproductopt=12345")
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

	if *chatsFilePath == "" {
		log.Fatal("chats file path is not specified")
	}

	if *searchUrl != "" {
		if !uhttp.IsValidUrl(*searchUrl) {
			log.Fatalf("invalid search url: %s", *searchUrl)
		}
		if *uProduct == 0 {
			log.Fatal("universal product id is not specified")
		} else if *uProductOpt == 0 {
			log.Fatal("universal product option is not specified")
		}
	}

	chatsFile, err := os.OpenFile(*chatsFilePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("failed to open chats file: %v", err)
	}

	logger := mustLogger(*debug)

	botCache, err := cache.New(logger, *sellerId, *loadCache)
	if err != nil {
		log.Fatal(err)
	}

	bot, err := xsbot.New(logger, botCache, chatsFile, *token,
		xsbot.WithSeller(*sellerId),
		xsbot.WithDebug(*debug),
		xsbot.WithSearch(*searchUrl, *uProduct, *uProductOpt),
	)
	if err != nil {
		log.Fatal("failed to create bot", err)
	}

	go bot.Run()

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)

	<-stopCh

	bot.Stop()
	chatsFile.Close()

	logger.Info("Bot gracefully stopped")
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
