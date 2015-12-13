package incoming

import (
	"fmt"
	"log"

	"github.com/jh-bate/fantail-bot/config"
	"github.com/jh-bate/fantail-bot/note"
	"github.com/jh-bate/fantail-bot/question"
	"github.com/jh-bate/fantail-bot/sticker"
)

type Action interface {
	getName() string
	getHint() string
	firstUp() Action
	getQuestions() question.Questions
	nextQuestion() *question.Question
	askQuestion()
}

const (
	start_action, help_action = "/start", "/help"

	say_action      = "/say"
	say_action_hint = "/say <message>"

	ask_action      = "/ask"
	ask_action_hint = "/ask <message>"

	chat_action = "/chat"

	review_action      = "/review"
	review_action_hint = "/review <days>"

	default_script = "default"
)

var Config struct {
	question.Questions `json:"QandA"`
}

func NewAction(msg *Message, runas string) Action {

	msgAction := ""

	if msg.HasAction() {
		msgAction = msg.Action
	} else {
		log.Println("run as", runas)
		msgAction = runas
	}

	if msgAction == say_action {
		return SayAction{Message: msg}
	} else if msgAction == ask_action {
		return AskAction{Message: msg}
	} else if msgAction == review_action {
		return &ReviewAction{Message: msg}
	} else if msgAction == chat_action {
		return &ChatAction{Message: msg}
	} else if msg.Sticker.Has {
		return &StickerChatAction{Message: msg}
	}
	return &HelpAction{Message: msg}
}

type SayAction struct {
	*Message
	*Session
}

func (a SayAction) getName() string {
	return say_action
}
func (a SayAction) getHint() string {
	return say_action_hint
}
func (a SayAction) firstUp() Action {
	if a.HasSubmisson {
		n := note.New(a.Text, a.Sender.ID, a.Date)
		n.Save()
	}
	return a
}

func (a SayAction) nextQuestion() *question.Question {

	q := a.getQuestions()

	if a.HasAction() {
		return q.First()
	}

	next, save := q.Next(a.Text)
	if save {
		n := note.New(a.Text, a.Sender.ID, a.Date, a.getName(), next.RelatesTo.SaveTag)
		n.Save()
	}
	return next
}

func (a SayAction) askQuestion() {
	if q := a.nextQuestion(); q != nil {
		a.send(a.Sender, q.Context...)
		a.sendWithKeyboard(a.Sender, q.QuestionText, q.MakeKeyboard())
	}
}

func (a SayAction) getQuestions() question.Questions {
	config.Load(&Config, "say.json")
	return Config.Questions
}

type AskAction struct {
	*Message
	*Session
}

func (a AskAction) getName() string {
	return ask_action
}
func (a AskAction) getHint() string {
	return ask_action_hint
}
func (a AskAction) firstUp() Action {
	if a.HasSubmisson {
		n := note.New(a.Text, a.Sender.ID, a.Date, note.HELP_TAG)
		n.Save()
	}
	return a
}

func (a AskAction) nextQuestion() *question.Question {

	q := a.getQuestions()

	if a.HasAction() {
		return q.First()
	}

	next, save := q.Next(a.Text)
	if save {
		n := note.New(a.Text, a.Sender.ID, a.Date, a.getName(), next.RelatesTo.SaveTag)
		n.Save()
	}
	return next
}

func (a AskAction) askQuestion() {
	if q := a.nextQuestion(); q != nil {
		a.send(a.Sender, q.Context...)
		a.sendWithKeyboard(a.Sender, q.QuestionText, q.MakeKeyboard())
	}
}

func (a AskAction) getQuestions() question.Questions {
	config.Load(&Config, "ask.json")
	return Config.Questions
}

type HelpAction struct {
	*Message
	*Session
}

func (a HelpAction) getName() string {
	return say_action
}
func (a HelpAction) getHint() string {
	return say_action_hint
}
func (a HelpAction) firstUp() Action {
	log.Println("help first up")
	helpInfo := fmt.Sprintf("%s %s %s %s %s %s ",
		fmt.Sprintf("Hey %s! We can't do it all but we can:\n\n", a.Sender.FirstName),
		chat_action+" - to have a *quick chat* about what your up-to \n\n",
		say_action_hint+" - to say *anything* thats on your mind \n\n",
		ask_action_hint+" - to ask *anything* thats on your mind \n\n",
		review_action_hint+" - to review what has been happening \n\n",
		"Stickers - we have those to help express yourself!! \n\n [Get them here](https://telegram.me/addstickers/betes)",
	)
	a.send(a.Sender, helpInfo)
	return a
}

