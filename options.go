package cache

import (
	"github.com/gocommon/cache/v2/codec"
	"github.com/gocommon/cache/v2/codec/gob"
	"github.com/gocommon/cache/v2/store"
)

// Options Options
type Options struct {
	prefix   string
	ttl      int64 // key 有效期
	touchTTL int64 // 多少秒内访问，自动续期
	tagTTL   int64 // tagkey 有效期，默认-1，永久，如果想省内容空间，可以设置值

	store store.Store
	codec codec.Codec
}

type Option func(*Options)

// WithPrefix key前缀
func WithPrefix(s string) Option {
	return func(o *Options) {
		o.prefix = s
	}
}

// WithTTL 数据key 有效期， 到期自动失效
func WithTTL(s int64) Option {
	return func(o *Options) {
		o.ttl = s
	}
}

// WithTagTTL tag 有效期，默认-1，永久，如果想省内容空间，可以设置值
func WithTagTTL(s int64) Option {
	return func(o *Options) {
		o.tagTTL = s
	}
}

// WithTouchTTL  多少秒内访问，自动续期
func WithTouchTTL(s int64) Option {
	return func(o *Options) {
		o.touchTTL = s
	}
}

// WithStore 默认没有
func WithStore(s store.Store) Option {
	return func(o *Options) {
		o.store = s
	}
}

// WithCodec 默认gob
func WithCodec(c codec.Codec) Option {
	return func(o *Options) {
		o.codec = c
	}
}

func defaultOptions(opts *Options) {
	if opts.ttl == 0 {
		opts.ttl = 7200
	}

	if opts.tagTTL == 0 {
		opts.tagTTL = -1
	}

	if opts.touchTTL == 0 {
		opts.touchTTL = 600
	}

	if len(opts.prefix) == 0 {
		opts.prefix = "tc."
	}

	if opts.store == nil {
		opts.store = store.DefaultErrStore
	}

	if opts.codec == nil {
		opts.codec = &gob.GobCodec{}
	}

}
