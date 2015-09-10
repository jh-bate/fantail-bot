package lib

import (
	"log"
	"math/rand"
	"time"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
)

const (
	yes_text = "Yeah"
	no_text  = "Nope"
	bye_text = "See you %s"

	typing_action = "typing"
)

type (
	Part struct {
		Func     func(incoming telebot.Message)
		Msg      func(msg string)
		Keyboard func(msg string, keys [][]string)
		ToBeRun  bool
	}

	Parts []*Part

	Step struct {
		SendMsg      func(incoming telebot.Message)
		SendKeyboard func(incoming telebot.Message)
	}

	Steps []*Step

	option struct {
		Text             string   `json:"text"`
		Feedback         []string `json:"feedback"`
		FollowUpQuestion []string `json:"followUp"`
	}

	question struct {
		Label    string      `json:"label"`
		Question string      `json:"question"`
		Children []*question `json:"children"`
	}

	questionYesNo struct {
		Question string `json:"question"`
		yesNo    `json:"answers"`
	}

	yesNo struct {
		Yes string `json:"yes"`
		No  string `json:"no"`
	}

	Details struct {
		Bot  *telebot.Bot
		User telebot.User
	}

	Process interface {
		Run(telebot.Message)
		CanRun() bool
	}

	Process2 interface {
		Run(input <-chan telebot.Message)
	}
)

func getLangText(opts []string) string {
	return opts[rand.Intn(len(opts))]
}

func (q *question) hasChildren() bool {
	return q.Children != nil && len(q.Children) > 0
}

func (q *question) findChild(label string) *question {
	if q.hasChildren() {
		for i := range q.Children {
			if q.Children[i].Label == label {
				log.Println("found ", label)
				return q.Children[i]
			} else if q.Children[i].hasChildren() {
				log.Println("check children ", label)
				match := q.Children[i].findChild(label)
				if match != nil {
					return match
				}
			}

		}
	}
	return nil
}

func (q *question) keyboard() [][]string {
	log.Println("keyboard for ", q.Label)
	keyboard := [][]string{}
	for i := range q.Children {
		keyboard = append(keyboard, []string{q.Children[i].Label})
	}
	return keyboard
}

func makeKeyBoard(keys ...string) [][]string {
	keyboard := [][]string{}
	for i := range keys {
		keyboard = append(keyboard, []string{keys[i]})
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
