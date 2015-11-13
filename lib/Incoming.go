package lib

import (
	"strings"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
)

type Incoming struct {
	msg    telebot.Message
	action Action
}

func newIncoming(msg telebot.Message) *Incoming {
	return &Incoming{msg: msg}
}

func (this Incoming) hasSubmisson() bool {
	//eg /say some stuff
	return this.isCmd() && len(strings.Fields(this.msg.Text)) > 1
}

func (this Incoming) getAction(s *session, prevActionName string) Action {
	return NewAction(&this, s, prevActionName)
}

func (this Incoming) getNote(tags ...string) Note {
	if strings.Contains(this.msg.Text, remind_action) {
		return NewReminderNote(this.msg)
	}
	return NewNote(this.msg, this.getCmd())
}

func (this Incoming) isSticker() bool {
	return this.msg.Sticker.Exists()
}

func (this Incoming) sender() telebot.User {
	return this.msg.Sender
}

func (this Incoming) getCmd() string {
	if this.isCmd() {
		if strings.Contains(this.msg.Text, "/") {
			return strings.Fields(this.msg.Text)[0]
		}
	}
	return ""
}

func (this Incoming) cmdMatches(cmds ...string) bool {
	if this.isCmd() {
		for i := range cmds {
			if cmds[i] == this.getCmd() {
				return true
			}
		}
	}
	return false
}

func (this Incoming) submissonMatches(cmds ...string) bool {
	if this.hasSubmisson() {
		for i := range cmds {
			if cmds[i] == this.getCmd() {
				return true
			}
		}
	}
	return false
}

func (this Incoming) isCmd() bool {
	//eg /cmd
	if this.msg.Text != "" && strings.Contains(this.msg.Text, "/") {
		return strings.Contains(strings.Fields(this.msg.Text)[0], "/")
	}

	return false
}
