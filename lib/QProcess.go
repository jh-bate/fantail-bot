package lib

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
)

type QProcess struct {
	Details *Details
	lang    struct {
		questions `json:"QandA"`
	}
	next *Question
}

func NewQProcess(d *Details) *QProcess {
	return &QProcess{Details: d}
}

func (this *QProcess) Run(input <-chan telebot.Message) {
	for msg := range input {
		this.Details.User = msg.Chat
		this.getNext(msg.Text).andAsk()
	}
}

func (this *QProcess) loadLanguage(name string) {

	file, err := os.Open(fmt.Sprintf("./config/%s.json", name))
	if err != nil {
		log.Panic("could not load QandA language file ", err.Error())
	}
	err = json.NewDecoder(file).Decode(&this.lang)
	if err != nil {
		log.Panic("could not decode QandA ", err.Error())
	}
}

func (this *QProcess) getNext(prevAnswer string) *QProcess {
	this.next = nil

	if strings.Contains(prevAnswer, "/") {
		log.Println("start fresh for command", prevAnswer)
		cmdName := strings.Split(prevAnswer, "/")[0]
		cmdName = strings.Split(cmdName, " ")[0]
		this.loadLanguage(cmdName)
		this.next = this.lang.questions[0]
	} else {
		for i := range this.lang.questions {
			for a := range this.lang.questions[i].RelatesTo {
				if this.lang.questions[i].RelatesTo[a] == prevAnswer {
					this.next = this.lang.questions[i]
					return this
				}
			}
		}
	}
	return this
}

func (this *QProcess) makeKeyboard() Keyboard {
	keyboard := Keyboard{}
	for i := range this.next.PossibleAnswers {
		keyboard = append(keyboard, []string{this.next.PossibleAnswers[i].Text})
	}
	return keyboard
}

func (this *QProcess) andAsk() {
	if this.next != nil {
		//context
		for i := range this.next.Context {
			this.Details.send(this.next.Context[i])
		}
		//the actual question
		this.Details.sendWithKeyboard(this.next.QuestionText, this.makeKeyboard())
	}
	return
}
