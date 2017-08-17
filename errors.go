package cache

import (
	"errors"

	"github.com/gocommon/cache/locker/locker"
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
	err error
}

func NewErrLocker(err error) locker.Locker {
	return &ErrLocker{err}
}

func (e *ErrLocker) NewLocker(key string) locker.Funcer {
	return NewErrLockerFuncer(e.err)
}

func (e *ErrLocker) NewWithConf(conf string) error {
	return e.err
}

type ErrLockerFuncer struct {
	err error
}

func NewErrLockerFuncer(err error) locker.Funcer {
	return &ErrLockerFuncer{err}
}

func (e *ErrLockerFuncer) Lock() error {
	return e.err
}

func (e *ErrLockerFuncer) Unlock() error {
	return e.err
}

var _ Cacher = &ErrCacher{}

type ErrCacher struct {
	err error
}

func (c *ErrCacher) Set(key string, val interface{}) error {
	return c.err
}
func (c *ErrCacher) Get(key string, val interface{}) (has bool, err error) {
	return false, c.err
}
func (c *ErrCacher) Forever(key string, val interface{}) error {
	return c.err
}
func (c *ErrCacher) Del(key string) error {
	return c.err
}
func (c *ErrCacher) Tags(tags ...string) TagCacher {
	return NewErrTagCacher(c.err)
}
func (c *ErrCacher) Locker(key string) locker.Funcer {
	return NewErrLockerFuncer(c.err)
}
func (c *ErrCacher) Options() *Options {
	return &Options{}
}

type ErrTagCacher struct {
	err error
}

func NewErrTagCacher(err error) TagCacher {
	return &ErrTagCacher{err}
}

func (c *ErrTagCacher) Set(key string, val interface{}) error {
	return c.err
}
func (c *ErrTagCacher) Get(key string, val interface{}) (has bool, err error) {
	return false, c.err
}
func (c *ErrTagCacher) Forever(key string, val interface{}) error {
	return c.err
}
func (c *ErrTagCacher) Del(key string) error {
	return c.err
}
func (c *ErrTagCacher) Flush() error {
	return c.err
}
func (c *ErrTagCacher) TagID(tag string) string {
	return ""
}
func (c *ErrTagCacher) SetTags(tags ...string) {
	return
}
