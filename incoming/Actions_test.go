package incoming

import (
	"testing"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
)

func makeTestAction(msg telebot.Message, tag string) Action {
	s := newSession(nil, nil)
	s.in = newIncoming(msg)
	return NewAction(s, tag)
}

func TestActions_getName(t *testing.T) {

	action := makeTestAction(newMsg("/say hi"), "")

	if action.getName() != "/say" {
		t.Errorf("expected %v got %v", "/say", action.getName())
	}

}

func TestActions_getName_usesPrevActionName(t *testing.T) {

	action := makeTestAction(newMsg("some more chatting"), chat_action)

	if action.getName() != chat_action {
		t.Errorf("expected %v got %v", chat_action, action.getName())
	}

}

func TestActions_findNext(t *testing.T) {

	action := makeTestAction(newMsg("/say hi"), "")

	q := action.nextQuestion()

	if q == nil {
		t.Errorf("expected nil got %v", q)
	}

}

func TestActions_getQuestions(t *testing.T) {

	action := makeTestAction(newMsg("/say hi"), "")

	q := action.getQuestions()

	if q == nil {
		t.Errorf("expected nil got %v", q)
	}

	if len(q) < 1 {
		t.Errorf("expected questions got to be loaded but got %d", len(q))
	}

}
