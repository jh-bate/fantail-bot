package note

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"
)

type (
	Note struct {
		UserId  int
		Added   time.Time
		Updated time.Time
		Deleted time.Time
		Tag     string
		Context []string
		Text    string
	}

	Notes []*Note
)

const (
	//tags we add to the notes
	SAID_TAG = "SAY"
	CHAT_TAG = "CHAT"
	HELP_TAG = "HELP"
)

func New(txt string, fromId int, date time.Time, tags ...string) *Note {

	answer := txt
	cmdTag := tagFromMsg(answer)

	if answer != "" && cmdTag != "" {
		//e.g. remove '/say' from the message
		if strings.Contains(answer, cmdTag) {
			answer = strings.TrimSpace(strings.Split(answer, cmdTag)[1])
		}
	}

	return &Note{
		UserId: fromId,
		Added:  date,
		Text:   answer,
		Tag:    strings.Join(append(tags, cmdTag), ","),
	}
}

func (this *Note) IsCurrent() bool {
	return this.Deleted.IsZero()
}

func (this *Note) Complete() {
	this.Deleted = time.Now()
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

func (this *Note) IsEmpty() bool {
	return this.Text == "" && this.Tag == ""
}

func (this Notes) FilterOnTag(tag string) Notes {
	n := Notes{}
	for i := range this {
		if this[i].Deleted.IsZero() && strings.Contains(strings.ToLower(this[i].Tag), strings.ToLower(tag)) {
			n = append(n, this[i])
		}
	}
	return n
}

func (this Notes) FilterOnTxt(txt string) Notes {
	n := Notes{}
	for i := range this {
		if this[i].Deleted.IsZero() && strings.Contains(strings.ToLower(this[i].Text), strings.ToLower(txt)) {
			n = append(n, this[i])
		}
	}
	return n
}

func (this Notes) GetWords() []string {
	var w []string
	for i := range this {
		w = append(w, strings.Fields(this[i].Text)...)
		w = append(w, strings.Fields(this[i].Tag)...)
	}
	return w
}

func (this Notes) NewerThan(daysAgo int) Notes {
	var r Notes
	daysAgoDate := time.Now().AddDate(0, 0, -daysAgo)

	for i := range this {
		if this[i].Added.After(daysAgoDate) {
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

func (this Notes) OldestFirst() Notes {
	sort.Sort(ByDate(this))
	return this
}

func (this Notes) MostRecent() *Note {
	if len(this) > 0 {
		this.OldestFirst()
		return this[len(this)-1]
	}
	return &Note{}
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
