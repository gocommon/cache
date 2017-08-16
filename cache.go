package cache

import (
	"bytes"

	"time"

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

// NewWithOptions NewWithOptions
func NewWithOptions(opts Options) TagCacher {

	options := defaultOptions(&opts)

	return &Cache{
		opts: options,
	}
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
		d, err = c.opts.Codec.Encode(val)
		if err != nil {
			return err
		}

	}

	// add unix to the end @
	unix := time.Now().Unix()
	d = c.joinUnix(d, unix)

	return c.opts.Store.Set(c.keyWithPrefix(key), d, c.opts.TTL)
}

// Get Get
func (c *Cache) Get(key string, val interface{}) (has bool, err error) {
	src, err := c.opts.Store.Get(c.keyWithPrefix(key))
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
		c.opts.Store.Set(c.keyWithPrefix(key), d, c.opts.TTL)
	}

	if bytes.Contains(d, EmptyValue) {
		// SetNil(val)
		return true, nil
	}

	return true, c.opts.Codec.Decode(d, val)

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

	// forever set unix = 0
	d = c.joinUnix(d, 0)

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
	if c.opts.Locker == nil {
		return NewErrLocker().NewLocker()
	}
	return c.opts.Locker.NewLocker(c.keyWithPrefix(key))
}
