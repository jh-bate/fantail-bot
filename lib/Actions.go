package lib

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Action interface {
	getName() string
	getHint() string
	firstUp() Action
	getQuestions() Questions
	nextQuestion() *Question
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

func NewAction(s *session, actionName string) Action {

	if s != nil {
		s.User = s.getSender()
	}

	cmd := actionName

	if s.sentAsCommand() {
		cmd = s.getSentCommand()
	}

	if cmd == say_action {
		return SayAction{session: s}
	} else if cmd == ask_action {
		return AskAction{session: s}
	} else if cmd == review_action {
		return &ReviewAction{session: s}
	} else if cmd == chat_action {
		return &ChatAction{session: s}
	} else if s.sentAsSticker() {
		return &StickerChatAction{session: s}
	}
	return &HelpAction{session: s}
}

func load(name string, q interface{}) {

	if strings.Contains(name, "/") {
		name = strings.Split(name, "/")[1]
	}

	absPath, _ := filepath.Abs(fmt.Sprintf("config/%s.json", name))

	log.Println("QandA", absPath)

	file, err := os.Open(absPath)
	if err != nil {
		log.Println("could not load QandA language file", err.Error())
		absPath, _ = filepath.Abs(fmt.Sprintf("lib/config/%s.json", name))
		log.Println("QandA path ", absPath)

		file, err = os.Open(absPath)
	}

	err = json.NewDecoder(file).Decode(&q)
	if err != nil {
		log.Panic("could not decode QandA ", err.Error())
	}

	return
}

type SayAction struct {
	*session
}

func (a SayAction) getName() string {
	return say_action
}
func (a SayAction) getHint() string {
	return say_action_hint
}
func (a SayAction) firstUp() Action {
	a.session.save(a.getName())
	return a
}

func (a SayAction) nextQuestion() *Question {

	q := a.getQuestions()

	if a.session.sentAsCommand() {
		return q.First()
	}

	next, save := q.next(a.getSentMsgText())
	if save {
		a.session.save(a.getName(), next.RelatesTo.SaveTag)
	}
	return next
}

func (a SayAction) askQuestion() {
	if q := a.nextQuestion(); q != nil {
		a.session.send(q.Context...)
		a.session.sendWithKeyboard(q.QuestionText, q.makeKeyboard())
	}
}

func (a SayAction) getQuestions() Questions {

	var q struct {
		Questions `json:"QandA"`
	}
	load(a.getName(), &q)
	return q.Questions
}

type AskAction struct {
	*session
}

func (a AskAction) getName() string {
	return ask_action
}
func (a AskAction) getHint() string {
	return ask_action_hint
}
func (a AskAction) firstUp() Action {
	a.session.save(a.getName(), help_tag)
	return a
}

func (a AskAction) nextQuestion() *Question {

	q := a.getQuestions()

	if a.session.sentAsCommand() {
		return q.First()
	}

	next, save := q.next(a.getSentMsgText())
	if save {
		a.session.save(a.getName(), next.RelatesTo.SaveTag)
	}
	return next
}

func (a AskAction) askQuestion() {
	if q := a.nextQuestion(); q != nil {
		a.session.send(q.Context...)
		a.session.sendWithKeyboard(q.QuestionText, q.makeKeyboard())
	}
}

func (a AskAction) getQuestions() Questions {

	var q struct {
		Questions `json:"QandA"`
	}

	if a.session.sentAsSubmission() {
		load(default_script, &q)
		return q.Questions
	}
	load(a.getName(), &q)
	return q.Questions
}

type HelpAction struct {
	*session
	q struct {
		Questions `json:"QandA"`
	}
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
		fmt.Sprintf("Hey %s! We can't do it all but we can:\n\n", a.session.getSentUsername()),
		chat_action+" - to have a *quick chat* about what your up-to \n\n",
		say_action_hint+" - to say *anything* thats on your mind \n\n",
		ask_action_hint+" - to ask *anything* thats on your mind \n\n",
		review_action_hint+" - to review what has been happening \n\n",
		"Stickers - we have those to help express yourself!! [Get them here](https://telegram.me/addstickers/betes) \n\n",
	)
	a.session.send(helpInfo)
	return a
}

