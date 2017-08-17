package locker

var (
	// DefaultLocker DefaultLocker
	DefaultLocker = NewWithAdapter("redis", DefaultRedisConfigString)
	// NewRedisLocker([]RedisOptions{RedisOptions{}})
)
