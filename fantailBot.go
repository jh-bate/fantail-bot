package main

import (
	"log"
	"os"
	"time"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
	"github.com/jh-bate/fantail-bot/lib"
)

type (
	fantailBot struct {
		bot   *telebot.Bot
		store *lib.Storage
	}
)

func newFantailBot() *fantailBot {
	botToken := os.Getenv("BOT_TOKEN")

	if botToken == "" {
		log.Fatal("$BOT_TOKEN must be set")
	}

	bot, err := telebot.NewBot(botToken)
	if err != nil {
		return nil
	}

	return &fantailBot{bot: bot, store: lib.NewStorage()}
}

func main() {

	fBot := newFantailBot()
	messages := make(chan telebot.Message)
	fBot.bot.Listen(messages, 1*time.Second)

	q := lib.NewQProcess(&lib.Details{Bot: fBot.bot, Storage: fBot.store})
	q.Run(messages)
}
