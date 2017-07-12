package cache

import "errors"

var (
	// ErrNil ErrNil
	ErrNil = errors.New("key not found")
)

// IsErrNil IsErrNil
func IsErrNil(err error) bool {
	return ErrNil == err
}
