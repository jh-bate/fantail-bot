package lib

import (
	"strings"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
)

const (
	default_script = "default"
	stickers_chat  = "chat"
)

type (
	Info struct {
		App       []string `json:"appInfo"`
		Reminders []string `json:"remindersInfo"`
		Chat      []string `json:"chatInfo"`
		Said      []string `json:"saidInfo"`
	}

	QProcess struct {
		s *session
	}
)

func NewQProcess(b *telebot.Bot, s *Storage) *QProcess {
	q := &QProcess{s: newSession(b, s)}
	return q
}

func getActionName(msg telebot.Message) (bool, string) {
	if strings.Contains(msg.Text, "/") {
		return true, strings.Fields(msg.Text)[0]
	}
	return false, ""
}

func (this *QProcess) Run(input <-chan telebot.Message) {

	prevActionName := ""

	for msg := range input {

		in := newIncoming(msg)
		in.getAction(this.s, prevActionName).firstUp().askQuestion()
		if update, name := getActionName(msg); update {
			prevActionName = name
		}
	}
}
