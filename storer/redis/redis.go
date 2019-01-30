package storer

import (
	"encoding/json"
	"time"

	"github.com/gocommon/cache/storer/storer"

	redigo "github.com/gomodule/redigo/redis"
)

// RedisStore RedisStore
type RedisStore struct {
	pool *redigo.Pool
}

// NewRedisStore NewRedisStore
func NewRedisStore(opts ...RedisOptions) *RedisStore {
	options := RedisOptions{}

	if len(opts) > 0 {
		options = opts[0]
	}
	options = defaultRedisOptions(options)

	return &RedisStore{
		pool: newRedisPool(options),
	}
}

// NewWithConf NewWithConf
func (s *RedisStore) NewWithConf(jsonconf string) error {

	var options RedisOptions
	if len(jsonconf) > 0 {
		err := json.Unmarshal([]byte(jsonconf), &options)
		if err != nil {
			return err
		}
	}

	options = defaultRedisOptions(options)

	s.pool = newRedisPool(options)

	return nil
}

func (s *RedisStore) do(commandName string, args ...interface{}) (reply interface{}, err error) {
	c := s.pool.Get()
	defer c.Close()
	return c.Do(commandName, args...)
}

// Set Set
func (s *RedisStore) Set(key string, val []byte, ttl int64) error {
	_, err := s.do("SETEX", key, ttl, val)
	if err != nil {
		return err
	}

	return nil
}

// Get Get
func (s *RedisStore) Get(key string) ([]byte, error) {
	ret, err := redigo.Bytes(s.do("GET", key))
	if err != nil {
		if err == redigo.ErrNil {
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

// Expire Expire
func (s *RedisStore) Expire(key string, ttl int64) error {
	_, err := s.do("EXPIRE", key, ttl)
	if err != nil {
		return err
	}

	return nil
}

// Forever Forever
func (s *RedisStore) Forever(key string, val []byte) error {
	_, err := s.do("SET", key, val)
	if err != nil {
		return err
	}

	return nil
}

// Del Del
func (s *RedisStore) Del(key string) error {
	_, err := s.do("DEL", key)
	if err != nil {
		return err
	}

	return nil

}

func newRedisPool(conf RedisOptions) *redigo.Pool {
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

func init() {
	storer.Register("redis", &RedisStore{})
}
