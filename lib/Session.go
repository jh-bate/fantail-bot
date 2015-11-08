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
		Action  *action
	}

	Keyboard [][]string
)

func newSession(b *telebot.Bot, s *Storage) *session {
	return &session{Bot: b, Storage: s}
}

func (s *session) addDetails(msg telebot.Message) {
	s.Action = newAction(msg)
	s.User = msg.Sender
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

func (s *session) save(n Note) {
	if s.Storage == nil {
		log.Println(FantailStorageErr.Error())
		return
	}
	if n.IsEmpty() {
		log.Println("Nothing to save")
		return
	}
	log.Println("Saving", n.Text)

	err := s.Storage.Save(fmt.Sprintf("%d", s.User.ID), n)

	if err != nil {
		log.Println(err.Error())
		log.Println(FantailSaveErr.Error())
	}
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

func (s *session) getNotes(msg telebot.Message) Notes {

	days := daysFromText(msg.Text)
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
		return all.FilterBy(topic).SortByDate()[0]
	}
	return nil
}
