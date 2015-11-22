package lib

import "testing"

func testUser(id int) *User {
	rn := NewReminderNote(newMsg("/remind 3 to do stuff"), test_tag)
	n := NewNote(newMsg("/say hi"), test_tag)

	return &User{
		id:       id,
		recent:   Notes{&rn, &n},
		lastChat: n.AddedOn,
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

	r := u1.GetReminders()

	if len(r) != 1 {
		t.Error("there should be one reminder")
	}

}

func TestUser_HelpAskedFor(t *testing.T) {
	u1 := testUser(111)

	r := u1.HelpAskedFor()

	if len(r) != 0 {
		t.Error("there should be NO help notes")
	}

}

func TestUser_AddOrUpdate(t *testing.T) {
	users := Users{}
	u1 := testUser(111)
	users = u1.AddOrUpdate(users)
	u2 := testUser(222)
	users = u2.AddOrUpdate(users)

	if len(users) != 2 {
		t.Error("there should be TWO users but have ", len(users))
	}

	u3 := testUser(333)
	users = u3.AddOrUpdate(users)

	if len(users) != 3 {
		t.Error("there should be THREE users but have ", len(users))
	}

	users = u2.AddOrUpdate(users)

	if len(users) != 3 {
		t.Error("there should still be THREE users but have ", len(users))
	}

}
