package locker

// Locker Locker
type Locker interface {
	NewLocker(key string) Funcer
	NewWithConf(jsonconf string) error
}

// Funcer Funcer
type Funcer interface {
	Lock() error
	Unlock() error
}
