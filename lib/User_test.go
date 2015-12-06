package lib

import (
	"testing"
	"time"
)

func testUser(id int) *User {
	n1 := NewNote(newMsg("to do stuff"), test_tag)
	n2 := NewNote(newMsg("/say hi"), help_tag)

	return &User{
		id:    id,
		notes: Notes{n1, n2},
	}
}

func TestUser_ToBotUser(t *testing.T) {

	u1 := testUser(111)

	bu1 := u1.ToBotUser()

	if bu1.ID != u1.id {
		t.Error("Ids should match")
	}
}

func TestUser_FollowUp(t *testing.T) {
	u1 := testUser(111)

	r := u1.FollowUpAbout()

	if len(r) != 1 {
		t.Error("there should be ONE help note", len(r))
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

	u3.notes = Notes{&Note{Added: time.Now().Add(3)}}

	users = u3.AddOrUpdate(users)

	if len(users) != 3 {
		t.Error("there should be THREE users but have ", len(users))
	}

	if users[2].notes[0].Added.YearDay() != u3.notes[0].Added.YearDay() {
		t.Errorf("expetced [%s] found [%s]", u3.notes[0].Added.String(), users[2].notes[0].Added.String())
	}

	if len(users[2].notes) != 1 {
		t.Errorf("expected one found [%d] recent notes", len(users[2].notes))
	}

}
