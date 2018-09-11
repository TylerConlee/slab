package datastore

import (
	"github.com/go-redis/redis"
	l "github.com/tylerconlee/slab/log"
)

var log = l.Log
var client *redis.Client

// RedisConnect establishes a connection to the localhost Redis instance.
func RedisConnect(db int) {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       db, // use default DB
	})

	pong, err := client.Ping().Result()
	if err != nil {
		log.Error("Error encountered attempting to connect to Redis.", map[string]interface{}{
			"error": err,
		})
	}
	log.Info("Redis connected at localhost:6379.", map[string]interface{}{
		"result": pong,
	})

}

// Save takes a key and value pair and saves it to the Redis instance.
func Save(key string, value string) (result bool) {

	err := client.Set(key, value, 0).Err()
	if err != nil {
		log.Error("Error attempting to save to Redis.", map[string]interface{}{
			"client": client,
			"error":  err,
		})
		return false
	}
	return true
}

// Load takes a key and returns the result of the lookup in Redis.
func Load(key string) (result string) {
	val, err := client.Get(key).Result()
	if err != nil {
		log.Error("Error attempting to load from Redis.", map[string]interface{}{
			"client": client,
			"error":  err,
		})
		return
	}
	return val
}
