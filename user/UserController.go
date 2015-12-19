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

func (this *User) Save() error {
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
