package user

import (
	"encoding/json"

	"github.com/jh-bate/fantail-bot/store"
)

var userStore store.Store

const user_store_name = "users"

func init() {
	userStore = store.NewRedisStore()
	return
}

func (this *User) Upsert() error {

	curr, err := GetUser(this.Id)
	if err != nil {
		return err
	}
	if curr != nil {
		//learnings
		for key, value := range curr.Learnt {
			if _, ok := this.Learnt[key]; !ok {
				this.Learnt[key] = value
			}
		}
		//helped
		for key, value := range curr.Helped {
			if _, ok := this.Helped[key]; !ok {
				this.Helped[key] = value
			}
		}
	}
	return userStore.Save(user_store_name, this)
}

func (this *User) Delete() error {
	return userStore.Delete(user_store_name, this)
}

func GetUsers() (Users, error) {
	var all Users
	items, err := userStore.ReadAll(user_store_name)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(items, &all)

	return all, nil
}

func GetUser(id int) (*User, error) {

	all, err := GetUsers()
	if err != nil {
		return nil, err
	}
	for i := range all {
		if all[i].Id == id {
			return all[i], nil
		}
	}

	return nil, nil
}
