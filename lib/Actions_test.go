package lib

import "testing"

func TestActions_getName(t *testing.T) {

	action := NewAction(newIncoming(newMsg("/say hi")), nil, "")

	if action.getName() != "/say" {
		t.Errorf("expected %v got %v", "/say", action.getName())
	}

}

func TestActions_getName_usesPrevActionName(t *testing.T) {

	action := NewAction(newIncoming(newMsg("some more chatting")), nil, chat_action)

	if action.getName() != chat_action {
		t.Errorf("expected %v got %v", chat_action, action.getName())
	}

}

func TestActions_findNext(t *testing.T) {

	action := NewAction(newIncoming(newMsg("/say hi")), nil, "")

	q := action.nextQuestion()

	if q == nil {
		t.Errorf("expected nil got %v", q)
	}

}

func TestActions_getQuestions(t *testing.T) {

	action := NewAction(newIncoming(newMsg("/say hi")), nil, "")

	q := action.getQuestions()

	if q == nil {
		t.Errorf("expected nil got %v", q)
	}

	if len(q) < 1 {
		t.Errorf("expected questions got to be loaded but got %d", len(q))
	}

}
