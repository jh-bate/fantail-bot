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

func hasSubmisson(txt string, cmds ...string) bool {
	if isCmd(txt, cmds...) {
		for i := range cmds {
			log.Println("Check if submisson", txt)
			if len(strings.SplitAfter(txt, cmds[i])) > 1 {
				return true
			}
		}
	}
	return false
}

func isCmd(txt string, cmds ...string) bool {
	log.Println("Check if cmd", txt)
	for i := range cmds {
		if strings.Contains(txt, cmds[i]) {
			return true
		}
	}
	return false
}

func (this *QProcess) saveAndFindNext(msg telebot.Message) *QProcess {
	this.next = nil

	if hasSubmisson(msg.Text, help_cmd, ask_cmd, tell_cmd) {
		this.Details.save(msg)
		langFile := strings.SplitAfter(msg.Text, "/")[1]
		log.Println("loading ...", langFile)
		this.loadLanguage(langFile)
		this.next = this.lang.questions[len(this.lang.questions)-1]
	} else if isCmd(msg.Text, help_cmd, ask_cmd, tell_cmd, chat_cmd) {
		langFile := strings.SplitAfter(msg.Text, "/")[1]
		log.Println("loading ...", langFile)
		this.loadLanguage(langFile)
		this.next = this.lang.questions[0]
	} else {
		for i := range this.lang.questions {
			for a := range this.lang.questions[i].RelatesTo.Answers {
				if this.lang.questions[i].RelatesTo.Answers[a] == msg.Text {
					//was the answer a remainder to save?
					if this.lang.questions[i].RelatesTo.Save {
						this.Details.save(msg)
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
