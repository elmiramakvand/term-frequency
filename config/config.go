package config

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

var (
	pool *redis.Pool
)

func GetRedisPool(addr string, dbNo int) *redis.Pool {
	pool := &redis.Pool{
		MaxActive:   12000,
		MaxIdle:     80,
		IdleTimeout: 240 * time.Second,
		// Dial or DialContext must be set. When both are set, DialContext takes precedence over Dial.
		Dial: func() (redis.Conn, error) { return redis.Dial("tcp", addr, redis.DialDatabase(dbNo)) },
	}
	return pool
}
