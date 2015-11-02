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
	//being a good citizen
	start_cmd, help_cmd = "/start", "/help"

	chat_cmd, remind_cmd, say_cmd, reminders_cmd, said_cmd = "/chat", "/remind", "/say", "/reminders", "/said"

	remind_cmd_hint, say_cmd_hint = "/remind me in <days> to <message>", "/say [what you want to say]"

	default_script = "default"
)

type QProcess struct {
	s    *session
	lang struct {
		questions `json:"QandA"`
	}
	next *Question
}

func NewQProcess(b *telebot.Bot, s *Storage) *QProcess {
	return &QProcess{s: newSession(b, s)}
}

func (this *QProcess) Run(input <-chan telebot.Message) {
	for msg := range input {
		this.s.User = msg.Chat
		this.quickWinFirst(msg).
			determineScript(msg).
			findNextQuestion(msg).
			andAsk()
	}
}

func (this *QProcess) checkForDefaults(msg telebot.Message) *QProcess {

	if isCmd(msg.Text, start_cmd, help_cmd) {
		this.s.send(
			"Fantail is your companion that is here to help you get the help you want quicker",
			"You can control me by sending these commands:",
			chat_cmd,
			say_cmd_hint,
			said_cmd,
			remind_cmd_hint,
			reminders_cmd,
		)
	}
	return this
}

func (this *QProcess) quickWinFirst(msg telebot.Message) *QProcess {

	if hasSubmisson(msg.Text, say_cmd) {
		log.Println("save something said ", msg.Text)
		this.s.save(msg, say_cmd)
	} else if hasSubmisson(msg.Text, remind_cmd) {
		log.Println("save a reminder ", msg.Text)
		this.s.saveAsReminder(msg)
	} else if hasSubmisson(msg.Text, reminders_cmd) || isCmd(msg.Text, reminders_cmd) {
		log.Println("showing reminders ", msg.Text)
		r := this.s.getReminders(msg)
		for i := range r {
			this.s.send(r[i].ToString())

		}
	} else if hasSubmisson(msg.Text, said_cmd) || isCmd(msg.Text, said_cmd) {
		log.Println("showing things said ", msg.Text)
		r := this.s.getNotes(msg)
		for i := range r {
			this.s.send(r[i].ToString())
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
						this.s.save(msg, chat_cmd, this.lang.questions[i].RelatesTo.SaveTag)
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
		this.s.send(this.next.Context...)
		this.s.sendWithKeyboard(this.next.QuestionText, this.makeKeyboard())
	}

	return
}
