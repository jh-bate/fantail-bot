package lib

import (
	"time"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
)

type (
	User struct {
		Id        int         `json:"id"`
		Learnings []Learnt    `json:"learnings"`
		Helped    []time.Time `json:"helped"`

		Notes `json:"-"`
	}

	Learnt struct {
		On       time.Time `json:"on"`
		Positive bool      `json:"positive"`
	}

	Users []*User
)

func (this *User) FollowUpAbout() Notes {
	this.Helped = append(this.Helped, time.Now())
	return this.Notes.FilterOnTag(help_tag).SortByDate()
}

func (this *User) IsPostive(days int) bool {
	classify := NewClassification()
	positive := classify.ArePositive(this.Notes.ForNextDays(days).GetWords())
	this.Learnings = append(this.Learnings, Learnt{On: time.Now(), Positive: positive})
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
