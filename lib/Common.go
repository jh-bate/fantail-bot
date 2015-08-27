package lib

import (
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
		Func    func(incoming telebot.Message)
		ToBeRun bool
	}

	Parts []*Part

	option struct {
		Text             string   `json:"text"`
		Feedback         []string `json:"feedback"`
		FollowUpQuestion []string `json:"followUp"`
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
)

func getLangText(opts []string) string {
	return opts[rand.Intn(len(opts))]
}

func makeKeyBoard(keys ...string) [][]string {
	keyboard := [][]string{}
	for i := range keys {
		keyboard = append(keyboard, []string{keys[i]})
	}
	return keyboard
}

func (d *Details) send(msg string) {
	d.Bot.SendMessage(
		d.User,
		msg,
		nil,
	)
	d.takeThoughtfulPause()
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
	time.Sleep(2 * time.Second)
	return
}
