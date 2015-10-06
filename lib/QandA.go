package lib

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
)

type QandA struct {
	Details *Details
	lang    struct {
		QandA questions `json:"QandA"`
	}
	toAsk *Question
}

func NewQandA(d *Details) *QandA {
	bg := &QandA{Details: d}
	bg.loadLanguage()
	return bg
}

func (this *QandA) Run(input <-chan telebot.Message) {
	for msg := range input {
		this.Details.User = msg.Chat
		this.getNext(msg.Text).andAsk()
	}
}

func (this *QandA) loadLanguage() {

	file, err := os.Open("./config/whys.json")
	if err != nil {
		log.Panic("could not load QandA language file ", err.Error())
	}
	err = json.NewDecoder(file).Decode(&this.lang)
	if err != nil {
		log.Panic("could not decode QandA ", err.Error())
	}
}

func (this *QandA) getNext(prevAnswer string) *QandA {
	this.toAsk = nil

	if strings.Contains(prevAnswer, "/chat") {
		this.toAsk = this.lang.QandA[0]
	} else {
		for i := range this.lang.QandA {
			for a := range this.lang.QandA[i].PossibleAnswers {
				if this.lang.QandA[i].PossibleAnswers[a] == prevAnswer {
					nextQNum := i + 1
					if len(this.lang.QandA) <= nextQNum {
						this.toAsk = this.lang.QandA[nextQNum]
						return this
					}
				}
			}
		}
	}
	return this
}

func (this *QandA) makeKeyboard() Keyboard {
	keyboard := Keyboard{}
	for i := range this.toAsk.PossibleAnswers {
		keyboard = append(keyboard, []string{this.toAsk.PossibleAnswers[i]})
	}
	return keyboard
}

func (this *QandA) andAsk() {
	if this.toAsk != nil {
		//context
		for i := range this.toAsk.Context {
			this.Details.send(fmt.Sprintf(this.toAsk.Context[i], this.Details.User.FirstName))
		}
		//the actual question
		this.Details.sendWithKeyboard(fmt.Sprintf(this.toAsk.QuestionText, this.Details.User.FirstName), this.makeKeyboard())
	}
	return
}
