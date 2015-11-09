package lib

import (
	"strings"
	"testing"
)

const (
	test_tag = "TODO"
)

func TestNote_NewNote(t *testing.T) {

	hi := newMsg("/say hi")
	n := NewNote(hi, test_tag)

	if n.IsEmpty() {
		t.Error("should not be empty")
	}

	if !strings.Contains(n.Tag, "/say") {
		t.Errorf("be tagged with cmd name Tag=%s", n.Tag)
	}

	if !strings.Contains(n.Tag, test_tag) {
		t.Errorf("be tagged with passed tags Tag=%s", n.Tag)
	}

	if n.Text != "hi" {
		t.Error("but is ", n.Text)
	}

}

func TestNote_NewReminderNote(t *testing.T) {

	hi := newMsg("/remind 3 to do stuff")
	n := NewReminderNote(hi, test_tag)

	if n.IsEmpty() {
		t.Error("should not be empty")
	}

	if !strings.Contains(n.Tag, remind_tag) {
		t.Errorf("be tagged with remind_tag Tag=%s", n.Tag)
	}

	if !strings.Contains(n.Tag, test_tag) {
		t.Errorf("be tagged with passed tags Tag=%s", n.Tag)
	}

	if n.Text != "to do stuff" {
		t.Error("but is ", n.Text)
	}

}
