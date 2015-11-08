package lib

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
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
			Questions `json:"QandA"`
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
		this.s.addDetails(msg)

		if msg.Sticker.Exists() {
			log.Println("incoming sticker", msg.Sticker.FileID)
			if s := this.sLib.FindSticker(msg.Sticker.FileID); s != nil {
				this.
					determineScript(msg).
					nextStickerQ(s, msg).
					andChat()
			}
		} else {

			this.
				quickWinFirst(msg).
				determineScript(msg).
				nextQ(msg).
				andChat()
		}

	}
}

func (this *QProcess) quickWinFirst(msg telebot.Message) *QProcess {

	if this.s.Action.isHelp() {

		this.s.send(fmt.Sprintf("%s %s", this.info.App, this.s.Action.typeOf))

	} else if this.s.Action.setTypeMatches(say_action) && this.s.Action.hasSubmisson() {

		this.s.save(NewNote(msg, this.s.Action.typeOf))

	} else if this.s.Action.setTypeMatches(remind_action) && this.s.Action.hasSubmisson() {

		this.s.save(NewReminderNote(msg))

	} else if this.s.Action.setTypeMatches(recap_action) {

		n := this.s.getNotes(msg)
		this.s.send(this.info.Reminders...)
		this.s.send(n.FilterBy(said_tag).ToString())
		this.s.send(n.FilterBy(chat_tag).ToString())

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

func (this *QProcess) determineScript(msg telebot.Message) *QProcess {
	this.s.Action.loadQuestions(this.lang.Questions)
	return this
}

func (this *QProcess) nextQ(msg telebot.Message) *QProcess {
	this.next = nil

	if this.s.Action.typeIsSet() {
		this.next = this.lang.Questions.First()
		return this
	} else {
		if nxt, sv := this.lang.Questions.next(msg.Text); sv {
			this.s.save(NewNote(msg, this.s.Action.typeOf, this.next.RelatesTo.SaveTag))
			this.next = nxt
			return this
		} else {
			this.next = nxt
			return this
		}
	}
}

func (this *QProcess) nextStickerQ(s *Sticker, msg telebot.Message) *QProcess {
	this.next = nil

	if nxt, sv := this.lang.Questions.nextFrom(s.Ids...); sv {
		this.s.save(s.ToNote(msg, this.s.Action.typeOf))
		this.next = nxt
		return this
	} else {
		this.next = nxt
		return this
	}
}

func (this *QProcess) andChat() {
	if this.next != nil {
		this.s.send(this.next.Context...)
		this.s.sendWithKeyboard(this.next.QuestionText, this.next.makeKeyboard())
	}
	return
}
