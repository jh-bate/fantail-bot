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
			//already exists so remove and then we will add the new one
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
