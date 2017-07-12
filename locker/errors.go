package locker

import "errors"

var (
	// ErrUnlockFailed ErrUnlockFailed
	ErrUnlockFailed = errors.New("Locker Unlock Failed")

	// ErrLockFailed ErrLockFailed
	ErrLockFailed = errors.New("Locker Lock Failed")
)

// IsErrLockFailed IsErrLockFailed
func IsErrLockFailed(err error) bool {
	return err == ErrLockFailed
}
