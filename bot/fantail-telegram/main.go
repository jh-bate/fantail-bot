package main

import (
	"log"
	"os"
	"time"

	"github.com/jh-bate/fantail-bot/bot"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
)

type telegram_bot struct {
	t *telebot.Bot
}

var tBot *telegram_bot

func init() {
	botToken := os.Getenv("TELEGRAM_TOKEN")

	if botToken == "" {
		log.Fatal("$TELEGRAM_TOKEN must be set")
	}
	ourBot, err := telebot.NewBot(botToken)
	if err != nil {
		log.Fatal("Bot setup failed: ", err.Error())
	}
	tBot = &telegram_bot{t: ourBot}
}

func (b *telegram_bot) Listen(subscription chan<- *bot.Payload) {

	messages := make(chan telebot.Message)
	b.t.Listen(messages, 1*time.Second)

	for msg := range messages {
		subscription <- bot.New(msg.Sender.ID, msg.Text, msg.Time())
	}

}

func (b *telegram_bot) SendMessage(recipientId int, message string) error {
	return b.t.SendMessage(telebot.User{ID: recipientId}, message, nil)
}

func main() {

	session := bot.NewSession(tBot)
	sub := make(chan *bot.Payload)
	tBot.Listen(sub)

	for payload := range sub {
		session.Respond(payload)
	}

}
