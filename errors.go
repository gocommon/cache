package cache

import (
	"errors"
)

var (
	// ErrNil ErrNil
	ErrNil = errors.New("key not found")
)

// IsErrNil IsErrNil
func IsErrNil(err error) bool {
	return ErrNil == err
}

var _ Cacher = &ErrCacher{}

type ErrCacher struct {
	err error
}

func NewErrCacher(err error) Cacher {
	return &ErrCacher{err}
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
