package lib

import (
	"encoding/json"
	"log"
	"os"

	"github.com/garyburd/redigo/redis"
)

type api struct {
	store redis.Conn
}

func newApi() *api {
	a := &api{}
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

func (a *api) save(userId string, s said) error {
	serialized, err := json.Marshal(s)
	if err != nil {
		return err
	}
	_, err = a.store.Do("LPUSH", userId, serialized)
	return err
}

func (a *api) get(userId string) (chat, error) {

	count, err := redis.Int(a.store.Do("LLEN", userId))

	if err != nil {
		return nil, err
	}

	items, err := redis.Values(a.store.Do("LRANGE", userId, 0, count))

	if err != nil {
		return nil, err
	}

	var c chat

	for i := range items {
		var s said
		serialized, _ := redis.Bytes(items[i], nil)
		json.Unmarshal(serialized, &s)
		c = append(c, &s)
	}
	return c, nil
}

func (a *api) getCurrentTodos(userId string) (chat, error) {

	all, err := a.get(userId)

	if err != nil {
		return nil, err
	}

	var c chat

	for i := range all {

		if all[i].Todo == true && all[i].TodoComplete.IsZero() {
			c = append(c, all[i])
		}
	}

	return c, nil
}

func (a *api) getCompleteTodos(userId string) (chat, error) {

	all, err := a.get(userId)

	if err != nil {
		return nil, err
	}

	var c chat

	for i := range all {

		if all[i].Todo == true && all[i].TodoComplete.IsZero() == false {
			c = append(c, all[i])
		}
	}

	return c, nil
}
