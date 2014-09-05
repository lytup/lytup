package lytup

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

var pool *redis.Pool

func Redis() redis.Conn {
	return pool.Get()
}

func init() {
	cfg := Config.Redis
	pool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				return nil, err
			}

			if cfg.Password != "" {
				if _, err := c.Do("AUTH", cfg.Password); err != nil {
					c.Close()
					return nil, err
				}
			}

			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
