package lib

import (
	"log"
	"time"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
)

const (
	typing_action = "typing"
)

type (
	question struct {
		Label    string      `json:"label"`
		Question string      `json:"question"`
		Children []*question `json:"children"`
	}

	Details struct {
		Bot  *telebot.Bot
		User telebot.User
	}

	Process interface {
		Run(input <-chan telebot.Message)
	}

	Keyboard [][]string
)

func (q *question) hasChildren() bool {
	return q.Children != nil && len(q.Children) > 0
}

func (q *question) find(label string) *question {
	if q.hasChildren() {

		if q.Label == label {
			log.Println("at the top so return this one's children")
			//return the first as we are the top
			return q
		}

		for i := range q.Children {
			if q.Children[i].Label == label {
				//log.Println("found ", label)
				return q.Children[i]
			} else if q.Children[i].hasChildren() {
				//log.Println("check children ", label)
				match := q.Children[i].find(label)
				if match != nil {
					return match
				}
			}
		}
	}
	return nil
}

func (q *question) keyboard() [][]string {
	keyboard := [][]string{}
	for i := range q.Children {
		keyboard = append(keyboard, []string{q.Children[i].Label})
	}
	return keyboard
}

func (d *Details) send(msg string) {
	d.takeThoughtfulPause()
	d.Bot.SendMessage(
		d.User,
		msg,
		nil,
	)
	return
}

func (d *Details) sendWithKeyboard(msg string, keyboard [][]string) {
	d.takeThoughtfulPause()
	d.Bot.SendMessage(
		d.User,
		msg,
		&telebot.SendOptions{
			ReplyMarkup: telebot.ReplyMarkup{
				ForceReply:      true,
				CustomKeyboard:  keyboard,
				ResizeKeyboard:  false,
				OneTimeKeyboard: true,
			},
		},
	)
	return
}

func (d *Details) takeThoughtfulPause() {
	d.Bot.SendChatAction(d.User, typing_action)
	time.Sleep(1 * time.Second)
	return
}
