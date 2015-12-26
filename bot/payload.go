package bot

import (
	"strings"
	"time"

	"github.com/jh-bate/fantail-bot/user"
)

//basic incomming payload for text based messages
type Payload struct {
	Sender       int
	Text         string
	Date         time.Time
	Action       string
	HasSubmisson bool
}

func New(senderId int, msgText string, date time.Time) *Payload {
	return &Payload{
		Sender:       senderId,
		Text:         msgText,
		Date:         date,
		Action:       setAction(msgText),
		HasSubmisson: hasSubmisson(msgText),
	}
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

type Session struct {
	actionRunAs string
	bot         Bot
}

var monitor *FollowUp

func NewSession(b Bot) *Session {

	s := &Session{bot: b}

	monitor = NewFollowUp(s)
	monitor.Start()
	return s
}

func (s *Session) Respond(data *Payload) {

	user.New(data.Sender).Upsert()

	a := NewAction(data, s.actionRunAs, s)
	s.actionRunAs = a.Name()
	a.Run()
}

func (s *Session) send(recipientId int, msgs ...string) {

	for i := range msgs {

		msg := msgs[i]
		/*if strings.Contains(msg, "%s") {
			msg = fmt.Sprintf(msg, recipient.FirstName)
		}*/
		s.bot.SendMessage(recipientId, msg)
	}
	return
}
