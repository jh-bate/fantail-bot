package lib

import (
	"encoding/json"
	"log"
	"os"

	"github.com/garyburd/redigo/redis"
)

type Storage struct {
	store redis.Conn
}

func NewStorage() *Storage {
	a := &Storage{}
	redisUrl := os.Getenv("REDIS_URL")

	if redisUrl == "" {
		log.Fatal("$REDIS_URL must be set")
	}

	c, err := redis.DialURL(redisUrl)
	if err != nil {
		log.Fatal(err)
	}
	a.store = c
	return a
}

func (a *Storage) Save(userId string, s Said) error {
	serialized, err := json.Marshal(s)
	if err != nil {
		return err
	}
	_, err = a.store.Do("LPUSH", userId, serialized)
	return err
}

func (a *Storage) Get(userId string) (Chat, error) {

	count, err := redis.Int(a.store.Do("LLEN", userId))

	if err != nil {
		return nil, err
	}

	items, err := redis.Values(a.store.Do("LRANGE", userId, 0, count))

	if err != nil {
		return nil, err
	}

	var c Chat

	for i := range items {
		var s Said
		serialized, _ := redis.Bytes(items[i], nil)
		json.Unmarshal(serialized, &s)
		c = append(c, &s)
	}
	return c, nil
}

func (a *Storage) GetCurrentTodos(userId string) (Chat, error) {

	all, err := a.Get(userId)

	if err != nil {
		return nil, err
	}

	var c Chat

	for i := range all {

		if all[i].Remind == true && all[i].RemindComplete.IsZero() {
			c = append(c, all[i])
		}
	}

	return c, nil
}

func (a *Storage) GetCompleteTodos(userId string) (Chat, error) {

	all, err := a.Get(userId)

	if err != nil {
		return nil, err
	}

	var c Chat

	for i := range all {

		if all[i].Remind == true && all[i].RemindComplete.IsZero() == false {
			c = append(c, all[i])
		}
	}

	return c, nil
}
