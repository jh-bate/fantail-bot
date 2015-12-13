package main

import (
	"log"
	"os"
	"time"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
	"github.com/jh-bate/fantail-bot/incoming"
)

func main() {

	botToken := os.Getenv("BOT_TOKEN")

	if botToken == "" {
		log.Fatal("$BOT_TOKEN must be set")
	}

	ourBot, err := telebot.NewBot(botToken)
	if err != nil {
		log.Fatal("Bot setup failed: ", err.Error())
	}

	messages := make(chan telebot.Message)
	ourBot.Listen(messages, 1*time.Second)

	session := incoming.NewSession(ourBot)

	for msg := range messages {
		session.Respond(msg)
	}

}
