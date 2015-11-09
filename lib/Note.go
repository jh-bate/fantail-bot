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

const said_tag = say_action
const chat_tag = chat_action
const remind_tag = remind_action

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
		str = fmt.Sprintf("%s", this[0].AddedOn.Format("Monday Jan 2"))
		day := this[0].AddedOn.YearDay()
		for i := range this {
			if day != this[i].AddedOn.YearDay() {
				log.Println("its a new day")
				str += fmt.Sprintf("\n\n %s", this[i].AddedOn.Format("Monday Jan 2"))
				day = this[i].AddedOn.YearDay()
			}
			str += fmt.Sprintf("\n- %s '%s'", this[i].AddedOn.Format("03:04pm"), this[i].Text)
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

func tagFromMsg(msgTxt string) string {
	words := strings.Fields(msgTxt)
	if strings.Contains(words[0], "/") {
		return words[0]
	}
	return ""
}

func NewNote(msg telebot.Message, tags ...string) Note {

	txt := msg.Text

	cmdTag := tagFromMsg(txt)

	//e.g. remove '/say' from the message
	if strings.Contains(txt, cmdTag) {
		txt = strings.TrimSpace(strings.Split(txt, cmdTag)[1])
	}

	return Note{
		WhoId:      msg.Sender.ID,
		AddedOn:    msg.Time(),
		Text:       txt,
		Tag:        strings.Join(append(tags, cmdTag), ","),
		RemindNext: time.Now().AddDate(0, 0, 7)}
}

func NewReminderNote(msg telebot.Message, tags ...string) Note {

	const remind_pos, time_pos, msg_pos = 0, 1, 2
	words := strings.Fields(msg.Text)

	days, err := strconv.Atoi(words[time_pos])

	if err != nil {
		log.Println(fmt.Errorf("Reminder format is %s", remind_action_hint).Error())
		return Note{}
	}

	what := strings.SplitAfterN(msg.Text, words[time_pos], 2)[1]

	if what == "" {
		log.Println(fmt.Errorf("Reminder format is %s", remind_action_hint).Error())
		return Note{}
	}

	return Note{
		WhoId:      msg.Sender.ID,
		AddedOn:    msg.Time(),
		Text:       strings.TrimSpace(what),
		Tag:        strings.Join(append(tags, remind_tag), ","),
		RemindNext: time.Now().AddDate(0, 0, days)}
}
