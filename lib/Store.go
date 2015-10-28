package lib

import (
	"encoding/json"
	"errors"
	"log"
	"os"

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

func (a *Storage) Save(userId string, n Note) error {
	serialized, err := json.Marshal(n)
	if err != nil {
		return err
	}
	_, err = a.store.Do("LPUSH", userId, serialized)
	return err
}

func (a *Storage) Get(userId string) (Notes, error) {

	count, err := redis.Int(a.store.Do("LLEN", userId))

	if err != nil {
		return nil, err
	}

	items, err := redis.Values(a.store.Do("LRANGE", userId, 0, count))

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
