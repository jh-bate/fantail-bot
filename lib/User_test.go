package lib

import (
	"testing"
	"time"
)

func testUser(id int) *User {
	rn := NewReminderNote(newMsg("/remind 3 to do stuff"), test_tag)
	n := NewNote(newMsg("/say hi"), test_tag)

	return &User{
		id:     id,
		recent: Notes{rn, n},
	}
}

func TestUser_ToBotUser(t *testing.T) {

	u1 := testUser(111)

	bu1 := u1.ToBotUser()

	if bu1.ID != u1.id {
		t.Error("Ids should match")
	}
}

func TestUser_GetReminders(t *testing.T) {
	u1 := testUser(111)

	r := u1.Reminders()

	if len(r) != 1 {
		t.Error("there should be one reminder")
	}

}

func TestUser_HelpAskedFor(t *testing.T) {
	u1 := testUser(111)

	r := u1.HelpWanted()

	if len(r) != 0 {
		t.Error("there should be NO help notes")
	}

}

func TestUser_AddOrUpdate(t *testing.T) {
	users := Users{}
	u1 := testUser(111)
	users = u1.AddOrUpdate(users)

	if len(users) != 1 {
		t.Error("there should be ONE users but have ", len(users))
	}

	if users[0].id != u1.id {
		t.Errorf("expected [%d] found [%d]", u1.id, users[0].id)
	}

	u2 := testUser(222)
	users = u2.AddOrUpdate(users)

	if len(users) != 2 {
		t.Error("there should be TWO users but have ", len(users))
	}

	if users[1].id != u2.id {
		t.Errorf("expected [%d] found [%d]", u2.id, users[1].id)
	}

	u3 := testUser(333)
	users = u3.AddOrUpdate(users)

	if len(users) != 3 {
		t.Error("there should be THREE users but have ", len(users))
	}

	if users[2].id != u3.id {
		t.Errorf("expected [%d] found [%d]", u3.id, users[2].id)
	}

	users = u2.AddOrUpdate(users)

	if len(users) != 3 {
		t.Error("there should still be THREE users but have ", len(users))
	}

}

func TestUser_AddOrUpdate_withUpdate(t *testing.T) {
	users := Users{}
	u1 := testUser(111)
	users = u1.AddOrUpdate(users)
	u2 := testUser(222)
	users = u2.AddOrUpdate(users)
	u3 := testUser(333)
	users = u3.AddOrUpdate(users)

	if len(users) != 3 {
		t.Error("there should be THREE users but have ", len(users))
	}

	u3.recent = Notes{&Note{AddedOn: time.Now().Add(3)}}

	users = u3.AddOrUpdate(users)

	if len(users) != 3 {
		t.Error("there should be TWO users but have ", len(users))
	}

	if users[2].LastChatted().YearDay() != u3.LastChatted().YearDay() {
		t.Errorf("expetced [%s] found [%s]", u3.LastChatted().String(), users[2].LastChatted().String())
	}

	if len(users[2].recent) != 1 {
		t.Errorf("expected one found [%d] recent notes", len(users[2].recent))
	}

}
