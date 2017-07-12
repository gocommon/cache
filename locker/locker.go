package locker

// Locker Locker
type Locker interface {
	NewLocker(key string) Funcer
}

// Funcer Funcer
type Funcer interface {
	Lock() error
	Unlock() error
}
