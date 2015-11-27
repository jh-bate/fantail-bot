package lib

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
)

const said_tag = say_action
const chat_tag = chat_action
const help_tag = "HELP"

type (
	Note struct {
		UserId    int
		Added     time.Time
		Updated   time.Time
		Remind    time.Time
		Completed time.Time
		Tag       string
		Context   []string
		Text      string
	}

	Notes []*Note
)

func (this *Note) RemindToday() bool {
	today := time.Now()

	if this.Remind.Before(today) == false {
		return this.Remind.Year() == today.Year() &&
			this.Remind.YearDay() == today.YearDay()
	}
	return true
}

func (this *Note) IsCurrent() bool {
	return this.Completed.IsZero()
}

func (this *Note) Complete() {
	this.Completed = time.Now()
	return
}

func (this *Note) Update() {
	this.Updated = time.Now()
	return
}

func (this *Note) SetContext(context ...string) {
	this.Context = context
	return
}

func (this *Note) ToString() string {
	return fmt.Sprintf("On %s you said '%s'", this.Added.Format("Mon Jan 2 03:04pm"), this.Text)
}

func (this Notes) FilterBy(tag string) Notes {
	n := Notes{}
	for i := range this {
		if this[i].Completed.IsZero() && strings.Contains(this[i].Tag, tag) {
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
		if this[i].Remind.Before(t) {
			r = append(r, this[i])
		}
	}
	return r
}

type ByDate Notes

func (this ByDate) Len() int           { return len(this) }
func (this ByDate) Swap(i, j int)      { this[i], this[j] = this[j], this[i] }
func (this ByDate) Less(i, j int) bool { return this[i].Added.Before(this[j].Added) }

func (this Notes) ToString() string {
	str := ""
	if len(this) > 0 {
		str = fmt.Sprintf("%s", this[0].Added.Format("Monday Jan 2"))
		day := this[0].Added.YearDay()
		for i := range this {
			if day != this[i].Added.YearDay() {
				log.Println("its a new day")
				str += fmt.Sprintf("\n\n %s", this[i].Added.Format("Monday Jan 2"))
				day = this[i].Added.YearDay()
			}
			str += fmt.Sprintf("\n- %s '%s'", this[i].Added.Format("03:04pm"), this[i].Text)
		}
	}

	return str
}

func (this Notes) SortByDate() Notes {
	sort.Sort(ByDate(this))
	return this
}

func (this Notes) MostRecent() *Note {
	if len(this) > 0 {
		this.SortByDate()
		return this[0]
	}
	return &Note{}
}

func (this Note) IsEmpty() bool {
	return this.Text == "" && this.Tag == ""
}

func tagFromMsg(msgTxt string) string {
	if msgTxt != "" {
		words := strings.Fields(msgTxt)
		if strings.Contains(words[0], "/") {
			return words[0]
		}
	}
	return ""
}

func NewNote(msg telebot.Message, tags ...string) *Note {

	answer := msg.Text
	cmdTag := tagFromMsg(answer)

	if answer != "" && cmdTag != "" {
		//e.g. remove '/say' from the message
		if strings.Contains(answer, cmdTag) {
			answer = strings.TrimSpace(strings.Split(answer, cmdTag)[1])
		}
	}

	return &Note{
		UserId: msg.Sender.ID,
		Added:  msg.Time(),
		Text:   answer,
		Tag:    strings.Join(append(tags, cmdTag), ","),
	}
}
