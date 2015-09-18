package lib

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
)

type QandA struct {
	Details *Details
	lang    struct {
		QandA question `json:"QandA"`
		Thank string   `json:"thank"`
	}
}

func (this *QandA) loadLanguage() {

	lifeJson, err := ioutil.ReadFile("./life.json")
	if err != nil {
		log.Panic("could not load QandA language file", err.Error())
	}

	err = json.Unmarshal(lifeJson, &this.lang)
	if err != nil {
		log.Panic("could not load QandA language ", err.Error())
	}
}

func NewQandA(d *Details) *QandA {
	bg := &QandA{Details: d}
	bg.loadLanguage()
	return bg
}

func (this *QandA) Run(input <-chan telebot.Message) {
	for msg := range input {
		this.Details.User = msg.Chat
		this.ask(msg)
	}
}

func (this *QandA) ask(msg telebot.Message) {
	log.Println("answer was", msg.Text)
	nextQ := this.lang.QandA.find(msg.Text)
	if nextQ == nil {
		this.Details.send(this.lang.Thank)
		log.Println("all done now")
		return
	}
	log.Println("asking ...", nextQ.Question, "labeled:", nextQ.Label)

	this.Details.sendWithKeyboard(nextQ.Question, nextQ.keyboard())
	return
}
