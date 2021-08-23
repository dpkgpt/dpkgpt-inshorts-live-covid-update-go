package config

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

var RedisPool *redis.Pool

func InitRedisConfig() {
	RedisPool = &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}
}
