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

	chat_action = "/chat"

	review_action      = "/review"
	review_action_hint = "/review <days>"

	remind_action      = "/remind"
	remind_action_hint = "/remind in <days> to <message>"
)

func NewAction(in *Incoming, s *session, actionName string) Action {

	if s != nil {
		s.User = in.msg.Sender
	}

	cmd := actionName

	if in.isCmd() {
		cmd = in.getCmd()
	}

	if cmd == say_action {
		return SayAction{in: in, s: s}
	} else if cmd == remind_action {
		return &RemindAction{in: in, s: s}
	} else if cmd == review_action {
		return &ReviewAction{in: in, s: s}
	} else if cmd == chat_action {
		return &ChatAction{in: in, s: s}
	} else if in.isSticker() {
		return &StickerChatAction{in: in, s: s}
	}
	return &HelpAction{in: in, s: s}
}

func load(name string, q interface{}) {

	if strings.Contains(name, "/") {
		name = strings.Split(name, "/")[1]
	}

	absPath, _ := filepath.Abs(fmt.Sprintf("./config/%s.json", name))

	log.Println("path ", absPath)

	file, err := os.Open(absPath)
	if err != nil {
		log.Panic("could not load QandA language file ", err.Error())
	}

	err = json.NewDecoder(file).Decode(&q)
	if err != nil {
		log.Panic("could not decode QandA ", err.Error())
	}

	return
}

type SayAction struct {
	s  *session
	in *Incoming
}

func (a SayAction) getName() string {
	return say_action
}
func (a SayAction) getHint() string {
	return say_action_hint
}
func (a SayAction) firstUp() Action {
	a.s.save(a.in.getNote())
	return a
}

func (a SayAction) nextQuestion() *Question {

	q := a.getQuestions()

	if a.in.isCmd() {
		return q.First()
	}

	next, save := q.next(a.in.msg.Text)
	if save {
		a.s.save(a.in.getNote(a.getName(), next.RelatesTo.SaveTag))
	}
	return next
}

func (a SayAction) askQuestion() {
	q := a.nextQuestion()
	a.s.send(q.Context...)
	a.s.sendWithKeyboard(q.QuestionText, q.makeKeyboard())
}

func (a SayAction) getQuestions() Questions {

	var q struct {
		Questions `json:"QandA"`
	}

	if a.in.hasSubmisson() {
		load(default_script, &q)
		return q.Questions
	}
	load(a.getName(), &q)
	return q.Questions
}

type HelpAction struct {
	s  *session
	in *Incoming
	q  struct {
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
	helpInfo := fmt.Sprintf("%s %s %s %s ",
		fmt.Sprintf("Hey %s! We can't doFirst it all but we can:\n\n", a.in.sender().Username),
		chat_action+" - to have a *quick chat* about what your up-to \n\n",
		say_action_hint+" - to say *anything* thats on your mind \n\n",
		review_action_hint+" - to review what has been happening \n\n",
	)
	a.s.send(helpInfo)
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
	s  *session
	in *Incoming
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

	if a.in.isCmd() {
		return q.First()
	}
	next, save := q.next(a.in.msg.Text)
	if save {
		a.s.save(a.in.getNote(a.getName(), next.RelatesTo.SaveTag))
	}
	return next
}

func (a ChatAction) askQuestion() {
	q := a.nextQuestion()
	a.s.send(q.Context...)
	a.s.sendWithKeyboard(q.QuestionText, q.makeKeyboard())
	return
}

type ReviewAction struct {
	s  *session
	in *Incoming
	q  struct {
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
	n := a.s.getNotes(a.in.msg)

	saidTxt := fmt.Sprintf("%s you said: \n\n %s", a.in.sender().Username, n.FilterBy(said_tag).ToString())
	a.s.send(saidTxt)
	talkedTxt := fmt.Sprintf("%s we talked about: \n\n %s", a.in.sender().Username, n.FilterBy(chat_tag).ToString())
	a.s.send(talkedTxt)
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

	if a.in.isCmd() {
		return q.First()
	}
	next, save := q.next(a.in.msg.Text)
	if save {
		a.s.save(a.in.getNote(a.getName(), next.RelatesTo.SaveTag))
	}
	return next
}

func (a ReviewAction) askQuestion() {
	q := a.nextQuestion()
	a.s.send(q.Context...)
	a.s.sendWithKeyboard(q.QuestionText, q.makeKeyboard())
	return
}

type RemindAction struct {
	s  *session
	in *Incoming
}

func (a RemindAction) getName() string {
	return remind_action
}
func (a RemindAction) getHint() string {
	return remind_action_hint
}
func (a RemindAction) firstUp() Action {
	log.Println("remind me ...")
	a.s.save(a.in.getNote())
	return a
}

func (a RemindAction) getQuestions() Questions {

	var q struct {
		Questions `json:"QandA"`
	}

	if a.in.hasSubmisson() {
		load(default_script, &q)
		return q.Questions
	}
	load(a.getName(), &q)

	return q.Questions
}

func (a RemindAction) nextQuestion() *Question {

	q := a.getQuestions()

	if a.in.isCmd() {
		return q.First()
	}
	next, save := q.next(a.in.msg.Text)
	if save {
		a.s.save(a.in.getNote(a.getName(), next.RelatesTo.SaveTag))
	}
	return next
}

func (a RemindAction) askQuestion() {
	q := a.nextQuestion()
	a.s.send(q.Context...)
	a.s.sendWithKeyboard(q.QuestionText, q.makeKeyboard())
	return
}

type StickerChatAction struct {
	s        *session
	in       *Incoming
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

	if a.in.isSticker() {

		sticker := a.stickers.FindSticker(a.in.msg.Sticker.FileID)
		a.in.msg.Text = sticker.Meaning
		next, save := q.nextFrom(sticker.Ids...)
		if save {
			a.s.save(a.in.getNote(a.getName(), next.RelatesTo.SaveTag))
		}
		return next
	}
	next, save := q.next(a.in.msg.Text)
	if save {
		a.s.save(a.in.getNote(a.getName(), next.RelatesTo.SaveTag))
	}
	return next
}

func (a StickerChatAction) askQuestion() {
	q := a.nextQuestion()
	a.s.send(q.Context...)
	a.s.sendWithKeyboard(q.QuestionText, q.makeKeyboard())
	return
}
