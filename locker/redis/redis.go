package redis

import (
	"github.com/gocommon/cache/locker"
)

// NewLocker NewLocker
func NewLocker(opts ...locker.RedisOptions) locker.Locker {
	return locker.NewRedisLocker(opts...)
}
