package user

import (
	"time"

	"github.com/jh-bate/fantail-bot/note"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
)

type (
	User struct {
		Id     int        `json:"id"`
		Learnt []Learning `json:"learnt"`
		Helped []Help     `json:"helped"`

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

func New(id int) *User {
	return &User{Id: id}
}

func (this *User) NeedsHelp() note.Notes {
	helpWith := this.Notes.FilterOnTag(note.HELP_TAG).OldestFirst()
	didHelp := false
	for i := range helpWith {
		didHelp = true
		this.Helped = append(this.Helped, Help{Date: time.Now(), Topic: helpWith[i].Text, AskedOn: helpWith[i].Added})
	}
	if didHelp {
		this.Save()
	}
	return helpWith
}

func (this *User) LearnAbout(days int) bool {
	classify := NewClassification()
	positive := classify.ArePositive(this.Notes.NewerThan(days).GetWords())
	this.Learnt = append(this.Learnt, Learning{Date: time.Now(), Positive: positive, Period: days})
	this.Save()
	return positive
}

func (this *User) ToBotUser() telebot.User {
	return telebot.User{ID: this.Id}
}

func (this *User) AddOrUpdate(users Users) Users {
	var updated Users

	for i := range users {
		if users[i].Id != this.Id {
			//rebuild the list from those that don't match the user we are trying to add or update
			updated = append(updated, users[i])
		}
	}
	return append(updated, this)
}

func (this Users) GetUser(id int) *User {

	for i := range this {
		if this[i].Id == id {
			return this[i]
		}
	}
	return nil
}
