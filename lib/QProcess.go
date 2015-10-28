package lib

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
)

const chat_cmd, say_cmd, remind_cmd, show_cmd = "/chat", "/say", "/remind", "/show"

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

func (this *QProcess) showReminders() *QProcess {
	reminders, err := this.Details.Storage.GetCurrentTodos(fmt.Sprintf("%d", this.Details.User.ID))
	if err != nil {
		log.Println("Error trying to get users reminders", err.Error())
	} else {

		for i := range reminders {
			if reminders[i].RemindToday() {
				this.Details.send(reminders[i].Text)
				reminders[i].SetNextReminder()
				this.Details.Storage.Save(fmt.Sprintf("%d", this.Details.User.ID), reminders[i])
			}
		}
	}
	return this
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

func (this *QProcess) saveAndFindNext(msg telebot.Message) *QProcess {
	this.next = nil

	if hasSubmisson(msg.Text, say_cmd) {
		this.Details.save(msg, say_cmd)
		langFile := strings.SplitAfter(msg.Text, "/")[1]
		langFile = strings.Fields(langFile)[0]
		log.Println("loading ...", langFile)
		this.loadLanguage(langFile)
		this.next = this.lang.questions[len(this.lang.questions)-1]
	} else if hasSubmisson(msg.Text, remind_cmd) {
		this.Details.saveReminder(msg)
		langFile := strings.SplitAfter(msg.Text, "/")[1]
		langFile = strings.Fields(langFile)[0]
		log.Println("loading ...", langFile)
		this.loadLanguage(langFile)
		this.next = this.lang.questions[len(this.lang.questions)-1]
	} else if isCmd(msg.Text, say_cmd, remind_cmd, chat_cmd) {
		langFile := strings.SplitAfter(msg.Text, "/")[1]
		log.Println("loading ...", langFile)
		this.loadLanguage(langFile)
		this.next = this.lang.questions[0]
	} else if isCmd(msg.Text, show_cmd) {
		this.showReminders()
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
						this.Details.save(msg, chat_cmd, this.lang.questions[i].RelatesTo.SaveTag)
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
