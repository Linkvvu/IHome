package dao

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

var RedisConnPool *redis.Pool

func init() {
	RedisConnPool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", ":6379")
		},

		MaxIdle:     5,
		MaxActive:   10,
		IdleTimeout: time.Minute * 3,
	}
}
