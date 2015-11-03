package lib

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"time"

	"github.com/garyburd/redigo/redis"
)

var FantailStorageErr = errors.New("Fantail storage is not enabled")
var FantailSaveErr = errors.New("Fantail issue trying to save to storage")

type Storage struct {
	store *redis.Pool
}

func NewStorage() *Storage {
	a := &Storage{}
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

func (a *Storage) Save(userId string, n Note) error {

	serialized, err := json.Marshal(n)
	if err != nil {
		return err
	}
	_, err = a.store.Get().Do("LPUSH", userId, serialized)
	return err
}

func (a *Storage) Get(userId string) (Notes, error) {

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
