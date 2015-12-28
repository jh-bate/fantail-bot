package store

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/garyburd/redigo/redis"
)

type Store interface {
	Save(name string, v interface{}) error
	Delete(name string, v interface{}) error
	ReadAll(name string) ([]byte, error)
}

type RedisStore struct {
	Pool *redis.Pool
}

func NewRedisStore() *RedisStore {
	store := &RedisStore{Pool: newPool()}
	//default to this ...
	store.Set(STORE_PROD_DB)
	return store
}

func newPool() *redis.Pool {

	redisUrl := os.Getenv("REDIS_URL")

	if redisUrl == "" {
		log.Println("REDIS_URL wasn't set, running as localhost")
		redisUrl = "redis://localhost:6379"
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

const (
	STORE_TEST_DB = 0
	STORE_PROD_DB = 1
)

//Allows us to set the database we are using
func (a *RedisStore) Set(db int) *RedisStore {
	log.Println("setting db as", db)
	_, err := a.Pool.Get().Do("Select", db)
	if err != nil {
		log.Panic("Error setting the database", err.Error())
	}
	return a
}

func (a *RedisStore) Save(name string, v interface{}) error {

	serialized, err := json.Marshal(v)
	if err != nil {
		return err
	}

	log.Println("saving...", string(serialized), "into", name)

	_, err = a.Pool.Get().Do("LPUSH", name, serialized)
	return err
}

func (a *RedisStore) Delete(name string, v interface{}) error {

	serialized, err := json.Marshal(v)
	if err != nil {
		return err
	}
	_, err = a.Pool.Get().Do("LREM", name, -1, serialized)
	return err
}

func bytes(v interface{}) []byte {
	switch v := v.(type) {
	case []byte:
		return v
	case string:
		return []byte(v)
	}
	return nil
}

func (a *RedisStore) ReadAll(name string) ([]byte, error) {

	c := a.Pool.Get()

	count, err := redis.Int(c.Do("LLEN", name))

	if err != nil {
		return nil, err
	}
	vals, err := redis.Values(c.Do("LRANGE", name, 0, count))
	if err != nil {
		return nil, err
	}
	data := []byte("[")
	for i := range vals {
		if string(data) != "[" {
			data = append(data, []byte(",")...)
		}
		data = append(data, bytes(vals[i])...)
	}
	data = append(data, []byte("]")...)
	return data, nil
}
