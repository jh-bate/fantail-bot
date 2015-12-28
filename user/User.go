package user

import (
	"time"

	"github.com/jh-bate/fantail-bot/note"
)

type (
	User struct {
		Id string `json:"id"`
		//don't want to persist this detail in the store
		Name   string                 `json:"-"`
		Learnt map[time.Time]Learning `json:"learnt"`
		Helped map[time.Time]Help     `json:"helped"`

		note.Notes `json:"-"`
	}

	Learning struct {
		Date     time.Time `json:"date"`
		Period   int       `json:"period"`
		Positive bool      `json:"positive"`
	}

	Help struct {
		Date    time.Time `json:"date"`
		AskedOn time.Time `json:"askedOn"`
		Topic   string    `json:"topic"`
	}

	Users []*User
)

func New(id string) *User {
	return &User{Id: id}
}

func (this *User) NeedsHelp() note.Notes {
	helpWith := this.Notes.FilterOnTag(note.HELP_TAG).OldestFirst()
	didHelp := false
	for i := range helpWith {
		didHelp = true
		now := time.Now()
		this.Helped[now] = Help{Date: now, Topic: helpWith[i].Text, AskedOn: helpWith[i].Added}
	}
	if didHelp {
		this.Upsert()
	}
	return helpWith
}

func (this *User) LearnAbout(days int) bool {
	classify := NewClassification()
	positive := classify.ArePositive(this.Notes.NewerThan(days).GetWords())
	now := time.Now()
	this.Learnt[now] = Learning{Date: now, Positive: positive, Period: days}
	this.Upsert()
	return positive
}

func (this Users) GetUser(id string) *User {

	for i := range this {
		if this[i].Id == id {
			return this[i]
		}
	}
	return nil
}
