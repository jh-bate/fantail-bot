package lib

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
)

const said_tag = say_cmd
const chat_tag = chat_cmd
const remind_tag = remind_cmd

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
	return fmt.Sprintf("On %s you said '%s'", this.AddedOn.Format("Mon Jan 2 03:04pm"), this.Text)
}

func (this Notes) FilterBy(tag string) Notes {
	var n Notes
	for i := range this {
		if this[i].CompletedOn.IsZero() && strings.Contains(this[i].Tag, tag) {
			n = append(n, this[i])
		}
	}
	return n
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

type ByDate Notes

func (this ByDate) Len() int           { return len(this) }
func (this ByDate) Swap(i, j int)      { this[i], this[j] = this[j], this[i] }
func (this ByDate) Less(i, j int) bool { return this[i].AddedOn.Before(this[j].AddedOn) }

func (this Notes) ToString() string {
	str := ""
	if len(this) > 0 {
		str = fmt.Sprintf("%s \n\n", this[0].AddedOn.Format("Mon Jan 2"))
		for i := range this {
			str += fmt.Sprintf("  - %s '%s'", this[i].AddedOn.Format("Mon 03:04pm"), this[i].Text)
		}
	}

	return str
}

func (this Notes) SortByDate() Notes {
	sort.Sort(ByDate(this))
	return this
}

func (this Note) IsEmpty() bool {
	return this.Text == "" && this.Tag == ""
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

func NewReminderNote(msg telebot.Message) Note {

	const remind_pos, time_pos, msg_pos = 0, 1, 2
	words := strings.Fields(msg.Text)

	days, err := strconv.Atoi(words[time_pos])

	if err != nil {
		log.Println(fmt.Errorf("Reminder format is %s", remind_cmd_hint).Error())
		return Note{}
	}

	what := strings.SplitAfterN(msg.Text, words[time_pos], 2)[1]

	if what == "" {
		log.Println(fmt.Errorf("Reminder format is %s", remind_cmd_hint).Error())
		return Note{}
	}

	return Note{
		WhoId:      msg.Sender.ID,
		AddedOn:    msg.Time(),
		Text:       strings.TrimSpace(what),
		Tag:        remind_tag,
		RemindNext: time.Now().AddDate(0, 0, days)}
}
