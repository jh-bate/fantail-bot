package lib

import (
	"errors"
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
	Reminder struct {
		WhoId       int
		AddedOn     time.Time
		RemindNext  time.Time
		CompletedOn time.Time
		Tag         string
		Text        string
	}

	Reminders []*Reminder

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

func (d *Details) save(msg telebot.Message, tags ...string) {
	if d.Storage == nil {
		log.Println(StorageInitErr.Error())
		return
	}
	log.Println("Saving", msg.Text)

	r := Reminder{
		WhoId:      msg.Sender.ID,
		AddedOn:    msg.Time(),
		Text:       msg.Text,
		Tag:        fmt.Sprintln(tags),
		RemindNext: time.Now().AddDate(0, 0, 7)}

	log.Println("saving... ", r)

	err := d.Storage.Save(fmt.Sprintf("%d", d.User.ID), r)

	if err != nil {
		log.Println(err.Error())
		log.Println(StorageSaveErr.Error())
	}
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
				//log.Println("Check if submisson", txt)
				return true
			}
		}
	}
	return false
}

func isCmd(txt string, cmds ...string) bool {
	//log.Println("Check if cmd", txt)
	for i := range cmds {
		if strings.Contains(txt, cmds[i]) {
			return true
		}
	}
	return false
}

func (d *Details) saveReminder(msg telebot.Message) error {

	const remind_pos, me_pos, in_pos, time_pos, to_pos, msg_pos = 0, 1, 2, 3, 4, 5
	const remind, me, in, to = "/remind", "me", "in", "to"
	words := strings.Fields(msg.Text)

	if strings.ToLower(words[remind_pos]) != remind ||
		strings.ToLower(words[me_pos]) != me ||
		strings.ToLower(words[in_pos]) != in ||
		strings.ToLower(words[to_pos]) != to {
		return errors.New("format is /remind me to <days> do <msg>")
	}

	days, err := strconv.Atoi(words[time_pos])
	if err != nil {
		return err
	}
	what := words[msg_pos]

	r := Reminder{
		WhoId:      msg.Sender.ID,
		AddedOn:    msg.Time(),
		Text:       what,
		Tag:        remind_cmd,
		RemindNext: time.Now().AddDate(0, 0, days)}

	return d.Storage.Save(fmt.Sprintf("%d", d.User.ID), r)

}
