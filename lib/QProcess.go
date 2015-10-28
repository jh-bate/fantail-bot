package lib

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
)

const chat_cmd, remind_cmd, show_cmd, free_form = "/chat", "/remind", "/show", "/free"

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
		this.quickWin(msg).loadScript(msg).findNextQuestion(msg).andAsk()
	}
}

func (this *QProcess) quickWin(msg telebot.Message) *QProcess {

	if hasSubmisson(msg.Text, remind_cmd) {
		log.Println("making submisson ", msg.Text)
		this.Details.saveAsReminder(msg)
	} else if strings.Contains(msg.Text, show_cmd) {
		log.Println("showing reminders ", msg.Text)
		r := this.Details.getReminders(fmt.Sprintf("%d", this.Details.User.ID))
		for i := range r {
			if r[i].RemindToday() {
				this.Details.send(r[i].Text)

			}
		}
	}
	return this
}

func (this *QProcess) loadScript(msg telebot.Message) *QProcess {

	scriptName := "default"

	if isCmd(msg.Text, remind_cmd, chat_cmd, show_cmd) {
		words := strings.Fields(msg.Text)
		scriptName = strings.SplitAfter(words[0], "/")[1]

		file, err := os.Open(fmt.Sprintf("./config/%s.json", scriptName))
		if err != nil {
			log.Panic("could not load QandA language file ", err.Error())
		}
		err = json.NewDecoder(file).Decode(&this.lang)
		if err != nil {
			log.Panic("could not decode QandA ", err.Error())
		}
	}

	return this
}

func (this *QProcess) findNextQuestion(msg telebot.Message) *QProcess {
	this.next = nil

	if isCmd(msg.Text, remind_cmd, chat_cmd) {
		//start at the beginning
		this.next = this.lang.questions[0]
		return this
	} else if hasSubmisson(msg.Text, remind_cmd) {
		//straight to the end
		this.next = this.lang.questions[len(this.lang.questions)-1]
		return this
	}

	//find the next question
	for i := range this.lang.questions {
		log.Println("looking next q ...")
		for a := range this.lang.questions[i].RelatesTo.Answers {
			if this.lang.questions[i].RelatesTo.Answers[a] == msg.Text {
				//was the answer a remainder to save?
				if this.lang.questions[i].RelatesTo.Save {
					this.Details.save(msg, chat_cmd, this.lang.questions[i].RelatesTo.SaveTag)
				}
				this.next = this.lang.questions[i]
				//all done
				return this
			}
		}
	}

	//not sure so we will just save it
	log.Println("unknown so will save as", free_form)
	this.Details.save(msg, free_form)

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
