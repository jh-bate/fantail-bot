package lib

import (
	"fmt"
	"log"
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

var Config struct {
	Questions `json:"QandA"`
}

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
	if a.sentAsSubmission() {
		a.session.save()
	}
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
	ConfigLoader(&Config, "say.json")
	return Config.Questions
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
	if a.sentAsSubmission() {
		a.session.save(help_tag)
	}
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
	ConfigLoader(&Config, "ask.json")
	return Config.Questions
}

type HelpAction struct {
	*session
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
		"Stickers - we have those to help express yourself!! \n\n [Get them here](https://telegram.me/addstickers/betes)",
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
	ConfigLoader(&Config, "chat.json")
	return Config.Questions
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

	saidTxt := fmt.Sprintf("%s you said: \n\n %s", a.session.getSentUsername(), n.FilterOnTag(said_tag).ToString())
	a.session.send(saidTxt)
	talkedTxt := fmt.Sprintf("%s we talked about: \n\n %s", a.session.getSentUsername(), n.FilterOnTag(chat_tag).ToString())
	a.session.send(talkedTxt)
	return a
}

func (a ReviewAction) getQuestions() Questions {
	ConfigLoader(&Config, "review.json")
	return Config.Questions
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
	ConfigLoader(&Config, "chat.json")
	return Config.Questions
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
