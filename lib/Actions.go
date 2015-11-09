package lib

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

type Action interface {
	getName() string
	getHint() string
	loadQuestions()
	doFirst()
	findNext() *Question
	chat(q *Question)
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

func NewAction(in Incoming, s *session) Action {

	if in.getCmd() == say_action {
		return &SayAction{in: &in, s: s}
	} else if in.getCmd() == remind_action {
		return &RemindAction{in: &in, s: s}
	} else if in.getCmd() == review_action {
		return &ReviewAction{in: &in, s: s}
	} else if in.getCmd() == chat_action {
		return &ChatAction{in: &in, s: s}
	} else if in.isSticker() {
		return &StickerChatAction{in: &in, s: s}
	}
	log.Println("asked ", in.getCmd())
	return &HelpAction{in: &in, s: s}
}

func load(name string, q interface{}) {

	if strings.Contains(name, "/") {
		name = strings.Split(name, "/")[1]
	}

	file, err := os.Open(fmt.Sprintf("./config/%s.json", name))
	if err != nil {
		log.Panic("could not load QandA language file ", err.Error())
	}
	err = json.NewDecoder(file).Decode(&q)
	if err != nil {
		log.Panic("could not decode QandA ", err.Error())
	}
}

type SayAction struct {
	s  *session
	in *Incoming
	q  struct {
		Questions `json:"QandA"`
	}
}

func (a SayAction) getName() string {
	return say_action
}
func (a SayAction) getHint() string {
	return say_action_hint
}
func (a SayAction) doFirst() {
	log.Println("say ...")
	a.s.save(a.in.getNote())
	return
}

func (a SayAction) findNext() *Question {
	if a.in.isCmd() {
		return a.q.Questions.First()
	}
	next, save := a.q.next(a.in.msg.Text)
	if save {
		a.s.save(a.in.getNote(a.getName(), next.RelatesTo.SaveTag))
	}
	return next
}

func (a SayAction) chat(q *Question) {
	a.s.send(q.Context...)
	a.s.sendWithKeyboard(q.QuestionText, q.makeKeyboard())
}

func (a SayAction) loadQuestions() {
	if a.in.hasSubmisson() {
		load(default_script, &a.q)
		return
	}
	load(a.getName(), &a.q)
	return
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
func (a HelpAction) doFirst() {
	helpInfo := fmt.Sprintf("%s %s %s %s ",
		fmt.Sprintf("Hey %s! We can't doFirst it all but we can:\n\n", a.in.sender().Username),
		chat_action+" - to have a *quick chat* about what your up-to \n\n",
		say_action_hint+" - to say *anything* thats on your mind \n\n",
		review_action_hint+" - to review what has been happening \n\n",
	)
	a.s.send(helpInfo)
	return
}

func (a HelpAction) loadQuestions() {
	return
}

func (a HelpAction) findNext() *Question {
	return nil
}

func (a HelpAction) chat(q *Question) {
	return
}

type ChatAction struct {
	s  *session
	in *Incoming
	q  struct {
		Questions `json:"QandA"`
	}
}

func (a ChatAction) getName() string {
	return chat_action
}
func (a ChatAction) getHint() string {
	return ""
}
func (a ChatAction) doFirst() {
	//nothing to doFirst
	return
}
func (a ChatAction) loadQuestions() {
	load(a.getName(), &a.q)
	return
}

func (a ChatAction) findNext() *Question {
	if a.in.isCmd() {
		return a.q.Questions.First()
	}
	next, save := a.q.next(a.in.msg.Text)
	if save {
		a.s.save(a.in.getNote(a.getName(), next.RelatesTo.SaveTag))
	}
	return next
}

func (a ChatAction) chat(q *Question) {
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
func (a ReviewAction) doFirst() {
	log.Println("doFirsting review ...")
	n := a.s.getNotes(a.in.msg)

	saidTxt := fmt.Sprintf("%s you said: \n\n %s", a.in.sender().Username, n.FilterBy(said_tag).ToString())
	a.s.send(saidTxt)
	talkedTxt := fmt.Sprintf("%s we talked about: \n\n %s", a.in.sender().Username, n.FilterBy(chat_tag).ToString())
	a.s.send(talkedTxt)
	return
}

func (a ReviewAction) loadQuestions() {
	load(a.getName(), &a.q)
	return
}

func (a ReviewAction) findNext() *Question {
	if a.in.isCmd() {
		return a.q.Questions.First()
	}
	next, save := a.q.next(a.in.msg.Text)
	if save {
		a.s.save(a.in.getNote(a.getName(), next.RelatesTo.SaveTag))
	}
	return next
}

func (a ReviewAction) chat(q *Question) {
	a.s.send(q.Context...)
	a.s.sendWithKeyboard(q.QuestionText, q.makeKeyboard())
	return
}

type RemindAction struct {
	s  *session
	in *Incoming
	q  struct {
		Questions `json:"QandA"`
	}
}

func (a RemindAction) getName() string {
	return remind_action
}
func (a RemindAction) getHint() string {
	return remind_action_hint
}
func (a RemindAction) doFirst() {
	log.Println("remind me ...")
	a.s.save(a.in.getNote())
	return
}

func (a RemindAction) loadQuestions() {
	if a.in.hasSubmisson() {
		load(default_script, &a.q)
		return
	}
	load(a.getName(), &a.q)
	return
}

func (a RemindAction) findNext() *Question {
	if a.in.isCmd() {
		return a.q.Questions.First()
	}
	next, save := a.q.next(a.in.msg.Text)
	if save {
		a.s.save(a.in.getNote(a.getName(), next.RelatesTo.SaveTag))
	}
	return next
}

func (a RemindAction) chat(q *Question) {
	a.s.send(q.Context...)
	a.s.sendWithKeyboard(q.QuestionText, q.makeKeyboard())
	return
}

type StickerChatAction struct {
	s  *session
	in *Incoming
	q  struct {
		Questions `json:"QandA"`
	}
	stickers Stickers
}

func (a StickerChatAction) getName() string {
	return chat_action
}
func (a StickerChatAction) getHint() string {
	return ""
}

func (a StickerChatAction) doFirst() {
	a.stickers = LoadKnownStickers()
	return
}

func (a StickerChatAction) loadQuestions() {
	load(a.getName(), &a.q)
	return
}

func (a StickerChatAction) findNext() *Question {
	if a.in.isSticker() {
		next, save := a.q.nextFrom(a.stickers.FindSticker(a.in.msg.Sticker.FileID).Ids...)
		if save {
			a.s.save(a.in.getNote(a.getName(), next.RelatesTo.SaveTag))
		}
		return next
	}
	next, save := a.q.next(a.in.msg.Text)
	if save {
		a.s.save(a.in.getNote(a.getName(), next.RelatesTo.SaveTag))
	}
	return next
}

func (a StickerChatAction) chat(q *Question) {
	a.s.send(q.Context...)
	a.s.sendWithKeyboard(q.QuestionText, q.makeKeyboard())
	return
}
