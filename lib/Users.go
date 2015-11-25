package lib

import (
	"time"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
)

type (
	User struct {
		id       int
		lastChat time.Time
		recent   Notes
	}

	Users []*User
)

func (this *User) GetReminders() Notes {
	return this.recent.FilterBy(remind_tag)
}

func (this *User) HelpAskedFor() Notes {
	return this.recent.FilterBy(help_tag)
}

func (this *User) ToBotUser() telebot.User {
	return telebot.User{ID: this.id}
}

func (this *User) AddOrUpdate(users Users) Users {
	var updated Users

	for i := range users {
		if users[i].id != this.id {
			//rebuild the list from those that don't match the user we are trying to add or update
			updated = append(updated, users[i])
		}
	}
	return append(updated, this)
}

func (this Users) GetUser(id int) *User {

	for i := range this {
		if this[i].id == id {
			return this[i]
		}
	}
	return nil
}
