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

//todo this need to be moved out
func bytes(v interface{}) []byte {
	switch v := v.(type) {
	case []byte:
		return v
	case string:
		return []byte(v)
	}
	return nil
}

func GetUsers() (Users, error) {
	var all Users
	items, err := userStore.ReadAll(user_store_name)
	if err != nil {
		return nil, err
	}

	for i := range items {
		var u User
		json.Unmarshal(bytes(items[i]), &u)
		all = append(all, &u)
	}

	return all, nil
}
