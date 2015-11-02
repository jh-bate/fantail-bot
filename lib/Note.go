package lib

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
)

const reminder_tag = remind_cmd
const said_tag = say_cmd
const chat_tag = chat_cmd

type (
	Note struct {
		WhoId       int
		AddedOn     time.Time
		RemindNext  time.Time
		CompletedOn time.Time
		Tag         string
		Text        string
	}

	Notes []*Note
)

func (this *Note) RemindToday() bool {
	today := time.Now()

	if this.RemindNext.Before(today) == false {
		return this.RemindNext.Year() == today.Year() &&
			this.RemindNext.YearDay() == today.YearDay()
	}
	return true
}

func (this *Note) IsReminder() bool {
	return strings.Contains(this.Tag, reminder_tag)
}

func (this *Note) IsCurrentReminder() bool {
	return this.IsReminder() && this.IsCurrent()
}

func (this *Note) IsCurrent() bool {
	return this.CompletedOn.IsZero()
}

func (this *Note) Complete() {
	this.CompletedOn = time.Now()
	return
}

func (this *Note) UpdateRemindNext() {
	today := time.Now()
	this.RemindNext = today.AddDate(0, 0, 7)
	return
}

func (this *Note) ToString() string {
	return strings.Join([]string{this.AddedOn.Format(time.Stamp), this.Text}, " ")
}

func (this Notes) FilterReminders() Notes {
	var r Notes
	for i := range this {
		if this[i].CompletedOn.IsZero() && strings.Contains(this[i].Tag, reminder_tag) {
			r = append(r, this[i])
		}
	}
	return r
}

func (this Notes) FilterNotes() Notes {
	var n Notes
	for i := range this {
		if this[i].CompletedOn.IsZero() && strings.Contains(this[i].Tag, said_tag) {
			n = append(n, this[i])
		}
	}
	return n
}

func (this Notes) FilterChat(topic string) Notes {
	var n Notes
	for i := range this {
		if this[i].CompletedOn.IsZero() && strings.Contains(this[i].Tag, chat_tag) && strings.Contains(this[i].Tag, topic) {
			n = append(n, this[i])
		}
	}

	return n
}

func (this Notes) ForToday() Notes {
	var r Notes
	for i := range this {
		if this[i].RemindToday() {
			r = append(r, this[i])
		}
	}
	return r
}

func (this Notes) ForNextDays(days int) Notes {
	var r Notes
	t := time.Now()

	t.AddDate(0, 0, days)

	log.Println("getting all before ", t.Format(time.Stamp))

	for i := range this {
		if this[i].RemindNext.Before(t) {
			r = append(r, this[i])
		}
	}
	return r
}

func NewNote(msg telebot.Message, tags ...string) Note {

	txt := msg.Text

	//e.g. remove '/say' from the message
	if strings.Contains(txt, tags[0]) {
		txt = strings.TrimSpace(strings.Split(txt, tags[0])[1])
	}

	return Note{
		WhoId:      msg.Sender.ID,
		AddedOn:    msg.Time(),
		Text:       txt,
		Tag:        strings.Join(tags, ","),
		RemindNext: time.Now().AddDate(0, 0, 7)}
}

func NewReminderNote(msg telebot.Message) (Note, error) {

	const remind_pos, time_pos, msg_pos = 0, 1, 2
	words := strings.Fields(msg.Text)

	days, err := strconv.Atoi(words[time_pos])

	if err != nil {
		return Note{}, errors.New("format is " + remind_cmd_hint)
	}

	what := strings.SplitAfterN(msg.Text, words[time_pos], 2)[1]

	if what == "" {
		return Note{}, errors.New("format is " + remind_cmd_hint)
	}

	return Note{
		WhoId:      msg.Sender.ID,
		AddedOn:    msg.Time(),
		Text:       strings.TrimSpace(what),
		Tag:        reminder_tag,
		RemindNext: time.Now().AddDate(0, 0, days)}, nil
}
