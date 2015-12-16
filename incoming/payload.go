package incoming

import (
	"strings"
	"time"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
)

type Payload struct {
	Sender       telebot.User
	Text         string
	Date         time.Time
	Action       string
	HasSubmisson bool
	Sticker      struct {
		Exists bool
		Id     string
	}
}

func New(msg telebot.Message) *Payload {
	return &Payload{
		Sender:       msg.Sender,
		Text:         msg.Text,
		Date:         msg.Time(),
		Action:       setAction(msg.Text),
		HasSubmisson: hasSubmisson(msg.Text),
		Sticker: struct {
			Exists bool
			Id     string
		}{
			Exists: msg.Sticker.Exists(),
			Id:     msg.Sticker.FileID,
		}}
}

func (this *Payload) HasAction() bool {
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
