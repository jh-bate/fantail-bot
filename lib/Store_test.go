package lib

import (
	"os"
	"testing"
)

func storeSetup(userIds ...string) Store {
	os.Setenv("REDIS_URL", "redis://localhost:6379")
	myStore := NewRedisStore()

	//cleanup first
	for i := range userIds {
		myStore.store.Get().Do("DEL", userIds[i])
	}

	return myStore
}

func testData(id int) *User {
	u := &User{
		Id: id,
	}

	u.Notes = append(u.Notes,
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
	store := storeSetup(string(user.Id))

	for i := range user.Notes {
		err := store.SaveNote(string(user.Id), user.Notes[i])
		if err != nil {
			t.Error("Error saving to store during tests", err.Error())
		}
	}

	retreived, err := store.GetNotes(string(user.Id))
	if err != nil {
		t.Error("Error getting notes from store during tests", err.Error())
	}

	if len(retreived) != len(user.Notes) {
		t.Errorf("expected %d got %d", len(user.Notes), len(retreived))
	}

	if len(retreived.FilterOnTag(help_tag)) != 1 {
		t.Errorf("expected %d got %d help items", 1, len(retreived.FilterOnTag(help_tag)))
	}
}

func TestStore_Update(t *testing.T) {

	user := testData(999)
	store := storeSetup(string(user.Id))

	for i := range user.Notes {
		err := store.SaveNote(string(user.Id), user.Notes[i])
		if err != nil {
			t.Error("Error saving to store during tests", err.Error())
		}
	}

	original := user.Notes[0]
	updated := original
	updated.Text = "updated"

	err := store.UpdateNote(string(user.Id), original, updated)

	if err != nil {
		t.Error("Error updating note during tests", err.Error())
	}

	notes, err := store.GetNotes(string(user.Id))

	if notes.MostRecent().Text != updated.Text {
		t.Errorf("expected %s got %s ", updated.Text, notes.MostRecent().Text)
	}

}
