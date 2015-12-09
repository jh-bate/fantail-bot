package lib

import (
	"os"
	"testing"
)

func storeSetup(userIds ...string) Store {
	const test_db = 1
	os.Setenv("REDIS_URL", "redis://localhost:6379")
	myStore := NewRedisStore().setDb(STORE_TEST_DB)

	//cleanup first
	myStore.pool.Get().Do("FLUSHDB")

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

func TestStore_SaveNote_and_GetNotes(t *testing.T) {

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

func TestStore_UpdateNote(t *testing.T) {

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

func TestStore_SaveUser_and_GetUsers(t *testing.T) {

	user := testData(123)
	store := storeSetup(string(user.Id))

	//init positive and followup
	user.IsPostive(1)
	user.FollowUpAbout()

	err := store.SaveUser(user)
	if err != nil {
		t.Error("Error saving the user", err.Error())
	}

	usrs, err := store.GetUsers()

	if err != nil {
		t.Fatal("Error getting users", err.Error())
	}

	if len(usrs) != 1 {
		t.Error("there should have been one user returned but got ", len(usrs))
	}

	if usrs[0].Id != user.Id {
		t.Errorf("expected %d got %d", user.Id, usrs[0].Id)
	}

	if len(usrs[0].Learnings) != 1 {
		t.Error("there should have been one learning ", len(usrs[0].Learnings))
	}

	if len(usrs[0].Helped) != 1 {
		t.Error("there should have been one help case ", len(usrs[0].Helped))
	}

}

func TestStore_UpdateUser(t *testing.T) {

	user := testData(999)
	store := storeSetup(string(user.Id))

	//init positive and followup
	user.IsPostive(1)
	user.FollowUpAbout()

	err := store.SaveUser(user)
	if err != nil {
		t.Error("Error saving the user", err.Error())
	}

	//do some updates
	updates := user
	updates.IsPostive(1)
	updates.FollowUpAbout()

	err = store.UpdateUser(user, updates)

	if err != nil {
		t.Error("Error saving the user", err.Error())
	}

	usrs, err := store.GetUsers()

	if err != nil {
		t.Fatal("Error getting users", err.Error())
	}

	if len(usrs) != 1 {
		t.Error("there should have been one user returned but got ", len(usrs))
	}

	if usrs[0].Id != user.Id {
		t.Errorf("expected %d got %d", user.Id, usrs[0].Id)
	}

	if len(usrs[0].Learnings) != 2 {
		t.Error("there should have been two learnings ", len(usrs[0].Learnings))
	}

	if len(usrs[0].Helped) != 2 {
		t.Error("there should have been two help case ", len(usrs[0].Helped))
	}

}
