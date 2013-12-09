package models

import (
	"github.com/garyburd/redigo/redis"
	"log"
	"time"
)

var (
	pool *redis.Pool
)

func Start(serverAddr string) {
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

// FlushAll takes the address of the Redis you want to flush, plus
// a string that MUST match exactly the string
// "Yes I am sure I want to flush all my Redis"
func FlushAll(serverAddr, areYouSure string) {
	if areYouSure != "Yes I am sure I want to flush all my Redis" {
		log.Fatalf("You tried to FlushAll the Redis instance at %s, but were not sure enough!",
			serverAddr)
		return
	}
	log.Printf("Flushing all of Redis at '%s'", serverAddr)
	c, err := redis.Dial("tcp", serverAddr)
	if err != nil {
		log.Fatal(err)
	}
	c.Do("FLUSHALL")
}
