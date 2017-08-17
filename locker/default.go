package locker

import (
	"github.com/gocommon/cache/locker/locker"
	redis "github.com/gocommon/cache/locker/redis"
)

var (
	// DefaultLocker DefaultLocker
	DefaultLocker = locker.NewWithAdapter("redis", redis.DefaultRedisConfigString)
	// NewRedisLocker([]RedisOptions{RedisOptions{}})
)
