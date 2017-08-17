package cache

import (
	"github.com/gocommon/cache/codec"
	lockerd "github.com/gocommon/cache/locker"
	"github.com/gocommon/cache/locker/locker"
	"github.com/gocommon/cache/storer"
)

// Options Options
type Options struct {
	Prefix    string
	TTL       int64 // key 有效期
	TouchTTL  int64 // 多少秒内访问，自动续期
	Store     storer.Storer
	Codec     codec.Codec
	Locker    locker.Locker
	TagTTL    int64 // tagkey 有效期，默认-1，永久，如果想省内容空间，可以设置值
	UseLocker bool
}

// Option Option
type Option func(*Options)

func defaultOptions(opts *Options) *Options {
	if opts.TTL == 0 {
		opts.TTL = 60
	}

	if opts.TagTTL == 0 {
		opts.TagTTL = -1
	}

	if opts.TouchTTL == 0 {
		opts.TouchTTL = 30
	}

	if len(opts.Prefix) == 0 {
		opts.Prefix = "tagcache."
	}

	if opts.Store == nil {
		opts.Store = storer.DefaultStore
	}

	if opts.Codec == nil {
		opts.Codec = codec.DefaultCodec
	}

	if opts.Locker == nil && opts.UseLocker {
		opts.Locker = lockerd.DefaultLocker
	}

	return opts
}

// Prefix Prefix
func Prefix(s string) Option {
	return func(o *Options) {
		o.Prefix = s
	}
}

// TTL TTL
func TTL(t int64) Option {
	return func(o *Options) {
		o.TTL = t
	}
}

// TagTTL TagTTL
func TagTTL(t int64) Option {
	return func(o *Options) {
		o.TagTTL = t
	}
}

// ForeverTagTTL ForeverTagTTL
func ForeverTagTTL() Option {
	return func(o *Options) {
		o.TagTTL = -1
	}
}

// Store Store
func Store(s storer.Storer) Option {
	return func(o *Options) {
		o.Store = s
	}
}

// Codec Codec
func Codec(s codec.Codec) Option {
	return func(o *Options) {
		o.Codec = s
	}
}

// UseLocker UseLocker
func UseLocker(use bool) Option {
	return func(o *Options) {
		o.UseLocker = use
	}
}
