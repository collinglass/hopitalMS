package models

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

var (
	serverAddr = "192.168.0.24:6379"
	pool       *redis.Pool
)

func init() {
	pool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", serverAddr)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
