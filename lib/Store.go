package lib

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"strings"

	"github.com/garyburd/redigo/redis"
)

var StorageInitErr = errors.New("Storage is not enabled")
var StorageSaveErr = errors.New("Error trying to save to storage")

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

func (a *Storage) Save(userId string, r Reminder) error {
	serialized, err := json.Marshal(r)
	if err != nil {
		return err
	}
	_, err = a.store.Do("LPUSH", userId, serialized)
	return err
}

func (a *Storage) Get(userId string) (Reminders, error) {

	count, err := redis.Int(a.store.Do("LLEN", userId))

	if err != nil {
		return nil, err
	}

	items, err := redis.Values(a.store.Do("LRANGE", userId, 0, count))

	if err != nil {
		return nil, err
	}

	var all Reminders

	for i := range items {
		var r Reminder
		serialized, _ := redis.Bytes(items[i], nil)
		json.Unmarshal(serialized, &r)
		all = append(all, &r)
	}
	return all, nil
}

func (a *Storage) GetReminders(userId string) (Reminders, error) {

	all, err := a.Get(userId)

	if err != nil {
		return nil, err
	}

	var curr Reminders

	for i := range all {
		if all[i].CompletedOn.IsZero() && strings.Contains(all[i].Tag, remind_cmd) {
			curr = append(curr, all[i])
		}
	}

	return curr, nil
}

func (a *Storage) GetCurrentTodos(userId string) (Reminders, error) {

	all, err := a.Get(userId)

	if err != nil {
		return nil, err
	}

	var curr Reminders

	for i := range all {

		if all[i].CompletedOn.IsZero() == true {
			curr = append(curr, all[i])
		}
	}

	return curr, nil
}

func (a *Storage) GetCompleteTodos(userId string) (Reminders, error) {

	all, err := a.Get(userId)

	if err != nil {
		return nil, err
	}

	var complete Reminders

	for i := range all {

		if all[i].CompletedOn.IsZero() != true {
			complete = append(complete, all[i])
		}
	}

	return complete, nil
}
