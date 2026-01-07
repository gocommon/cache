package cache

import (
	"time"

	"github.com/gocommon/cache/v2/codec"
	"github.com/gocommon/cache/v2/codec/gob"
	"github.com/gocommon/cache/v2/store"
	"github.com/gocommon/cache/v2/store/memory"
)

const (
	// DefaultTTL 默认缓存过期时间 (2小时)
	DefaultTTL = 7200

	// DefaultTagTTL 默认标签过期时间 (-1表示永久)
	DefaultTagTTL = -1

	// DefaultTouchTTL 默认自动续期时间 (10分钟)
	DefaultTouchTTL = 600

	// DefaultPrefix 默认key前缀
	DefaultPrefix = "tc."
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

// WithMemoryStore 使用内存存储
func WithMemoryStore() Option {
	return func(o *Options) {
		o.store = memory.NewMemory()
	}
}

// WithMemoryStoreWithCleanup 使用内存存储并指定清理间隔
func WithMemoryStoreWithCleanup(cleanupInterval time.Duration) Option {
	return func(o *Options) {
		o.store = memory.NewMemoryWithCleanupInterval(cleanupInterval)
	}
}

func defaultOptions(opts *Options) {
	if opts.ttl == 0 {
		opts.ttl = DefaultTTL
	}

	if opts.tagTTL == 0 {
		opts.tagTTL = DefaultTagTTL
	}

	if opts.touchTTL == 0 {
		opts.touchTTL = DefaultTouchTTL
	}

	if len(opts.prefix) == 0 {
		opts.prefix = DefaultPrefix
	}

	if opts.store == nil {
		opts.store = store.DefaultErrStore
	}

	if opts.codec == nil {
		opts.codec = &gob.GobCodec{}
	}

}

// Validate 验证配置的有效性
func (o *Options) Validate() error {
	if o.ttl < 0 {
		return NewCacheError(ErrCodeValidationError, ErrMsgInvalidTTL)
	}
	if o.tagTTL < 0 && o.tagTTL != -1 {
		return NewCacheError(ErrCodeValidationError, "tagTTL must be non-negative or -1")
	}
	if o.touchTTL < 0 {
		return NewCacheError(ErrCodeValidationError, "touchTTL must be non-negative")
	}
	if o.store == nil {
		return NewCacheError(ErrCodeValidationError, ErrMsgStoreRequired)
	}
	return nil
}
