package cache

import (
	"context"
)

// Cacher Cacher
type Cacher interface {
	Tags(ctx context.Context, tags ...string) Session
}

type Session interface {
	Get(key string, val interface{}) (has bool, err error)
	GetWithVersion(key string, val interface{}) (has bool, version string, err error)
	Set(key string, val interface{}) error
	Del(key string) error
	// Version tag ver
	Version() (string, error)
	Flush() error
}

var _ Cacher = &Cache{}

// Cache Cache
type Cache struct {
	opts *Options
}

// New New
func New(opts ...Option) *Cache {
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

func (c *Cache) Tags(ctx context.Context, tags ...string) Session {
	return &session{ctx: ctx, tags: tags, opts: c.opts}
}
