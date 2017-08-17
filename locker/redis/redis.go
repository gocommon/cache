package locker

import (
	"encoding/json"
	"time"

	redigo "github.com/garyburd/redigo/redis"
	"github.com/gocommon/cache/locker/locker"
	redsync "gopkg.in/redsync.v1"
)

// DefaultRedisConfigString DefaultRedisConfigString
var DefaultRedisConfigString = "[{}]"

// RedisLocker RedisLocker
type RedisLocker struct {
	redsync *redsync.Redsync
}

// NewRedisLocker NewRedisLocker
func NewRedisLocker(opts ...[]RedisOptions) *RedisLocker {
	options := []RedisOptions{}

	if len(opts) > 0 {
		options = opts[0]
	}

	pools := make([]redsync.Pool, len(options))

	for i := 0; i < len(options); i++ {

		pools[i] = newRedisPool(defaultRedisOptions(options[i]))
	}

	redsync := redsync.New(pools)

	return &RedisLocker{redsync: redsync}
}

// NewWithConf NewWithConf
func (l *RedisLocker) NewWithConf(jsonconf string) error {
	var options []RedisOptions

	if len(jsonconf) == 0 {
		jsonconf = DefaultRedisConfigString
	}

	err := json.Unmarshal([]byte(jsonconf), &options)
	if err != nil {
		return err
	}

	pools := make([]redsync.Pool, len(options))

	for i := 0; i < len(options); i++ {

		pools[i] = newRedisPool(defaultRedisOptions(options[i]))
	}

	redsync := redsync.New(pools)

	l.redsync = redsync

	return nil
}

// NewLocker NewLocker
// NewMutex returns a new distributed mutex with given name.
// func (r *Redsync) NewMutex(name string, options ...Option) *Mutex {
// 	m := &Mutex{
// 		name:   name,
// 		expiry: 8 * time.Second,
// 		tries:  32,
// 		delay:  500 * time.Millisecond,
// 		factor: 0.01,
// 		quorum: len(r.pools)/2 + 1,
// 		pools:  r.pools,
// 	}
// 	for _, o := range options {
// 		o.Apply(m)
// 	}
// 	return m
// }
func (l *RedisLocker) NewLocker(key string) locker.Funcer {
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
			return locker.ErrLockFailed
		}

		return err
	}

	return nil
}

// Unlock Unlock
func (l *RedisLockFunc) Unlock() error {

	ok := l.mutex.Unlock()
	if !ok {
		return locker.ErrUnlockFailed
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
	locker.Register("redis", &RedisLocker{})
}
