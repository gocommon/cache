package redis

import (
	"time"

	redigo "github.com/garyburd/redigo/redis"
)

// Store Store
type Store struct {
	pool *redigo.Pool
}

// NewStore NewStore
func NewStore(opts ...Options) *Store {
	options := Options{}

	if len(opts) > 0 {
		options = opts[0]
	}
	options = defaultOptions(options)

	return &Store{
		pool: newRedisPool(options),
	}
}

func (s *Store) do(commandName string, args ...interface{}) (reply interface{}, err error) {
	c := s.pool.Get()
	defer c.Close()
	return c.Do(commandName, args...)
}

// Set Set
func (s *Store) Set(key string, val string, ttl int64) error {
	_, err := s.do("SETEX", key, ttl, val)
	if err != nil {
		return err
	}

	return nil
}

// Get Get
func (s *Store) Get(key string) (string, error) {
	ret, err := redigo.String(s.do("GET", key))
	if err != nil {
		if err == redigo.ErrNil {
			return "", nil
		}
		return "", err
	}

	return ret, nil
}

// Expire Expire
func (s *Store) Expire(key string, ttl int64) error {
	_, err := s.do("EXPIRE", key, ttl)
	if err != nil {
		return err
	}

	return nil
}

// Forever Forever
func (s *Store) Forever(key string, val string) error {
	_, err := s.do("SET", key, val)
	if err != nil {
		return err
	}

	return nil
}

// Del Del
func (s *Store) Del(key string) error {
	_, err := s.do("DEL", key)
	if err != nil {
		return err
	}

	return nil

}

func newRedisPool(conf Options) *redigo.Pool {
	return &redigo.Pool{
		MaxIdle:     conf.MaxIdle,
		IdleTimeout: time.Duration(conf.IdleTimeout) * time.Second,
		Dial: func() (redigo.Conn, error) {
			c, err := redigo.Dial("tcp", conf.Addr)
			if err != nil {
				return nil, err
			}
			if len(conf.Passwd) > 0 {
				if _, err := c.Do("AUTH", conf.Passwd); err != nil {
					c.Close()
					return nil, err
				}
			}
			_, err = c.Do("SELECT", conf.SelectDB)
			if err != nil {
				c.Close()
				return nil, err
			}

			return c, err
		},
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
