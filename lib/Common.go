package lib

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
)

const (
	typing_action = "typing"
)

type (
	Question struct {
		RelatesTo struct {
			Answers []string `json:"answers"`
			Save    bool     `json:"save"`
			SaveTag string   `json:"saveTag"`
		} `json:"relatesTo"`
		Context         []string `json:"context"`
		QuestionText    string   `json:"question"`
		PossibleAnswers []string `json:"answers"`
	}

	questions []*Question

	Details struct {
		Bot     *telebot.Bot
		User    telebot.User
		Storage *Storage
	}

	Process interface {
		Run(input <-chan telebot.Message)
	}

	Keyboard [][]string
)

func (d *Details) send(msg string) {
	d.takeThoughtfulPause()

	if strings.Contains(msg, "%s") {
		msg = fmt.Sprintf(msg, d.User.FirstName)
	}

	d.Bot.SendMessage(
		d.User,
		msg,
		nil,
	)
	return
}

func (d *Details) sendWithKeyboard(msg string, kb Keyboard) {
	d.takeThoughtfulPause()

	if strings.Contains(msg, "%s") {
		msg = fmt.Sprintf(msg, d.User.FirstName)
	}

	d.Bot.SendMessage(
		d.User,
		msg,
		&telebot.SendOptions{
			ReplyMarkup: telebot.ReplyMarkup{
				ForceReply:      false,
				CustomKeyboard:  kb,
				ResizeKeyboard:  true,
				OneTimeKeyboard: true,
			},
		},
	)
	return
}

func (d *Details) takeThoughtfulPause() {
	d.Bot.SendChatAction(d.User, typing_action)
	time.Sleep(1 * time.Second)
	return
}

func hasSubmisson(txt string, cmds ...string) bool {
	if isCmd(txt, cmds...) {
		for i := range cmds {
			if len(strings.SplitAfter(txt, cmds[i])) > 2 {
				return true
			}
		}
	}
	return false
}

func isCmd(txt string, cmds ...string) bool {
	for i := range cmds {
		if strings.Contains(txt, cmds[i]) {
			return true
		}
	}
	return false
}

func (d *Details) saveAsReminder(msg telebot.Message) error {
	r, err := NewReminderNote(msg)
	if err != nil {
		return err
	}
	return d.Storage.Save(fmt.Sprintf("%d", d.User.ID), r)
}

func (d *Details) save(msg telebot.Message, tags ...string) {
	if d.Storage == nil {
		log.Println(StorageInitErr.Error())
		return
	}
	log.Println("Saving", msg.Text)

	err := d.Storage.Save(fmt.Sprintf("%d", d.User.ID), NewNote(msg, tags...))

	if err != nil {
		log.Println(err.Error())
		log.Println(StorageSaveErr.Error())
	}
	return
}

func (d *Details) getReminders(userId string) Notes {
	all, err := d.Storage.Get(userId)
	if err == nil {
		return all.FilterReminders()
	}
	return Notes{}
}

func (d *Details) getNotes(userId string) Notes {
	all, err := d.Storage.Get(userId)
	if err == nil {
		return all.FilterNotes()
	}
	return Notes{}
}
