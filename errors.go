package cache

import (
	"errors"

	"github.com/gocommon/cache/locker"
)

var (
	// ErrNil ErrNil
	ErrNil             = errors.New("key not found")
	ErrLockerUndefined = errors.New("locker undefined")
)

// IsErrNil IsErrNil
func IsErrNil(err error) bool {
	return ErrNil == err
}

type ErrLocker struct {
}

func NewErrLocker() *ErrLocker {
	return &ErrLocker{}
}

func (e *ErrLocker) NewLocker() locker.Funcer {
	return &ErrLockerFuncer{}
}

type ErrLockerFuncer struct {
}

func (e *ErrLockerFuncer) Lock() error {
	return ErrLockerUndefined
}

func (e *ErrLockerFuncer) Unlock() error {
	return ErrLockerUndefined
}
