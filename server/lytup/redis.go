package lytup

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

var pool *redis.Pool

func Redis() redis.Conn {
	return pool.Get()
}

// SetKeyWithExpiry sets a key with expiry.
func SetKeyWithExpiry(key, val string, expiry uint) error {
	r := Redis()
	defer r.Close()
	r.Send("MULTI")
	r.Send("SET", key, val)
	r.Send("EXPIRE", key, C.VerifyEmailExpiry)
	if _, err := r.Do("EXEC"); err != nil {
		return err
	}
	return nil
}

func init() {
	pool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				return nil, err
			}

			if C.Redis.Password != "" {
				if _, err := c.Do("AUTH", C.Redis.Password); err != nil {
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
