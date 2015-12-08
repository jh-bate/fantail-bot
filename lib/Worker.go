package lib

import (
	"log"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
)

type (
	Worker struct {
		*session
		followup *FollowUp
	}
)

func NewWorker(b *telebot.Bot, s Store) *Worker {
	sess := newSession(b, s)
	q := &Worker{session: sess, followup: NewFollowUp(sess)}
	return q
}

func (this *Worker) ProcessMessages(input <-chan telebot.Message) {

	prevActionName := ""

	for msg := range input {

		log.Println("prev name", prevActionName)

		this.session.setIncoming(newIncoming(msg))

		this.session.getActionForSent(prevActionName).firstUp().askQuestion()

		if update, name := this.session.getActionNameForSent(); update {
			prevActionName = name
		}

		log.Println("prev name after", prevActionName)
	}
}

func (this *Worker) DoFollowUp() {
	log.Println("running the followup process ...")
	this.followup.Start()
}
