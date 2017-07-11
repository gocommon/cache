package cache

import (
	"github.com/gocommon/cache/codec"
	"github.com/gocommon/cache/storer"
)

// Options Options
type Options struct {
	Prefix string
	TTL    int64
	Store  storer.Storer
	Codec  codec.Codec
	TagTTL int64
}

// Option Option
type Option func(*Options)

func defaultOptions(opts *Options) *Options {
	if opts.TTL == 0 {
		opts.TTL = 3600
	}

	if opts.TagTTL == 0 {
		opts.TagTTL = -1
	}

	if len(opts.Prefix) == 0 {
		opts.Prefix = "tagcache:"
	}

	if opts.Store == nil {
		opts.Store = storer.DefaultStore
	}

	if opts.Codec == nil {
		opts.Codec = codec.DefaultCodec
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
