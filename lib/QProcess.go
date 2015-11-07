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

	chat_cmd, say_cmd, review_cmd, remind_cmd = "/chat", "/say", "/review", "/remind"

	remind_cmd_hint, review_cmd_hint, say_cmd_hint = "/remind in <days> to <msg>", "/review <days>", "/say [what you want to say]"

	default_script = "default"
	stickers_chat  = "chat"
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

		this.s.User = msg.Sender
		if msg.Sticker.Exists() {
			log.Println("incoming sticker", msg.Sticker.FileID)
			if s := this.sLib.FindSticker(msg.Sticker.FileID); s != nil {
				log.Println("We know what to do with this sticker ...", s.Meaning, s.SaveTag)
				this.loadScript(stickers_chat)
				this.findNextStickerQ(s, msg).andChat()
			}
			return
		}
		this.
			quickWinFirst(msg).
			determineScript(msg).
			findNextQuestion(msg).
			andChat()

	}
}

func (this *QProcess) quickWinFirst(msg telebot.Message) *QProcess {

	if isCmd(msg.Text, start_cmd, help_cmd) {
		appInfo := fmt.Sprintf("%s %s %s %s ",
			this.info.App,
			chat_cmd+" - to have a *quick chat* about what your upto \n\n",
			say_cmd_hint+" - to say *anything* thats on your mind \n\n",
			review_cmd_hint+" - to review what has been happening \n\n",
		)
		this.s.send(appInfo)
	} else if hasSubmisson(msg.Text, say_cmd) {
		log.Println("save something said ", msg.Text)
		this.s.save(NewNote(msg, say_cmd))
	} else if hasSubmisson(msg.Text, review_cmd) || isCmd(msg.Text, review_cmd) {
		log.Println("doing review ", msg.Text)
		n := this.s.getNotes(msg)
		this.s.send(this.info.Reminders...)
		this.s.send(n.FilterBy(said_tag).ToString())
		this.s.send(n.FilterBy(chat_tag).ToString())
	} else if hasSubmisson(msg.Text, remind_cmd) {
		log.Println("save reminder ", msg.Text)
		this.s.save(NewReminderNote(msg))
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
	} else if isCmd(msg.Text, chat_cmd, say_cmd, remind_cmd) {
		words := strings.Fields(msg.Text)
		scriptName := strings.SplitAfter(words[0], "/")[1]
		log.Println("load command script", scriptName)
		this.loadScript(scriptName)
	}
	return this
}

func (this *QProcess) findNextQuestion(msg telebot.Message) *QProcess {
	this.next = nil

	if isCmd(msg.Text, chat_cmd, say_cmd, remind_cmd) {
		//start at the beginning - covers submissons also
		this.next = this.lang.Questions[0]
		return this
	} else {
		//find the next question
		if nxt, sv := this.lang.Questions.next(msg.Text); sv {
			this.s.save(NewNote(msg, chat_cmd, this.next.RelatesTo.SaveTag))
			log.Println("After save")
			this.next = nxt
			log.Println("After setting next ", nxt.QuestionText)
			return this
		} else {
			this.next = nxt
			return this
		}
	}
}

func (this *QProcess) findNextStickerQ(s *Sticker, msg telebot.Message) *QProcess {
	this.next = nil

	if nxt, sv := this.lang.Questions.nextFrom(s.Ids...); sv {
		log.Println("About to save save")
		this.s.save(s.ToNote(msg, chat_tag))
		log.Println("After save")
		this.next = nxt
		log.Println("After setting next ", nxt.QuestionText)
		return this
	} else {
		this.next = nxt
		return this
	}
}

func (this *QProcess) andChat() {

	log.Println("are we chatting? ", this.next != nil)
	if this.next != nil {
		log.Println("chatting sending context")
		this.s.send(this.next.Context...)
		log.Println("chatting sending question")
		this.s.sendWithKeyboard(this.next.QuestionText, this.next.makeKeyboard())
		log.Println("all good")
		//return
	}

	/*if len(this.lastTime) > 0 {
		lastTimeTxt := strings.Join(this.info.Chat, "\n\n")
		for i := range this.lastTime {
			lastTimeTxt = lastTimeTxt + "\n\n" + this.lastTime[i].ToString()
		}
		//this.s.send(lastTimeTxt)
		log.Println("Last Time ", lastTimeTxt)
		this.lastTime = nil
	}*/
	return
}
