package user

import (
	"encoding/json"
	"fmt"

	"github.com/jh-bate/fantail-bot/store"
)

var userStore store.Store

const user_store_name = "user_"

func init() {
	userStore = store.NewRedisStore()
	return
}

func nameForStore(id int) string {
	return fmt.Sprintf("%s%d", user_store_name, id)
}

func (this *User) Save() error {
	return userStore.Save(nameForStore(this.Id), this)
}

func (this *User) Delete() error {
	return userStore.Delete(nameForStore(this.Id), this)
}

func GetUsers() (Users, error) {
	var all Users
	jsonD, err := userStore.ReadAll(user_store_name + "*")
	if err != nil {
		return nil, err
	}
	json.Unmarshal(jsonD, &all)
	return all, nil
}
