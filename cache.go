package cache

import (
	"context"
)

// EmptyValue EmptyValue
var EmptyValue = []byte("##empty- -!##")

// Cacher Cacher
type Cacher interface {
	Tag(ctx context.Context, tags ...string) Session
}

type Session interface {
	Get(key string, val interface{}) (has bool, err error)
	Set(key string, val interface{}) error
	Del(key string) error
	Flush() error
}

var _ Cacher = &Cache{}

// Cache Cache
type Cache struct {
	opts *Options
}

// New New
func New(opts ...Option) Cacher {
	options := &Options{}
	for _, op := range opts {
		op(options)
	}

	defaultOptions(options)

	c := &Cache{
		opts: options,
	}

	return c
}

func (c *Cache) Tag(ctx context.Context, tags ...string) Session {

	return &session{ctx: ctx, tags: tags, opts: c.opts}
}
