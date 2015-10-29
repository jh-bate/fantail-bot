package lib

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
)

const (
	chat_cmd, remind_cmd, reminders_cmd, notes_cmd = "/chat", "/remind", "/reminders", "/notes"

	free_form      = "JUST_SAYING"
	default_script = "default"
)

type QProcess struct {
	Details *Details
	lang    struct {
		questions `json:"QandA"`
	}
	next *Question
	done bool
}

func NewQProcess(d *Details) *QProcess {
	return &QProcess{Details: d}
}

func (this *QProcess) Run(input <-chan telebot.Message) {
	for msg := range input {
		this.Details.User = msg.Chat
		this.quickWin(msg).determineScript(msg).findNextQuestion(msg).andAsk()
	}
}

func (this *QProcess) quickWin(msg telebot.Message) *QProcess {

	/*if hasSubmisson(msg.Text, remind_cmd) {
		log.Println("making submisson ", msg.Text)
		this.Details.saveAsReminder(msg)
	} else */
	if isCmd(msg.Text, reminders_cmd) {
		log.Println("showing reminders ", msg.Text)

		r := this.Details.getReminders(fmt.Sprintf("%d", this.Details.User.ID))
		for i := range r {
			if r[i].RemindToday() {
				this.Details.send(r[i].Text)

			}
		}
	} else if isCmd(msg.Text, notes_cmd) {
		log.Println("showing notes ", msg.Text)

		r := this.Details.getNotes(fmt.Sprintf("%d", this.Details.User.ID))
		for i := range r {
			if r[i].IsCurrent() {
				this.Details.send(r[i].Text)
			}
		}
	}
	return this
}

func (this *QProcess) loadScript(scriptName string) {
	file, err := os.Open(fmt.Sprintf("./config/%s.json", scriptName))
	if err != nil {
		log.Panic("could not load QandA language file ", err.Error())
	}
	err = json.NewDecoder(file).Decode(&this.lang)
	if err != nil {
		log.Panic("could not decode QandA ", err.Error())
	}
}

func (this *QProcess) determineScript(msg telebot.Message) *QProcess {

	if isCmd(msg.Text, chat_cmd, notes_cmd) { //remind_cmd, chat_cmd, reminders_cmd, notes_cmd) {
		words := strings.Fields(msg.Text)
		scriptName := strings.SplitAfter(words[0], "/")[1]
		this.loadScript(scriptName)
	}
	return this
}

func (this *QProcess) findNextQuestion(msg telebot.Message) *QProcess {
	this.next = nil

	/*if this.done || this.lang.questions == nil {
		log.Println("unknown so will save as", free_form)
		this.Details.save(msg, free_form)
		//load default and start at the beginning
		this.loadScript(default_script)
		this.next = this.lang.questions[0]
		this.done = false
		return this
	} else */
	if isCmd(msg.Text, chat_cmd, notes_cmd) { //|| hasSubmisson(msg.Text, remind_cmd) {
		//start at the beginning
		this.next = this.lang.questions[0]
		//this.done = false
		return this
	} else {
		//find the next question
		searched := false
		for i := range this.lang.questions {
			searched = true
			log.Println("looking next q ...")
			for a := range this.lang.questions[i].RelatesTo.Answers {
				if this.lang.questions[i].RelatesTo.Answers[a] == msg.Text {
					//was the answer a remainder to save?
					if this.lang.questions[i].RelatesTo.Save {
						this.Details.save(msg, chat_cmd, this.lang.questions[i].RelatesTo.SaveTag)
					}
					this.next = this.lang.questions[i]
					//this.done = false
					return this
				}
			}
		}
		if !searched {
			log.Println("was just saying so we are just saving")
			this.Details.save(msg, free_form)
			//load default and start at the beginning
			this.loadScript(default_script)
			this.next = this.lang.questions[0]
			//this.done = false
			return this
		}
	}
	log.Println("looks like we are all done!")
	//this.done = true
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
