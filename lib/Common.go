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
	/*question struct {
		Label    string      `json:"label"`
		Question string      `json:"question"`
		Children []*question `json:"children"`
	}*/

	questions struct {
		Label     string      `json:"label"`
		Questions []string    `json:"questions"`
		Children  []*question `json:"children"`
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

func (q *questions) hasChildren() bool {
	return q.Children != nil && len(q.Children) > 0
}

func (q *questions) find(label string) *questions {
	if q.hasChildren() {

		if q.Label == label {
			log.Println("at the top so return this one's children")
			return q
		}

		for i := range q.Children {
			if q.Children[i].Label == label {
				return q.Children[i]
			} else if q.Children[i].hasChildren() {
				match := q.Children[i].find(label)
				if match != nil {
					return match
				}
			}
		}
	}
	log.Println("nothing else found")
	return nil
}

func (q *questions) makeKeyboard() Keyboard {
	keyboard := Keyboard{}
	for i := range q.Children {
		keyboard = append(keyboard, []string{q.Children[i].Label})
	}
	return keyboard
}

func (d *Details) send(m string) {
	d.takeThoughtfulPause()
	d.Bot.SendMessage(
		d.User,
		msg,
		nil,
	)
	return
}

func (d *Details) askQuestion(q *questions) {
	for i := range q.Questions {
		if i == len(q.Questions) {
			d.send(q.Questions[i])
		} else {
			d.sendWithKeyboard(q.Questions[i], q.makeKeyboard())
		}
	}
	return
}

func (d *Details) sendWithKeyboard(msg string, kb Keyboard) {
	d.takeThoughtfulPause()
	d.Bot.SendMessage(
		d.User,
		msg,
		&telebot.SendOptions{
			ReplyMarkup: telebot.ReplyMarkup{
				ForceReply:      true,
				CustomKeyboard:  kb,
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
