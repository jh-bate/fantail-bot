package lib

import (
	"log"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
)

type (
	QProcess struct {
		s *session
		f *FollowUp
	}
)

func NewQProcess(b *telebot.Bot, s *Storage) *QProcess {
	sess := newSession(b, s)
	q := &QProcess{s: sess, f: NewFollowUp(sess)}
	return q
}

func (this *QProcess) Run(input <-chan telebot.Message) {

	prevActionName := ""

	for msg := range input {

		log.Println("prev name", prevActionName)

		//in := newIncoming(msg)

		this.s.setIncoming(newIncoming(msg))

		this.s.getActionForSent(prevActionName).firstUp().askQuestion()

		if update, name := this.s.getActionNameForSent(); update {
			prevActionName = name
		}

		log.Println("prev name after", prevActionName)
	}
}

func (this *QProcess) DoFollowUp() {
	log.Println("running the followup process ...")
	this.f.Start()
}
