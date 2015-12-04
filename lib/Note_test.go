package lib

import (
	"strings"
	"testing"
	"time"
)

const (
	test_tag = "TODO"
)

func TestNote_NewNote(t *testing.T) {

	hi := newMsg("/say hi")
	n := NewNote(hi, "", test_tag)

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

func TestNote_IsEmpty(t *testing.T) {

	n := Note{}

	if n.IsEmpty() == false {
		t.Error("should be empty")
	}

}

func TestNote_Complete(t *testing.T) {

	n := NewNote(newMsg("/say hi"), "", test_tag)

	if !n.Completed.IsZero() {
		t.Error("note should not be complete on creation")
	}

	n.Complete()

	if n.Completed.IsZero() {
		t.Error("note should now be complete")
	}
}

func TestNote_Update(t *testing.T) {

	n := NewNote(newMsg("/say hi"), "", test_tag)

	if !n.Updated.IsZero() {
		t.Error("note should not be updated on creation")
	}

	n.Update()

	if n.Updated.IsZero() {
		t.Error("note should now be updated")
	}
}

func TestNote_tagFromMsg(t *testing.T) {

	empty := tagFromMsg("")

	if empty != "" {
		t.Errorf("should return no tage but got %s", empty)
	}

	hi := newMsg("/hi here I am")

	hiTag := tagFromMsg(hi.Text)

	if hiTag != "/hi" {
		t.Errorf("expected %s got %s", "/hi", hiTag)
	}

	stuff := newMsg("/stuff")

	stuffTag := tagFromMsg(stuff.Text)

	if stuffTag != "/stuff" {
		t.Errorf("expected %s got %s", "/stuff", stuffTag)
	}

}

func TestNote_Notes_SortByDate(t *testing.T) {

	n := NewNote(newMsg("/say stuff"), "", test_tag)
	n.Added = time.Now().Add(-5)
	n2 := NewNote(newMsg("/say moar"), "", test_tag)
	n2.Added = time.Now()
	n3 := NewNote(newMsg("/say moar stuff"), "", test_tag)
	n3.Added = time.Now().Add(5)

	notes := Notes{n, n2, n3}

	if notes.SortByDate()[0].Added.YearDay() != n3.Added.YearDay() {
		t.Errorf("expected %d got %d", n3.Added.YearDay(), notes.SortByDate()[0].Added.YearDay())
	}

	if notes.SortByDate()[1].Added.YearDay() != n2.Added.YearDay() {
		t.Errorf("expected %d got %d", n2.Added.YearDay(), notes.SortByDate()[1].Added.YearDay())
	}

	if notes.SortByDate()[2].Added.YearDay() != n.Added.YearDay() {
		t.Errorf("expected %d got %d", n.Added.YearDay(), notes.SortByDate()[2].Added.YearDay())
	}

}

func TestNote_Notes_MostRecent(t *testing.T) {

	n := NewNote(newMsg("/say stuff"), "", test_tag)
	n.Added = time.Now().Add(-5)
	n2 := NewNote(newMsg("/say moar"), "", test_tag)
	n2.Added = time.Now()
	n3 := NewNote(newMsg("/say moar stuff"), "", test_tag)
	n3.Added = time.Now().Add(5)

	notes := Notes{n, n2, n3}

	latest := notes.MostRecent()

	if latest.Added.YearDay() != n3.Added.YearDay() {
		t.Errorf("expected %d got %d", n3.Added.YearDay(), latest.Added.YearDay())
	}

}

func TestNote_Notes_MostRecent_whenEmpty(t *testing.T) {

	notes := Notes{}

	latest := notes.MostRecent()

	if !latest.Added.IsZero() {
		t.Error("there are no notes so the most recent should be empty")
	}

}
