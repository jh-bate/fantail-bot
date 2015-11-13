package lib

import "github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"

const (
	default_script = "default"
	stickers_chat  = "chat"
)

type (
	Info struct {
		App       []string `json:"appInfo"`
		Reminders []string `json:"remindersInfo"`
		Chat      []string `json:"chatInfo"`
		Said      []string `json:"saidInfo"`
	}

	QProcess struct {
		s *session
	}
)

func NewQProcess(b *telebot.Bot, s *Storage) *QProcess {
	q := &QProcess{s: newSession(b, s)}
	return q
}

func (this *QProcess) Run(input <-chan telebot.Message) {

	prevActionName := ""

	for msg := range input {
		in := newIncoming(msg)

		in.getAction(this.s, prevActionName).firstUp().askQuestion()
		prevActionName = in.action.getName()
	}
}
