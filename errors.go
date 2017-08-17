package cache

import (
	"errors"

	"github.com/gocommon/cache/storer/storer"
)

var (
	_ Cacher    = &ErrCacher{}
	_ TagCacher = &ErrTagCache{}
	// ErrNil ErrNil
	ErrNil = errors.New("key not found")
)

// IsErrNil IsErrNil
func IsErrNil(err error) bool {
	return ErrNil == err
}

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
	return NewErrTagCache(c.err)
}

func (c *ErrCacher) Options() Options {
	return Options{}
}
func (c *ErrCacher) Store() storer.Storer {
	return storer.NewErrStorer(c.err)
}

type ErrTagCache struct {
	err error
}

func NewErrTagCache(err error) TagCacher {
	return &ErrTagCache{err}
}

func (c *ErrTagCache) Set(key string, val interface{}) error {
	return c.err
}
func (c *ErrTagCache) Get(key string, val interface{}) (has bool, err error) {
	return false, c.err
}
func (c *ErrTagCache) Forever(key string, val interface{}) error {
	return c.err
}
func (c *ErrTagCache) Del(key string) error {
	return c.err
}
func (c *ErrTagCache) Flush() error {
	return c.err
}
func (c *ErrTagCache) TagID(tag string) string {
	return ""
}
func (c *ErrTagCache) SetTags(tags ...string) {
	return
}
