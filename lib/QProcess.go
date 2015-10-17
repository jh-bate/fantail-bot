package lib

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
)

const chat_cmd, ask_cmd, tell_cmd, help_cmd = "/chat", "/ask", "/tell", "/help"

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
		this.saveAndFindNext(msg).andAsk()
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

func (this *QProcess) saveThis(msg telebot.Message) {
	if strings.Contains(msg.Text, ask_cmd) {
		if len(strings.SplitAfter(msg.Text, ask_cmd)) > 0 {
			log.Println("Saving...", ask_cmd)
			this.Details.save(msg)
		}
		return
	} else if strings.Contains(msg.Text, tell_cmd) {
		if len(strings.SplitAfter(msg.Text, tell_cmd)) > 0 {
			log.Println("Saving...", tell_cmd)
			this.Details.save(msg)
		}
		return
	} else if strings.Contains(msg.Text, help_cmd) {
		if len(strings.SplitAfter(msg.Text, help_cmd)) > 0 {
			log.Println("Saving...", help_cmd)
			this.Details.save(msg)
		}
		return
	}
	log.Println("Saving...", chat_cmd)
	this.Details.save(msg)
	return
}

func (this *QProcess) saveAndFindNext(msg telebot.Message) *QProcess {
	this.next = nil

	if strings.Contains(msg.Text, chat_cmd) ||
		strings.Contains(msg.Text, ask_cmd) ||
		strings.Contains(msg.Text, tell_cmd) {
		this.saveThis(msg)
		this.loadLanguage("thank")
		this.next = this.lang.questions[0]
	} else if strings.Contains(msg.Text, chat_cmd) {
		this.loadLanguage("chat")
		this.next = this.lang.questions[0]
	} else {
		for i := range this.lang.questions {
			for a := range this.lang.questions[i].RelatesTo.Answers {
				if this.lang.questions[i].RelatesTo.Answers[a] == msg.Text {
					//was the answer a remainder to save?
					if this.lang.questions[i].RelatesTo.Save {
						this.saveThis(msg)
					}
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
		keyboard = append(keyboard, []string{this.next.PossibleAnswers[i]})
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
