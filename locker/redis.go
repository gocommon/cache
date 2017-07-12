package locker

import (
	"encoding/json"
	"time"

	redigo "github.com/garyburd/redigo/redis"
	redsync "gopkg.in/redsync.v1"
)

// RedisLocker RedisLocker
type RedisLocker struct {
	redsync *redsync.Redsync
	opts    RedisOptions
}

// NewRedisLocker NewRedisLocker
func NewRedisLocker(opts ...RedisOptions) *RedisLocker {
	options := RedisOptions{}

	if len(opts) > 0 {
		options = opts[0]
	}
	options = defaultOptions(options)

	pool := newRedisPool(options)

	redsync := redsync.New([]redsync.Pool{pool})

	return &RedisLocker{redsync: redsync, opts: options}
}

// NewWithConf NewWithConf
func (l *RedisLocker) NewWithConf(jsonconf string) error {
	var options RedisOptions

	err := json.Unmarshal([]byte(jsonconf), &options)
	if err != nil {
		return err
	}

	options = defaultOptions(options)

	pool := newRedisPool(options)

	redsync := redsync.New([]redsync.Pool{pool})

	l.redsync = redsync
	l.opts = options

	return nil
}

// NewLocker NewLocker
func (l *RedisLocker) NewLocker(key string) Funcer {
	// do not retry!
	m := l.redsync.NewMutex(key, redsync.SetTries(1))
	return &RedisLockFunc{
		mutex: m,
	}
}

// RedisLockFunc RedisLockFunc
type RedisLockFunc struct {
	mutex *redsync.Mutex
}

// Lock Lock
func (l *RedisLockFunc) Lock() error {
	err := l.mutex.Lock()
	if err != nil {
		if err == redsync.ErrFailed {
			return ErrLockFailed
		}

		return err
	}

	return nil
}

// Unlock Unlock
func (l *RedisLockFunc) Unlock() error {

	ok := l.mutex.Unlock()
	if !ok {
		return ErrUnlockFailed
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
