package locker

var (
	// DefaultLocker DefaultLocker
	DefaultLocker = NewRedisLocker()
)
