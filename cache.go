package cache

import (
	"bytes"
	"context"

	codecd "github.com/gocommon/cache/v2/codec"
	"github.com/gocommon/cache/v2/codec/codec"
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
	opts  Options
	store Storer
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

	c.store = NewErrStorer(ErrStorerNotFound)
	if len(options.StoreAdapter) > 0 {
		c.store = InitStore(options.StoreAdapter, options.StoreAdapterConfig)
	}

	c.codec = codecd.DefaultCodec
	return c
}

func (c *Cache) Tag(ctx context.Context, tags ...string) Session {

	return &session{ctx: ctx, tags: tags, c: c}
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
