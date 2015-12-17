package incoming

import (
	"fmt"
	"strings"
	"time"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
	"github.com/jh-bate/fantail-bot/user"
)

const (
	typing_action = "typing"
)

type (
	Session struct {
		actionRunAs string
		*telebot.Bot
	}

	Keyboard [][]string
)

var monitor *FollowUp

func NewSession(ourBot *telebot.Bot) *Session {

	s := &Session{Bot: ourBot}

	monitor = NewFollowUp(s)
	monitor.Start()
	return s
}

func (s *Session) Respond(msg telebot.Message) {

	//TODO should not just keep saving...
	sessionUser := user.New(msg.Sender.ID)
	sessionUser.Save()

	a := NewAction(New(msg), s.actionRunAs, s)
	s.actionRunAs = a.getName()
	a.firstUp().askQuestion()
}

func (s *Session) send(to telebot.User, msgs ...string) {

	for i := range msgs {
		s.takeThoughtfulPause(to)

		msg := msgs[i]
		if strings.Contains(msg, "%s") {
			msg = fmt.Sprintf(msg, to.FirstName)
		}

		s.Bot.SendMessage(
			to,
			msg,
			&telebot.SendOptions{ParseMode: telebot.ModeMarkdown},
		)
	}
	return
}

func (s *Session) sendWithKeyboard(to telebot.User, msg string, kb Keyboard) {
	s.takeThoughtfulPause(to)

	if strings.Contains(msg, "%s") {
		msg = fmt.Sprintf(msg, to.FirstName)
	}

	s.Bot.SendMessage(
		to,
		msg,
		&telebot.SendOptions{
			ReplyMarkup: telebot.ReplyMarkup{
				ForceReply:      false,
				CustomKeyboard:  kb,
				ResizeKeyboard:  false,
				OneTimeKeyboard: true,
			},
			ParseMode: telebot.ModeMarkdown,
		},
	)
	return
}

func (s *Session) takeThoughtfulPause(to telebot.User) {
	s.Bot.SendChatAction(to, typing_action)
	time.Sleep(1 * time.Second)
	return
}
