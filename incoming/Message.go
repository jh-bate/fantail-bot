package incoming

import (
	"strings"
	"time"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
)

type Message struct {
	Sender       telebot.User
	Text         string
	Date         time.Time
	Action       string
	HasSubmisson bool
	Sticker      struct {
		Has bool
		Id  string
	}
}

func New(msg telebot.Message) *Message {
	return &Message{
		Sender:       msg.Sender,
		Text:         msg.Text,
		Date:         msg.Time(),
		Action:       setAction(msg.Text),
		HasSubmisson: hasSubmisson(msg.Text),
		Sticker: struct {
			Has bool
			Id  string
		}{
			Has: msg.Sticker.Exists(),
			Id:  msg.Sticker.FileID,
		}}
}

func (this *Message) HasAction() bool {
	return this.Action != ""
}

func isAction(txt string) bool {
	// `/chat` is an action
	if txt != "" && strings.Contains(txt, "/") {
		return strings.Contains(strings.Fields(txt)[0], "/")
	}
	return false
}

func setAction(txt string) string {
	if isAction(txt) {
		return strings.ToLower(strings.Fields(txt)[0])
	}
	return ""
}

func hasSubmisson(txt string) bool {
	return isAction(txt) && len(strings.Fields(txt)) > 1
}
