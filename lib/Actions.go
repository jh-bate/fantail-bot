package lib

import (
	"fmt"
	"log"
)

type Action interface {
	getName() string
	getHint() string
	quickWin()
}

const (
	say_action      = "/say"
	say_action_hint = "/say [what you want to say]"

	chat_action = "/chat"

	review_action      = "/review"
	review_action_hint = "/review <days>"

	remind_action      = "/remind"
	remind_action_hint = "/remind in <days> to <msg>"
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
	}
	return &HelpAction{in: &in, s: s}
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
func (a SayAction) quickWin() {
	log.Println("say ...")
	a.s.save(a.in.getNote())
	return
}

type HelpAction struct {
	s  *session
	in *Incoming
}

func (a HelpAction) getName() string {
	return say_action
}
func (a HelpAction) getHint() string {
	return say_action_hint
}
func (a HelpAction) quickWin() {
	helpInfo := fmt.Sprintf("%s %s %s %s ",
		"Options:\n\n",
		chat_action+" - to have a *quick chat* about what your upto \n\n",
		say_action_hint+" - to say *anything* thats on your mind \n\n",
		review_action_hint+" - to review what has been happening \n\n",
	)
	a.s.send(helpInfo)
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
func (a ChatAction) quickWin() {
	return
}

type ReviewAction struct {
	s  *session
	in *Incoming
}

func (a ReviewAction) getName() string {
	return review_action
}
func (a ReviewAction) getHint() string {
	return review_action_hint
}
func (a ReviewAction) quickWin() {
	log.Println("doing review ...")
	n := a.s.getNotes(a.in.msg)

	saidTxt := "You said: \n\n" + n.FilterBy(said_tag).ToString()
	a.s.send(saidTxt)
	chattedTxt := "We chatted: \n\n" + n.FilterBy(chat_tag).ToString()
	a.s.send(chattedTxt)
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
func (a RemindAction) quickWin() {
	log.Println("remind me ...")
	a.s.save(a.in.getNote())
	return
}

func (a RemindAction) getDays() int {
	return 0
}

func (a RemindAction) getMessage() string {
	return ""
}