func (a HelpAction) nextQuestion() *Question {
	return nil
}

func (a HelpAction) getQuestions() Questions {
	return nil
}

func (a HelpAction) askQuestion() {
	return
}

type ChatAction struct {
	*session
}

func (a ChatAction) getName() string {
	return chat_action
}
func (a ChatAction) getHint() string {
	return ""
}
func (a ChatAction) firstUp() Action {
	//nothing to doFirst
	return a
}

func (a ChatAction) getQuestions() Questions {

	var q struct {
		Questions `json:"QandA"`
	}

	load(a.getName(), &q)
	return q.Questions
}

func (a ChatAction) nextQuestion() *Question {

	q := a.getQuestions()

	if a.session.sentAsCommand() {
		return q.First()
	}
	next, save := q.next(a.session.getSentMsgText())
	if save {
		a.session.save(a.getName(), next.RelatesTo.SaveTag)
	}
	return next
}

func (a ChatAction) askQuestion() {
	if q := a.nextQuestion(); q != nil {
		a.session.send(q.Context...)
		a.session.sendWithKeyboard(q.QuestionText, q.makeKeyboard())
	}
	return
}

type ReviewAction struct {
	*session
	q struct {
		Questions `json:"QandA"`
	}
}

func (a ReviewAction) getName() string {
	return review_action
}
func (a ReviewAction) getHint() string {
	return review_action_hint
}
func (a ReviewAction) firstUp() Action {
	log.Println("doFirsting review ...")
	n := a.session.getNotes()

	saidTxt := fmt.Sprintf("%s you said: \n\n %s", a.session.getSentUsername(), n.FilterBy(said_tag).ToString())
	a.session.send(saidTxt)
	talkedTxt := fmt.Sprintf("%s we talked about: \n\n %s", a.session.getSentUsername(), n.FilterBy(chat_tag).ToString())
	a.session.send(talkedTxt)
	return a
}

func (a ReviewAction) getQuestions() Questions {

	var q struct {
		Questions `json:"QandA"`
	}

	load(a.getName(), &q)

	return q.Questions
}

func (a ReviewAction) nextQuestion() *Question {

	q := a.getQuestions()

	if a.session.sentAsCommand() {
		return q.First()
	}
	next, save := q.next(a.session.getSentMsgText())
	if save {
		a.session.save(a.getName(), next.RelatesTo.SaveTag)
	}
	return next
}

func (a ReviewAction) askQuestion() {
	if q := a.nextQuestion(); q != nil {
		a.session.send(q.Context...)
		a.session.sendWithKeyboard(q.QuestionText, q.makeKeyboard())
	}
	return
}

type StickerChatAction struct {
	*session
	stickers Stickers
}

func (a StickerChatAction) getName() string {
	return chat_action
}
func (a StickerChatAction) getHint() string {
	return ""
}

func (a StickerChatAction) firstUp() Action {
	a.stickers = LoadKnownStickers()
	return a
}

func (a StickerChatAction) getQuestions() Questions {
	var q struct {
		Questions `json:"QandA"`
	}
	load(a.getName(), &q)
	return q.Questions
}

func (a StickerChatAction) nextQuestion() *Question {

	q := a.getQuestions()

	if a.session.sentAsSticker() {

		sticker := a.stickers.FindSticker(a.session.getSentStickerId())
		a.session.setSentMsgText(sticker.Meaning)
		next, save := q.nextFrom(sticker.Ids...)
		if save {
			a.session.save(a.getName(), next.RelatesTo.SaveTag)
		}
		return next
	}
	next, save := q.next(a.session.getSentMsgText())
	if save {
		a.session.save(a.getName(), next.RelatesTo.SaveTag)
	}
	return next
}

func (a StickerChatAction) askQuestion() {
	if q := a.nextQuestion(); q != nil {
		a.session.send(q.Context...)
		a.session.sendWithKeyboard(q.QuestionText, q.makeKeyboard())
	}
	return
}
