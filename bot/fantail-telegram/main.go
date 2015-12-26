package main

import (
	"log"
	"os"
	"time"

	"github.com/jh-bate/fantail-bot/bot"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
)

type telegram struct {
	t *telebot.Bot
}

func newTelegramBot() *telegram {
	botToken := os.Getenv("BOT_TOKEN")

	if botToken == "" {
		log.Fatal("$BOT_TOKEN must be set")
	}
	ourBot, err := telebot.NewBot(botToken)
	if err != nil {
		log.Fatal("Bot setup failed: ", err.Error())
	}
	return &telegram{t: ourBot}
}

func (b *telegram) Listen(subscription chan<- *bot.Payload) {

	messages := make(chan telebot.Message)
	b.t.Listen(messages, 1*time.Second)

	for msg := range messages {
		subscription <- bot.New(msg.Sender.ID, msg.Text, msg.Time())
	}

}

func (b *telegram) SendMessage(recipientId int, message string) error {
	return b.t.SendMessage(telebot.User{ID: recipientId}, message, nil)
}

func main() {

	telegramBot := newTelegramBot()
	session := bot.NewSession(telegramBot)
	sub := make(chan *bot.Payload)
	telegramBot.Listen(sub)

	for payload := range sub {
		session.Respond(payload)
	}

}
