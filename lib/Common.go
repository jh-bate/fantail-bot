package lib

import (
	"fmt"
	"strings"
	"time"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
)

const (
	typing_action = "typing"
)

type (
	questionTree struct {
		Label     string          `json:"label"`
		Questions []string        `json:"questions"`
		Children  []*questionTree `json:"children"`
	}

	Question struct {
		Context         []string `json:"context"`
		QuestionText    string   `json:"question"`
		PossibleAnswers []string `json:"answers"`
	}

	questions []*Question

	Details struct {
		Bot  *telebot.Bot
		User telebot.User
	}

	Process interface {
		Run(input <-chan telebot.Message)
	}

	Keyboard [][]string
)

func (d *Details) send(msg string) {
	d.takeThoughtfulPause()

	if strings.Contains(msg, "%s") {
		msg = fmt.Sprintf(msg, d.User.FirstName)
	}

	d.Bot.SendMessage(
		d.User,
		msg,
		nil,
	)
	return
}

func (d *Details) sendWithKeyboard(msg string, kb Keyboard) {
	d.takeThoughtfulPause()

	if strings.Contains(msg, "%s") {
		msg = fmt.Sprintf(msg, d.User.FirstName)
	}

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
