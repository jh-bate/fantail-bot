package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/garyburd/redigo/redis"
)

var FantailStorageErr = errors.New("Fantail storage is not enabled")
var FantailSaveErr = errors.New("Fantail issue trying to save to storage")

type Store interface {
	SaveNote(userId string, note *Note) error
	UpdateNote(userId string, original, updated *Note) error
	GetNotes(userId string) (Notes, error)
	SaveUser(u *User) error
	UpdateUser(original, updated *User) error
	GetUsers() (Users, error)
}

type RedisStore struct {
	store *redis.Pool
}

func NewRedisStore() *RedisStore {
	a := &RedisStore{}
	redisUrl := os.Getenv("REDIS_URL")

	if redisUrl == "" {
		log.Fatal("$REDIS_URL must be set")
	}

	a.store = newPool()
	return a
}

func newPool() *redis.Pool {

	redisUrl := os.Getenv("REDIS_URL")

	if redisUrl == "" {
		log.Fatal("$REDIS_URL must be set")
	}
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(redisUrl)
			if err != nil {
				return nil, err
			}
			return c, err
		},
	}
}

func (a *RedisStore) SaveNote(userId string, note *Note) error {

	serialized, err := json.Marshal(note)
	if err != nil {
		return err
	}
	_, err = a.store.Get().Do("LPUSH", userId, serialized)
	return err
}

func (a *RedisStore) UpdateNote(userId string, original, updated *Note) error {

	serializedOriginal, err := json.Marshal(original)
	if err != nil {
		return err
	}
	_, err = a.store.Get().Do("LREM", userId, -1, serializedOriginal)
	if err != nil {
		return err
	}
	return a.SaveNote(userId, updated)
}

func (a *RedisStore) GetNotes(userId string) (Notes, error) {

	c := a.store.Get()

	count, err := redis.Int(c.Do("LLEN", userId))

	if err != nil {
		return nil, err
	}

	items, err := redis.Values(c.Do("LRANGE", userId, 0, count))

	if err != nil {
		return nil, err
	}

	var all Notes

	for i := range items {
		var n Note
		serialized, _ := redis.Bytes(items[i], nil)
		json.Unmarshal(serialized, &n)
		all = append(all, &n)
	}
	return all, nil
}

/*
func (a *RedisStore) GetUsers() ([]int, error) {
	c := a.store.Get()
	return redis.Ints(c.Do("KEYS", "*"))
}
*/

func userKey(query string) string {
	const user_pattern = "user_%s"
	return fmt.Sprintf(user_pattern, query)
}

func (a *RedisStore) GetUsers() (Users, error) {
	c := a.store.Get()

	keys, err := redis.String(c.Do("KEYS", userKey("*")))

	if err != nil {
		return nil, err
	}

	var all Users

	for i := range keys {

		items, err := redis.Values(c.Do("LRANGE", keys[i], 0, 0))

		if err != nil {
			return nil, err
		}

		for i := range items {
			var u User
			serialized, _ := redis.Bytes(items[i], nil)
			json.Unmarshal(serialized, &u)
			all = append(all, &u)
			break
		}
	}
	return all, nil
}

func (a *RedisStore) SaveUser(u *User) error {

	serialized, err := json.Marshal(u)
	if err != nil {
		return err
	}
	_, err = a.store.Get().Do("LPUSH", userKey(string(u.Id)), serialized)
	return err

}
func (a *RedisStore) UpdateUser(original, updated *User) error {
	serializedOriginal, err := json.Marshal(original)
	if err != nil {
		return err
	}
	_, err = a.store.Get().Do("LREM", userKey(string(original.Id)), -1, serializedOriginal)
	if err != nil {
		return err
	}
	return a.SaveUser(updated)
}
