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
	chat_cmd, remind_cmd, say_cmd, reminders_cmd, said_cmd = "/chat", "/remind", "/say", "/reminders", "/said"

	remind_cmd_hint, say_cmd_hint = "/remind me in <days> to <message>", "/say [what you want to say]"

	free_form      = "JUST_SAYING"
	default_script = "default"
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
		this.quickWinFirst(msg).
			determineScript(msg).
			findNextQuestion(msg).
			andAsk()
	}
}

func (this *QProcess) quickWinFirst(msg telebot.Message) *QProcess {

	if hasSubmisson(msg.Text, say_cmd) {
		log.Println("save something said ", msg.Text)
		this.Details.save(msg, say_cmd)
	} else if hasSubmisson(msg.Text, remind_cmd) {
		log.Println("save a reminder ", msg.Text)
		this.Details.saveAsReminder(msg)
	} else if isCmd(msg.Text, reminders_cmd) {
		log.Println("showing reminders ", msg.Text)
		r := this.Details.getReminders(fmt.Sprintf("%d", this.Details.User.ID))
		for i := range r {
			if r[i].RemindToday() {
				this.Details.send(r[i].Text)
			}
		}
	} else if isCmd(msg.Text, said_cmd) {
		log.Println("showing things said ", msg.Text)
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

	if hasSubmisson(msg.Text, remind_cmd, say_cmd) {
		log.Println("load default script after submisson")
		this.loadScript(default_script)
	} else if isCmd(msg.Text, chat_cmd, said_cmd, say_cmd, remind_cmd, reminders_cmd) {
		words := strings.Fields(msg.Text)
		scriptName := strings.SplitAfter(words[0], "/")[1]
		log.Println("load command script", scriptName)
		this.loadScript(scriptName)
	}
	return this
}

func (this *QProcess) findNextQuestion(msg telebot.Message) *QProcess {
	this.next = nil

	if isCmd(msg.Text, chat_cmd, said_cmd, say_cmd, remind_cmd, reminders_cmd) {
		//start at the beginning - covers submissons also
		this.next = this.lang.questions[0]
		return this
	} else {
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
					return this
				}
			}
		}
	}
	log.Println("looks like we are all done!")
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
