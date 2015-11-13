package lib

import (
	"log"
	"strings"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
)

type (
	QProcess struct {
		s *session
	}
)

func NewQProcess(b *telebot.Bot, s *Storage) *QProcess {
	q := &QProcess{s: newSession(b, s)}
	return q
}

func getActionName(msg telebot.Message) (bool, string) {
	log.Println("Check getActionName ", msg.Text)
	if strings.Contains(msg.Text, "/") {
		log.Println("Check getActionName ", strings.Fields(msg.Text)[0])
		return true, strings.Fields(msg.Text)[0]
	}
	return false, ""
}

func (this *QProcess) Run(input <-chan telebot.Message) {

	prevActionName := ""

	for msg := range input {

		log.Println("prev name", prevActionName)

		in := newIncoming(msg)
		in.getAction(this.s, prevActionName).firstUp().askQuestion()
		if update, name := getActionName(msg); update {
			prevActionName = name
		}

		log.Println("prev name after", prevActionName)
	}
}
