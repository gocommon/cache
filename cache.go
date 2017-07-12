package cache

import (
	"github.com/gocommon/cache/locker"
)

var _ TagCacher = &Cache{}

// Cache Cache
type Cache struct {
	opts *Options
}

// NewCache NewCache
func NewCache(opts ...Option) TagCacher {
	options := &Options{}
	for i := range opts {
		opts[i](options)
	}

	options = defaultOptions(options)

	return &Cache{
		opts: options,
	}

}

func (c *Cache) keyWithPrefix(key string) string {
	return c.opts.Prefix + key
}

// Set Set
func (c *Cache) Set(key string, val interface{}) error {
	d := EmptyValue

	if !IsNil(val) {
		var err error
		d, err = c.opts.Codec.Encode(val)
		if err != nil {
			return err
		}
	}

	return c.opts.Store.Set(c.keyWithPrefix(key), d, c.opts.TTL)
}

// Get Get
func (c *Cache) Get(key string, val interface{}) error {
	d, err := c.opts.Store.Get(c.keyWithPrefix(key))
	if err != nil {
		return err
	}

	if len(d) == 0 {
		return ErrNil
	}

	if d == EmptyValue {
		SetNil(val)
		return nil

	}

	return c.opts.Codec.Decode(d, val)

}

// Forever Forever
func (c *Cache) Forever(key string, val interface{}) error {
	d := EmptyValue
	if !IsNil(val) {
		var err error
		d, err = c.opts.Codec.Encode(val)
		if err != nil {
			return err
		}
	}
	return c.opts.Store.Forever(c.keyWithPrefix(key), d)

}

// Del Del
func (c *Cache) Del(key string) error {
	return c.opts.Store.Del(c.keyWithPrefix(key))

}

// Tags Tags
func (c *Cache) Tags(tags []string) Cacher {
	return &TagCache{
		tagSet: &TagSet{names: tags, opts: c.opts},
		cache:  c,
	}
}

// TagID TagID
func (c *Cache) TagID(tag string) string {
	return (&TagSet{names: []string{}, opts: c.opts}).TagID(tag)
}

// Flush Flush
func (c *Cache) Flush(tags []string) error {
	tagSet := &TagSet{names: []string{}, opts: c.opts}
	for i := range tags {
		tagSet.ResetTag(tags[i])
	}

	return nil
}

// NewLocker NewLocker
func (c *Cache) NewLocker(key string) locker.Funcer {
	return c.opts.Locker.NewLocker(key)
}
