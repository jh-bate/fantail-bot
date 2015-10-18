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
	Said struct {
		FromId         int
		ToId           int
		When           time.Time
		Remind         bool
		RemindComplete time.Time
		Text           string
	}

	Chat []*Said

	Question struct {
		RelatesTo struct {
			Answers []string `json:"answers"`
			Save    bool     `json:"save"`
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

func (d *Details) save(msg telebot.Message) {
	if d.Storage == nil {
		log.Println(StorageInitErr.Error())
		return
	}
	log.Println("Saving", msg.Text)
	err := d.Storage.Save(fmt.Sprintf("%d", d.User.ID), Said{FromId: msg.Sender.ID, ToId: msg.Chat.ID, When: msg.Time(), Text: msg.Text, Remind: true})
	if err != nil {
		log.Println(StorageSaveErr.Error())
	}
	return
}

func (d *Details) takeThoughtfulPause() {
	d.Bot.SendChatAction(d.User, typing_action)
	time.Sleep(1 * time.Second)
	return
}
