package locker

import "errors"

var (
	_ Locker = &ErrLocker{}
	// ErrUnlockFailed ErrUnlockFailed
	ErrUnlockFailed = errors.New("Locker Unlock Failed")

	// ErrLockFailed ErrLockFailed
	ErrLockFailed = errors.New("Locker Lock Failed")

	// ErrLockerNotFound ErrLockerNotFound
	ErrLockerNotFound = errors.New("locker not found")
)

// IsErrLockFailed IsErrLockFailed
func IsErrLockFailed(err error) bool {
	return err == ErrLockFailed
}

// ErrLocker ErrLocker
type ErrLocker struct {
	err error
}

// NewErrLocker NewErrLocker
func NewErrLocker(err error) *ErrLocker {
	return &ErrLocker{err}
}

// NewLocker NewLocker
func (e *ErrLocker) NewLocker(key string) Funcer {
	return NewErrFuncer(e.err)
}

// NewWithConf NewWithConf
func (e *ErrLocker) NewWithConf(jsonconf string) error {
	return e.err
}

// ErrFuncer ErrFuncer
type ErrFuncer struct {
	err error
}

// NewErrFuncer NewErrFuncer
func NewErrFuncer(err error) *ErrFuncer {
	return &ErrFuncer{err}
}

// Lock Lock
func (e *ErrFuncer) Lock() error {
	return e.err
}

// Unlock Unlock
func (e *ErrFuncer) Unlock() error {
	return e.err

}
