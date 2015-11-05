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
	stickers_chat  = "stickers_chat"
)

type (
	Info struct {
		App       []string `json:"appInfo"`
		Reminders []string `json:"remindersInfo"`
		Chat      []string `json:"chatInfo"`
		Said      []string `json:"saidInfo"`
	}

	QProcess struct {
		s    *session
		lang struct {
			questions `json:"QandA"`
		}
		info     *Info
		next     *Question
		lastTime Notes
		sLib     Stickers
	}
)

func NewQProcess(b *telebot.Bot, s *Storage) *QProcess {
	q := &QProcess{s: newSession(b, s), sLib: LoadKnownStickers()}
	q.loadInfo()
	return q
}

func (this *QProcess) Run(input <-chan telebot.Message) {
	for msg := range input {

		this.s.User = msg.Sender
		if msg.Sticker.Exists() {
			log.Println("incoming sticker", msg.Sticker.FileID)

			s := this.sLib.FindSticker(msg.Sticker.FileID)
			if s != nil {
				log.Println("We know what to do with this sticker ...", s.Meaning, s.SaveTag)
				this.loadScript(stickers_chat)
				//get going
				this.findNextStickerQ(s, msg).andChat()
			}
		} else {
			this.
				quickWinFirst(msg).
				determineScript(msg).
				findNextQuestion(msg).
				andChat()
		}
	}
}

func (this *QProcess) quickWinFirst(msg telebot.Message) *QProcess {

	if isCmd(msg.Text, start_cmd, help_cmd) {

		appInfo := fmt.Sprintf("%s %s %s %s %s %s",
			this.info.App,
			chat_cmd+" - to have a *quick chat* about what your upto \n\n",
			say_cmd_hint+" - to say *anything* thats on your mind \n\n",
			said_cmd+" - to show all the things you have said \n\n",
			remind_cmd_hint+" - to keep track of the things you need to be reminded about \n\n",
			reminders_cmd+" - to show all the reminders you have made",
		)

		this.s.send(appInfo)
	} else if hasSubmisson(msg.Text, say_cmd) {
		log.Println("save something said ", msg.Text)
		this.s.save(NewNote(msg, say_cmd))
	} else if hasSubmisson(msg.Text, remind_cmd) {
		log.Println("save a reminder ", msg.Text)
		this.s.save(NewReminderNote(msg))
	} else if hasSubmisson(msg.Text, reminders_cmd) || isCmd(msg.Text, reminders_cmd) {
		log.Println("showing reminders ", msg.Text)
		r := this.s.getReminders(msg)
		this.s.send(this.info.Reminders...)
		for i := range r {
			this.s.send(r[i].ToString())

		}
	} else if hasSubmisson(msg.Text, said_cmd) || isCmd(msg.Text, said_cmd) {
		log.Println("showing things said ", msg.Text)
		r := this.s.getNotes(msg)
		this.s.send(this.info.Said...)
		for i := range r {
			this.s.send(r[i].ToString())
		}
	}
	return this
}

func (this *QProcess) loadInfo() {
	file, err := os.Open("./config/fantail.json")
	if err != nil {
		log.Panic("could not load App info", err.Error())
	}
	err = json.NewDecoder(file).Decode(&this.info)
	if err != nil {
		log.Panic("could not decode App info ", err.Error())
	}
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
			for a := range this.lang.questions[i].RelatesTo.Answers {
				if this.lang.questions[i].RelatesTo.Answers[a] == msg.Text {
					//was the answer a remainder to save?
					if this.lang.questions[i].RelatesTo.Save {
						this.s.save(NewNote(msg, chat_cmd, this.lang.questions[i].RelatesTo.SaveTag))
						this.lastTime = append(this.lastTime, this.s.getLastChat(this.lang.questions[i].RelatesTo.SaveTag))
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

func (this *QProcess) findNextStickerQ(s *Sticker, msg telebot.Message) *QProcess {
	this.next = nil

	for i := range this.lang.questions {
		for a := range this.lang.questions[i].RelatesTo.Answers {
			for si := range s.Ids {
				if this.lang.questions[i].RelatesTo.Answers[a] == s.Ids[si] {
					if this.lang.questions[i].RelatesTo.Save {
						this.s.save(s.ToNote(msg))
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

func (this *QProcess) andChat() {
	if this.next != nil {
		this.s.send(this.next.Context...)
		this.s.sendWithKeyboard(this.next.QuestionText, this.makeKeyboard())
		return
	}

	if len(this.lastTime) > 0 {
		this.s.send(this.info.Chat...)
		for i := range this.lastTime {
			this.s.send(this.lastTime[i].ToString())
		}
		this.lastTime = nil
	}
	return
}
