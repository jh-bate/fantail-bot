package lib

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
)

const reminder_tag = remind_cmd

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

	return this.RemindNext.Year() == today.Year() &&
		this.RemindNext.YearDay() == today.YearDay()
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

func (this *Note) UpdateRemindNext() {
	today := time.Now()
	this.RemindNext = today.AddDate(0, 0, 7)
	return
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
		if this[i].CompletedOn.IsZero() && strings.Contains(this[i].Tag, reminder_tag) == false {
			n = append(n, this[i])
		}
	}
	return n
}

func NewNote(msg telebot.Message, tags ...string) Note {
	return Note{
		WhoId:      msg.Sender.ID,
		AddedOn:    msg.Time(),
		Text:       msg.Text,
		Tag:        strings.Join(tags, ","),
		RemindNext: time.Now().AddDate(0, 0, 7)}
}

func NewReminderNote(msg telebot.Message) (Note, error) {

	const remind_pos, me_pos, in_pos, time_pos, to_pos, msg_pos = 0, 1, 2, 3, 4, 5
	const remind, me, in, to = "/remind", "me", "in", "to"
	words := strings.Fields(msg.Text)

	if strings.ToLower(words[remind_pos]) != remind ||
		strings.ToLower(words[me_pos]) != me ||
		strings.ToLower(words[in_pos]) != in ||
		strings.ToLower(words[to_pos]) != to {
		return Note{}, errors.New("format is " + remind_cmd_hint)
	}

	days, err := strconv.Atoi(words[time_pos])
	if err != nil {
		return Note{}, err
	}
	what := words[msg_pos]

	return Note{
		WhoId:      msg.Sender.ID,
		AddedOn:    msg.Time(),
		Text:       what,
		Tag:        reminder_tag,
		RemindNext: time.Now().AddDate(0, 0, days)}, nil
}
