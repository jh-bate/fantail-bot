package user

import (
	"time"

	"github.com/jh-bate/fantail-bot/note"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
)

type (
	User struct {
		Id     int        `json:"id"`
		Learnt []learning `json:"learnt"`
		Helped []help     `json:"helped"`

		note.Notes `json:"-"`
	}

	learning struct {
		Date     time.Time `json:"date"`
		Period   int       `json:"period"`
		Positive bool      `json:"positive"`
	}

	help struct {
		Date    time.Time `json:"date"`
		AskedOn time.Time `json:"askedOn"`
		Topic   string    `json:"topic"`
	}

	Users []*User
)

func New() *User {
	return &User{}
}

func (this *User) NeedsHelp() note.Notes {
	helpWith := this.Notes.FilterOnTag(note.HELP_TAG).SortByDate()

	for i := range helpWith {
		this.Helped = append(this.Helped, help{Date: time.Now(), Topic: helpWith[i].Text, AskedOn: helpWith[i].Added})
	}
	return helpWith
}

func (this *User) LearnAbout(days int) bool {
	classify := NewClassification()
	positive := classify.ArePositive(this.Notes.ForNextDays(days).GetWords())
	this.Learnt = append(this.Learnt, learning{Date: time.Now(), Positive: positive, Period: days})
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
