package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jh-bate/fantail-bot/bot"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/nlopes/slack"
)

type slack_bot struct {
	rtm *slack.RTM
}

var sBot *slack_bot

func init() {
	token := os.Getenv("SLACK_TOKEN")
	if token == "" {
		log.Fatal("$SLACK_TOKEN must be set")
	}

	api := slack.New(token)
	api.SetDebug(true)
	sBot = &slack_bot{rtm: api.NewRTM()}
}

func (b *slack_bot) Listen(subscription chan<- *bot.Payload) {

	go b.rtm.ManageConnection()

	for {
		select {
		case msg := <-b.rtm.IncomingEvents:
			fmt.Print("Event Received: ")
			switch ev := msg.Data.(type) {

			case *slack.MessageEvent:
				fmt.Printf("Message: %v\n", ev)
				when, err := time.Parse(time.RFC3339, ev.Msg.Timestamp)
				if err == nil {
					subscription <- bot.New(777, ev.Msg.Text, when)
				}

			default:

				// Ignore other events..
				// fmt.Printf("Unexpected: %v\n", msg.Data)
			}
		}
	}

}

func (b *slack_bot) SendMessage(recipientId int, message string) error {
	b.rtm.SendMessage(b.rtm.NewOutgoingMessage(message, "test"))
	return nil
}

func main() {

	session := bot.NewSession(sBot)
	sub := make(chan *bot.Payload)
	sBot.Listen(sub)

	for payload := range sub {
		session.Respond(payload)
	}

}
