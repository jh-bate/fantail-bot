package lib

import (
	"os"
	"testing"
)

func storeSetup(userIds ...string) *Storage {
	os.Setenv("REDIS_URL", "redis://localhost:6379")
	myStore := NewStorage()

	//cleanup first
	for i := range userIds {
		myStore.store.Get().Do("DEL", userIds[i])
	}

	return myStore
}

func testData(id int) *User {
	u := &User{
		id: id,
	}

	u.notes = append(u.notes,
		NewNote(newMsg("/say hi")),
		NewNote(newMsg("/say hello")),
		NewNote(newMsg("answering another question")),
		NewNote(newMsg("/say need help to figure stuff out"), help_tag),
		NewNote(newMsg("/say bye")),
	)

	return u
}

func TestStore_Save_and_Get(t *testing.T) {

	user := testData(123)
	store := storeSetup(string(user.id))

	for i := range user.notes {
		err := store.Save(string(user.id), user.notes[i])
		if err != nil {
			t.Error("Error saving to store during tests", err.Error())
		}
	}

	retreived, err := store.Get(string(user.id))
	if err != nil {
		t.Error("Error getting notes from store during tests", err.Error())
	}

	if len(retreived) != len(user.notes) {
		t.Errorf("expected %d got %d", len(user.notes), len(retreived))
	}

	if len(retreived.FilterBy(help_tag)) != 1 {
		t.Errorf("expected %d got %d help items", 1, len(retreived.FilterBy(help_tag)))
	}
}

func TestStore_Update(t *testing.T) {

	user := testData(999)
	store := storeSetup(string(user.id))

	for i := range user.notes {
		err := store.Save(string(user.id), user.notes[i])
		if err != nil {
			t.Error("Error saving to store during tests", err.Error())
		}
	}

	original := user.notes[0]
	updated := original
	updated.Text = "updated"

	err := store.Update(string(user.id), original, updated)

	if err != nil {
		t.Error("Error updating note during tests", err.Error())
	}

	notes, err := store.Get(string(user.id))

	if notes.MostRecent().Text != updated.Text {
		t.Errorf("expected %s got %s ", updated.Text, notes.MostRecent().Text)
	}

}
