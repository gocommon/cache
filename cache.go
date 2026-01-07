package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gocommon/cache/v2/store"
	redisStore "github.com/gocommon/cache/v2/store/redis"
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

	// 验证配置
	if err := options.Validate(); err != nil {
		// 如果配置无效，返回错误缓存
		return &Cache{opts: &Options{store: store.NewErrStore(err)}}
	}

	c := &Cache{
		opts: options,
	}

	return c
}

func (c *Cache) Tags(ctx context.Context, tags ...string) Session {
	return NewSession(ctx, tags, c.opts)
}

// NewMemoryCache 创建使用内存存储的缓存实例
func NewMemoryCache(ttl int64, opts ...Option) *Cache {
	options := []Option{
		WithMemoryStore(),
		WithTTL(ttl),
	}
	options = append(options, opts...)
	return New(options...)
}

// NewMemoryCacheWithCleanup 创建使用内存存储并指定清理间隔的缓存实例
func NewMemoryCacheWithCleanup(ttl int64, cleanupInterval time.Duration, opts ...Option) *Cache {
	options := []Option{
		WithMemoryStoreWithCleanup(cleanupInterval),
		WithTTL(ttl),
	}
	options = append(options, opts...)
	return New(options...)
}

// NewRedisCache 创建使用Redis存储的缓存实例
func NewRedisCache(rdb redis.Cmdable, ttl int64, opts ...Option) *Cache {
	options := []Option{
		WithStore(redisStore.NewRedis(rdb.(*redis.Client))),
		WithTTL(ttl),
	}
	options = append(options, opts...)
	return New(options...)
}
