package lib

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
)

const (
	typing_action = "typing"
)

type (
	session struct {
		Bot     *telebot.Bot
		User    telebot.User
		Storage *Storage
		in      *Incoming
	}

	Keyboard [][]string
)

func newSession(b *telebot.Bot, s *Storage) *session {
	return &session{Bot: b, Storage: s}
}

// Utility functions that wrap sent messages

func (s *session) setIncoming(in *Incoming) *session {
	s.in = in
	return s
}

func (s *session) getIncoming() *Incoming {
	return s.in
}

func (s *session) getActionForSent(prevAction string) Action {
	return s.in.getAction(s, prevAction)
}

func (s *session) getActionNameForSent() (bool, string) {
	log.Println("Check getActionName ", s.getSentMsgText())
	if strings.Contains(s.getSentMsgText(), "/") {
		log.Println("Check getActionName ", strings.Fields(s.getSentMsgText())[0])
		return true, strings.Fields(s.getSentMsgText())[0]
	}
	return false, ""
}

func (s *session) getSentUsername() string {
	return s.in.sender().Username
}

func (s *session) getSentMsg() telebot.Message {
	return s.in.msg
}

func (s *session) getSender() telebot.User {
	return s.in.msg.Sender
}

func (s *session) getSentMsgText() string {
	return s.in.msg.Text
}

func (s *session) sentAsCommand() bool {
	return s.in.isCmd()
}

func (s *session) getSentCommand() string {
	return s.in.getCmd()
}

func (s *session) sentAsSticker() bool {
	return s.in.isSticker()
}

func (s *session) sentAsSubmission() bool {
	return s.in.hasSubmisson()
}

func (s *session) setSentMsgText(text string) {
	s.in.msg.Text = text
	return
}

func (s *session) getSentStickerId() string {
	if s.sentAsSticker() {
		return s.in.msg.Sticker.FileID
	}
	return ""
}

func (s *session) saveWithContext(context []string, tags ...string) {

	if s.Storage == nil {
		log.Println(FantailStorageErr.Error())
		return
	}
	n := s.in.createNote(tags...)
	n.SetContext(context...)
	if n.IsEmpty() {
		log.Println("Nothing to save")
		return
	}

	err := s.Storage.Save(fmt.Sprintf("%d", s.User.ID), n)

	if err != nil {
		log.Println(err.Error())
		log.Println(FantailSaveErr.Error())
	}
	return
}

func (s *session) save(tags ...string) {

	if s.Storage == nil {
		log.Println(FantailStorageErr.Error())
		return
	}

	n := s.in.createNote(tags...)

	if n.IsEmpty() {
		log.Println("Nothing to save")
		return
	}

	err := s.Storage.Save(fmt.Sprintf("%d", s.User.ID), n)

	if err != nil {
		log.Println(err.Error())
		log.Println(FantailSaveErr.Error())
	}
	return
}

func (s *session) send(msgs ...string) {

	for i := range msgs {
		s.takeThoughtfulPause()

		msg := msgs[i]
		if strings.Contains(msg, "%s") {
			msg = fmt.Sprintf(msg, s.User.FirstName)
		}

		s.Bot.SendMessage(
			s.User,
			msg,
			&telebot.SendOptions{ParseMode: telebot.ModeMarkdown},
		)
	}
	return
}

func (s *session) sendWithKeyboard(msg string, kb Keyboard) {
	s.takeThoughtfulPause()

	if strings.Contains(msg, "%s") {
		msg = fmt.Sprintf(msg, s.User.FirstName)
	}

	s.Bot.SendMessage(
		s.User,
		msg,
		&telebot.SendOptions{
			ReplyMarkup: telebot.ReplyMarkup{
				ForceReply:      false,
				CustomKeyboard:  kb,
				ResizeKeyboard:  false,
				OneTimeKeyboard: true,
			},
			ParseMode: telebot.ModeMarkdown,
		},
	)
	return
}

func (s *session) takeThoughtfulPause() {
	s.Bot.SendChatAction(s.User, typing_action)
	time.Sleep(1 * time.Second)
	return
}

func daysFromText(txt string) int {
	const cmd_pos, time_pos = 0, 1
	words := strings.Fields(txt)

	if len(words) == 1 {
		//dafault if zero
		return 0
	}

	days, err := strconv.Atoi(words[time_pos])

	if err != nil {
		log.Println("error getting number of days", err.Error())
	}

	log.Println("days", days)

	return days
}

func (s *session) getNotes() Notes {
	days := daysFromText(s.in.msg.Text)
	all, err := s.Storage.Get(fmt.Sprintf("%d", s.User.ID))
	if err == nil {
		if days > 0 {
			return all.ForNextDays(days)
		}
		return all
	}
	return Notes{}
}

func (s *session) getLastChatForTopic(topic string) *Note {

	all, err := s.Storage.Get(fmt.Sprintf("%d", s.User.ID))
	if err == nil {
		return all.FilterOnTag(topic).SortByDate()[0]
	}
	return nil
}
