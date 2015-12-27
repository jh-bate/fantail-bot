package bot

import (
	"fmt"
	"log"
	"strings"

	"github.com/jh-bate/fantail-bot/config"
	"github.com/jh-bate/fantail-bot/note"
	"github.com/jh-bate/fantail-bot/question"
)

type Action interface {
	Name() string
	Payload() *Payload
	Run()
}

const (
	start_action   = "/start"
	default_script = "default"
)

func NewAction(payload *Payload, runas string, session *Session) Action {

	msgAction := ""

	if payload.HasAction() {
		msgAction = payload.Action
	} else {
		log.Println("running as ...", runas)
		msgAction = runas
	}

	if msgAction == say_action {
		return &SayAction{payload: payload, Session: session}
	} else if msgAction == ask_action {
		return &AskAction{payload: payload, Session: session}
	} else if msgAction == review_action {
		return &ReviewAction{payload: payload, Session: session}
	} else if msgAction == chat_action {
		return &ChatAction{payload: payload, Session: session}
	}
	return &HelpAction{payload: payload, Session: session}
}

func nextQuestion(a Action) *question.Question {

	var Config struct {
		question.Questions `json:"QandA"`
	}

	payload := a.Payload()

	config.Load(&Config, strings.Split(a.Name(), "/")[1]+".json")

	if payload.HasAction() {
		return Config.Questions.First()
	}

	next, save := Config.Questions.Next(payload.Text)
	if save {
		n := note.New(payload.Text, payload.Sender, payload.Date, a.Name(), next.RelatesTo.SaveTag)
		n.Save()
	}
	return next
}

const (
	say_action      = "/say"
	say_action_hint = "/say <message>"
)

type SayAction struct {
	payload *Payload
	Session *Session
}

func (a SayAction) Name() string {
	return say_action
}

func (a SayAction) Payload() *Payload {
	return a.payload
}

func (a SayAction) Run() {
	if a.payload.HasSubmisson {
		n := note.New(a.payload.Text, a.payload.Sender, a.payload.Date)
		n.Save()
		return
	}
	if q := nextQuestion(a); q != nil {
		a.Session.send(a.payload.Sender, q.Context...)
		a.Session.send(a.payload.Sender, q.QuestionText)
		//a.sendWithKeyboard(a.Sender, q.QuestionText, q.MakeKeyboard())
	}
	return
}

const (
	review_action      = "/review"
	review_action_hint = "/review <days>"
)

type ReviewAction struct {
	payload *Payload
	Session *Session
}

func (a ReviewAction) Name() string {
	return review_action
}

func (a ReviewAction) Payload() *Payload {
	return a.payload
}

func (a ReviewAction) Run() {
	log.Println("do review ...")

	notes, _ := note.GetNotes(fmt.Sprintf("%d", a.payload.Sender))

	saidTxt := fmt.Sprintf("You said: \n\n %s", notes.FilterOnTag(note.SAID_TAG).ToString())
	a.Session.send(a.payload.Sender, saidTxt)
	talkedTxt := fmt.Sprintf("We talked about: \n\n %s", notes.FilterOnTag(note.CHAT_TAG).ToString())
	a.Session.send(a.payload.Sender, talkedTxt)
	return
}

const (
	help_action = "/help"
)

type HelpAction struct {
	payload *Payload
	Session *Session
}

func (a HelpAction) Name() string {
	return help_action
}

func (a HelpAction) Payload() *Payload {
	return a.payload
}

func (a HelpAction) Run() {
	log.Println("help first up")
	helpInfo := fmt.Sprintf("%s %s %s %s %s %s ",
		//fmt.Sprintf("Hey %s! We can't do it all but we can:\n\n", a.Sender.FirstName),
		"Hey! We can't do it all but we can:\n\n",
		chat_action+" - to have a *quick chat* about what your up-to \n\n",
		say_action_hint+" - to say *anything* thats on your mind \n\n",
		ask_action_hint+" - to ask *anything* thats on your mind \n\n",
		review_action_hint+" - to review what has been happening \n\n",
		"Stickers - we have those to help express yourself!! \n\n [Get them here](https://telegram.me/addstickers/betes)",
	)
	a.Session.send(a.payload.Sender, helpInfo)
	return
}

const (
	chat_action = "/chat"
)

type ChatAction struct {
	Session *Session
	payload *Payload
}

func (a ChatAction) Name() string {
	return chat_action
}

func (a ChatAction) Payload() *Payload {
	return a.payload
}

func (a ChatAction) Run() {
	if q := nextQuestion(a); q != nil {
		a.Session.send(a.payload.Sender, q.Context...)
		a.Session.send(a.payload.Sender, q.QuestionText)
		//a.sendWithKeyboard(a.Sender, q.QuestionText, q.MakeKeyboard())
	}
	return
}

const (
	ask_action      = "/ask"
	ask_action_hint = "/ask <message>"
)

type AskAction struct {
	payload *Payload
	Session *Session
}

func (a AskAction) Name() string {
	return ask_action
}

func (a AskAction) Payload() *Payload {
	return a.payload
}

func (a AskAction) Run() {
	if a.payload.HasSubmisson {
		n := note.New(a.payload.Text, a.payload.Sender, a.payload.Date, note.HELP_TAG)
		n.Save()
		return
	}
	if q := nextQuestion(a); q != nil {
		a.Session.send(a.payload.Sender, q.Context...)
		a.Session.send(a.payload.Sender, q.QuestionText)
		//a.sendWithKeyboard(a.Sender, q.QuestionText, q.MakeKeyboard())
	}
	return
}
