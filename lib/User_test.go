package lib

import (
	"testing"
	"time"
)

func testUser(Id int) *User {
	n1 := NewNote(newMsg("to do stuff"), test_tag)
	n2 := NewNote(newMsg("/say hi"), help_tag)

	return &User{
		Id:    Id,
		Notes: Notes{n1, n2},
	}
}

func TestUser_ToBotUser(t *testing.T) {

	u1 := testUser(111)

	bu1 := u1.ToBotUser()

	if bu1.ID != u1.Id {
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

func TestUser_FollowUp_Helped(t *testing.T) {
	u1 := testUser(111)

	u1.FollowUpAbout()

	if len(u1.Helped) != 1 {
		t.Error("helped should have been set")
	}

	if u1.Helped[0].IsZero() {
		t.Error("the helped date should have been set")
	}
}

func TestUser_IsPostive_Learnings(t *testing.T) {
	u1 := testUser(111)

	pos := u1.IsPostive(1)

	if len(u1.Learnings) != 1 {
		t.Error("should have learnt something")
	}

	if u1.Learnings[0].Date.IsZero() {
		t.Error("the learnings date should have been set")
	}

	if u1.Learnings[0].Positive != pos {
		t.Error("the learnings should have been same as what was returned from IsPostive")
	}
}

func TestUser_AddOrUpdate(t *testing.T) {
	users := Users{}
	u1 := testUser(111)
	users = u1.AddOrUpdate(users)

	if len(users) != 1 {
		t.Error("there should be ONE users but have ", len(users))
	}

	if users[0].Id != u1.Id {
		t.Errorf("expected [%d] found [%d]", u1.Id, users[0].Id)
	}

	u2 := testUser(222)
	users = u2.AddOrUpdate(users)

	if len(users) != 2 {
		t.Error("there should be TWO users but have ", len(users))
	}

	if users[1].Id != u2.Id {
		t.Errorf("expected [%d] found [%d]", u2.Id, users[1].Id)
	}

	u3 := testUser(333)
	users = u3.AddOrUpdate(users)

	if len(users) != 3 {
		t.Error("there should be THREE users but have ", len(users))
	}

	if users[2].Id != u3.Id {
		t.Errorf("expected [%d] found [%d]", u3.Id, users[2].Id)
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

	u3.Notes = Notes{&Note{Added: time.Now().Add(3)}}

	users = u3.AddOrUpdate(users)

	if len(users) != 3 {
		t.Error("there should be THREE users but have ", len(users))
	}

	if users[2].Notes[0].Added.YearDay() != u3.Notes[0].Added.YearDay() {
		t.Errorf("expetced [%s] found [%s]", u3.Notes[0].Added.String(), users[2].Notes[0].Added.String())
	}

	if len(users[2].Notes) != 1 {
		t.Errorf("expected one found [%d] recent notes", len(users[2].Notes))
	}

}
