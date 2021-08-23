package config

import (
	"crud/env"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

var RedisPool *redis.Pool

func InitRedisConfig() {
	RedisPool = &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", env.GetValue("REDIS_URL"))
		},
	}
	conn := RedisPool.Get()
	defer conn.Close()
	_, err := conn.Do("AUTH", env.GetValue("REDIS_USER"), env.GetValue("REDIS_PWD"))
	if err != nil {
		log.Fatal("Auth to the Redis database failed", err)
	}
	// Test the connection
	_, err = conn.Do("PING")
	if err != nil {
		log.Fatal("Can't connect to the Redis database", err)
	}
}
