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
	ReadAll(name string) ([]interface{}, error)
}

type RedisStore struct {
	Pool *redis.Pool
}

func NewRedisStore() *RedisStore {
	return &RedisStore{Pool: newPool()}
}

func newPool() *redis.Pool {

	redisUrl := os.Getenv("REDIS_URL")

	if redisUrl == "" {
		log.Panic("REDIS_URL wasn't set")
		return nil
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
	STORE_TEST_DB = iota
	STORE_PROD_DB
)

//Allows us to set the database we are using
func (a *RedisStore) Set(db int) *RedisStore {
	_, err := a.Pool.Get().Do("Select", db)
	if err != nil {
		log.Panic("Error setting the database", err.Error())
	}
	return a
}

func (a *RedisStore) Save(name string, v interface{}) error {

	json.Marshal(v)

	serialized, err := json.Marshal(v)
	if err != nil {
		return err
	}

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

func (a *RedisStore) ReadAll(name string) ([]interface{}, error) {

	c := a.Pool.Get()

	count, err := redis.Int(c.Do("LLEN", name))

	log.Println("found", count)

	if err != nil {
		return nil, err
	}
	return redis.Values(c.Do("LRANGE", name, 0, count))
}
