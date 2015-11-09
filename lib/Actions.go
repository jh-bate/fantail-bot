package lib

import (
	"fmt"
	"log"
)

type Action interface {
	getName() string
	getHint() string
	do()
}

const (
	start_action, help_action = "/start", "/help"

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
	log.Println("asked ", in.getCmd())
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
func (a SayAction) do() {
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
func (a HelpAction) do() {
	helpInfo := fmt.Sprintf("%s %s %s %s ",
		fmt.Sprintf("Hey %s! We can't do it all but we can:\n\n", a.in.sender().Username),
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
func (a ChatAction) do() {
	//nothing to do
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
func (a ReviewAction) do() {
	log.Println("doing review ...")
	n := a.s.getNotes(a.in.msg)

	saidTxt := fmt.Sprintf("%s you said: \n\n %s", a.in.sender().Username, n.FilterBy(said_tag).ToString())
	a.s.send(saidTxt)
	talkedTxt := fmt.Sprintf("%s we talked about: \n\n %s", a.in.sender().Username, n.FilterBy(chat_tag).ToString())
	a.s.send(talkedTxt)
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
func (a RemindAction) do() {
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
