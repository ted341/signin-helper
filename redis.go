package main

import (
	"os"
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
	//"github.com/nitishm/go-rejson"
)

type Redis struct {
	client  *redis.Client
	//handler *rejson.Handler
}

var __redis Redis

const (
	HOUR   = time.Hour
	MINUTE = time.Minute
	SECOND = time.Second
)

func initRedis() {

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)
	redisOpts := &redis.Options{
		Addr: redisAddr,
	}

	__redis = Redis{
		client: redis.NewClient(redisOpts),
		//handler: rejson.NewReJSONHandler(),
	}

	return
}

func (r Redis) Get(key string) (val string, found bool, err error) {
	val, err = r.client.Get(key).Result()
	if err == redis.Nil {
		err = nil
		return
	}
	found = true
	return
}
/*
func (r Redis) GetJSON(key, path string) (res interface{}, err error) {
	res, err = r.handler.JSONGet(key, path)
	return
}
*/
func (r Redis) Set(key string, val interface{}, expr time.Duration) (err error) {
	err = r.client.Set(key, val, expr).Err()
	return
}
/*
func (r Redis) SetJSON(key, path string, value interface{}) (res interface{}, err error) {
	res, err = r.handler.JSONSet(key, path, value)
	return
}
*/
func (r Redis) Del(key string) (err error) {
	err = r.client.Del(key).Err()
	return
}

func (r Redis) Incr(key string) (result int64, err error) {
	result, err = r.client.Incr(key).Result()
	if err != nil {
		panic(err)
	}
	return
}

func (r Redis) SAdd(key string, members ...interface{}) (err error) {
	err = r.client.SAdd(key, members...).Err()
	return
}

func (r Redis) SMembers(key string) (list []string, err error) {
	list, err = r.client.SMembers(key).Result()
	return
}