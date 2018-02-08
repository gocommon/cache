package cache

import (
	"bytes"
	"sync"

	"time"

	codecd "github.com/gocommon/cache/codec"
	"github.com/gocommon/cache/codec/codec"
	storerd "github.com/gocommon/cache/storer"
	"github.com/gocommon/cache/storer/storer"
)

var _ Cacher = &Cache{}

// Cache Cache
type Cache struct {
	opts  Options
	pool  sync.Pool
	store storer.Storer
	codec codec.Codec
}

// New New
func New(opts ...Options) Cacher {
	options := Options{}
	if len(opts) > 0 {
		options = opts[0]
	}

	options = defaultOptions(options)

	c := &Cache{
		opts: options,
	}

	c.pool.New = func() interface{} {
		return &TagCache{}
	}

	c.store = storerd.DefaultStore
	if len(options.StoreAdapter) > 0 {
		c.store = storer.NewWithAdapter(options.StoreAdapter, options.StoreAdapterConfig)
	}

	c.codec = codecd.DefaultCodec
	return c
}

func (c *Cache) keyWithPrefix(key string) string {
	return c.opts.Prefix + key
}

func (c *Cache) splitUnix(src []byte) (data []byte, unix int64) {
	idx := len(src) - 9

	flag := src[idx : idx+1]
	if idx < 0 || flag[0] != '@' {
		return src, 0
	}

	return src[0:idx], int64(BytesToUint64(src[idx+1:]))

}

func (c *Cache) joinUnix(data []byte, unix int64) []byte {
	buf := bytes.NewBuffer(data)
	buf.WriteByte(byte('@'))
	buf.Write(Uint64ToBytes(uint64(unix)))

	return buf.Bytes()
}

// Set Set
func (c *Cache) Set(key string, val interface{}) error {
	d := EmptyValue

	if !IsNil(val) {
		var err error
		d, err = c.codec.Encode(val)
		if err != nil {
			return err
		}

	}

	// add unix to the end @
	unix := time.Now().Unix()
	d = c.joinUnix(d, unix)

	return c.store.Set(c.keyWithPrefix(key), d, c.opts.TTL)
}

// Get Get
func (c *Cache) Get(key string, val interface{}) (has bool, err error) {
	src, err := c.store.Get(c.keyWithPrefix(key))
	if err != nil {
		return false, err
	}

	if len(src) == 0 {
		return false, nil
	}

	d, unix := c.splitUnix(src)

	// near expire
	if unix > 0 && unix+c.opts.TTL-time.Now().Unix() < c.opts.TouchTTL {
		unix := time.Now().Unix()
		d = c.joinUnix(d, unix)
		c.store.Set(c.keyWithPrefix(key), d, c.opts.TTL)
	}

	if bytes.Contains(d, EmptyValue) {
		// SetNil(val)
		return true, nil
	}

	return true, c.codec.Decode(d, val)

}

// Forever Forever
func (c *Cache) Forever(key string, val interface{}) error {
	d := EmptyValue
	if !IsNil(val) {
		var err error
		d, err = c.codec.Encode(val)
		if err != nil {
			return err
		}
	}

	// forever set unix = 0
	d = c.joinUnix(d, 0)

	return c.store.Forever(c.keyWithPrefix(key), d)

}

// Del Del
func (c *Cache) Del(key string) error {
	return c.store.Del(c.keyWithPrefix(key))

}

// Tags Tags
func (c *Cache) Tags(tags ...string) TagCacher {
	tc := c.getTagCache()
	tc.SetTags(tags...)
	return tc
}

// GetTagCache GetTagCache
func (c *Cache) getTagCache() TagCacher {
	// tc := c.pool.Get().(*TagCache)
	// if tc.cache == nil {
	// 	tc.cache = c
	// }
	// return tc

	return &TagCache{
		cache: c,
	}
}

// ReleaseTagCache ReleaseTagCache
func (c *Cache) ReleaseTagCache(tc TagCacher) {
	// c.pool.Put(tc)
}

// Options Options
func (c *Cache) Options() Options {
	return c.opts
}

// Store Store
func (c *Cache) Store() storer.Storer {
	return c.store
}