func (a HelpAction) nextQuestion() *question.Question {
	return nil
}

func (a HelpAction) getQuestions() question.Questions {
	return nil
}

func (a HelpAction) askQuestion() {
	return
}

type ChatAction struct {
	*Session
	*Message
}

func (a ChatAction) getName() string {
	return chat_action
}
func (a ChatAction) getHint() string {
	return ""
}
func (a ChatAction) firstUp() Action {
	return a
}

func (a ChatAction) getQuestions() question.Questions {
	config.Load(&Config, "chat.json")
	return Config.Questions
}

func (a ChatAction) nextQuestion() *question.Question {

	q := a.getQuestions()

	if a.HasAction() {
		return q.First()
	}
	next, save := q.Next(a.Text)
	if save {
		n := note.New(
			a.Text,
			a.Sender.ID,
			a.Date,
			a.getName(),
			next.RelatesTo.SaveTag,
		)
		n.Save()
	}
	return next
}

func (a ChatAction) askQuestion() {
	if q := a.nextQuestion(); q != nil {
		a.send(a.Sender, q.Context...)
		a.sendWithKeyboard(a.Sender, q.QuestionText, q.MakeKeyboard())
	}
	return
}

type ReviewAction struct {
	*Message
	*Session
}

func (a ReviewAction) getName() string {
	return review_action
}
func (a ReviewAction) getHint() string {
	return review_action_hint
}
func (a ReviewAction) firstUp() Action {
	log.Println("doFirsting review ...")

	notes, _ := note.GetNotes(fmt.Sprintf("%d", a.Sender.ID))

	saidTxt := fmt.Sprintf("%s you said: \n\n %s", a.Sender.FirstName, notes.FilterOnTag(note.SAID_TAG).ToString())
	a.send(a.Sender, saidTxt)
	talkedTxt := fmt.Sprintf("%s we talked about: \n\n %s", a.Sender.FirstName, notes.FilterOnTag(note.CHAT_TAG).ToString())
	a.send(a.Sender, talkedTxt)
	return a
}

func (a ReviewAction) getQuestions() question.Questions {
	config.Load(&Config, "review.json")
	return Config.Questions
}

func (a ReviewAction) nextQuestion() *question.Question {

	q := a.getQuestions()

	if a.HasAction() {
		return q.First()
	}
	next, save := q.Next(a.Text)
	if save {
		n := note.New(a.Text, a.Sender.ID, a.Date, a.getName(), next.RelatesTo.SaveTag)
		n.Save()
	}
	return next
}

func (a ReviewAction) askQuestion() {
	if q := a.nextQuestion(); q != nil {
		a.send(a.Sender, q.Context...)
		a.sendWithKeyboard(a.Sender, q.QuestionText, q.MakeKeyboard())
	}
	return
}

type StickerChatAction struct {
	*Session
	sticker.Stickers
	*Message
}

func (a StickerChatAction) getName() string {
	return chat_action
}
func (a StickerChatAction) getHint() string {
	return ""
}

func (a StickerChatAction) firstUp() Action {
	a.Stickers = sticker.Load()
	return a
}

func (a StickerChatAction) getQuestions() question.Questions {
	config.Load(&Config, "chat.json")
	return Config.Questions
}

func (a StickerChatAction) nextQuestion() *question.Question {

	q := a.getQuestions()

	if a.Sticker.Has {

		sticker := a.Stickers.Find(a.Sticker.Id)

		next, save := q.NextFrom(sticker.Ids...)
		if save {
			n := note.New(sticker.Meaning, a.Sender.ID, a.Date, a.getName(), next.RelatesTo.SaveTag)
			n.Save()
		}
		return next
	}
	next, save := q.Next(a.Text)
	if save {
		n := note.New(a.Text, a.Sender.ID, a.Date, a.getName(), next.RelatesTo.SaveTag)
		n.Save()
	}
	return next
}

func (a StickerChatAction) askQuestion() {
	if q := a.nextQuestion(); q != nil {
		a.send(a.Sender, q.Context...)
		a.sendWithKeyboard(a.Sender, q.QuestionText, q.MakeKeyboard())
	}
	return
}
