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

func (this Users) GetUser(id int) *User {

	for i := range this {
		if this[i].id == id {
			return this[i]
		}
	}
	return nil
}

func (this Users) AddUser(u *User) {
	this = append(this, u)
	return
}
