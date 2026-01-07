package cache

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gocommon/cache/v2/store"
	"github.com/gocommon/cache/v2/store/memory"
	redisStore "github.com/gocommon/cache/v2/store/redis"
)

// CacheConfig 表示从DSN解析出的缓存配置
type CacheConfig struct {
	Store    string            // 存储类型: "memory" 或 "redis"
	Username string            // Redis用户名
	Password string            // Redis密码
	Address  string            // Redis地址 (host:port)
	DB       int               // Redis数据库编号
	TTL      int64             // 默认TTL (秒)
	Prefix   string            // Key前缀
	TagTTL   int64             // 标签TTL (秒，-1表示永不过期)
	TouchTTL int64             // 自动续期时间 (秒)
	Cleanup  time.Duration     // 内存清理间隔 (仅memory有效)
	Options  map[string]string // 其他选项
}

// ParseDSN 解析cache DSN字符串
//
// DSN格式:
// [store://][username[:password]@][address[:port]][/db][?query][#fragment]
//
// 支持的存储类型:
// - memory:// - 内存存储
// - redis:// - Redis存储
//
// Query参数:
// - ttl: 默认缓存过期时间(秒)
// - prefix: key前缀
// - tag_ttl: 标签过期时间(秒，-1表示永不过期)
// - touch_ttl: 自动续期时间(秒)
// - cleanup: 内存清理间隔(仅memory，如30s, 5m)
//
// 示例:
// - memory://?ttl=300&prefix=app:
// - redis://localhost:6379/0?ttl=3600&prefix=cache:
// - redis://user:pass@remote.host:6379/1?ttl=7200&tag_ttl=-1
func ParseDSN(dsn string) (*CacheConfig, error) {
	if dsn == "" {
		return nil, NewCacheError(ErrCodeValidationError, "DSN cannot be empty")
	}

	u, err := url.Parse(dsn)
	if err != nil {
		return nil, WrapError(ErrCodeValidationError, "invalid DSN format", err)
	}

	config := &CacheConfig{
		Options: make(map[string]string),
	}

	// 解析存储类型
	config.Store = strings.TrimSuffix(u.Scheme, "://")
	if config.Store == "" {
		return nil, NewCacheError(ErrCodeValidationError, "store type is required in DSN")
	}

	switch config.Store {
	case "memory":
		if err := config.parseMemoryDSN(u); err != nil {
			return nil, err
		}
	case "redis":
		if err := config.parseRedisDSN(u); err != nil {
			return nil, err
		}
	default:
		return nil, NewCacheError(ErrCodeValidationError, fmt.Sprintf("unsupported store type: %s", config.Store))
	}

	// 解析通用查询参数
	if err := config.parseQueryParams(u.Query()); err != nil {
		return nil, err
	}

	return config, nil
}

// parseMemoryDSN 解析内存存储DSN
func (c *CacheConfig) parseMemoryDSN(u *url.URL) error {
	// 内存存储不需要用户名密码和地址
	if u.User != nil {
		return NewCacheError(ErrCodeValidationError, "memory store does not support authentication")
	}
	return nil
}

// parseRedisDSN 解析Redis存储DSN
func (c *CacheConfig) parseRedisDSN(u *url.URL) error {
	// 解析认证信息
	if u.User != nil {
		c.Username = u.User.Username()
		c.Password, _ = u.User.Password()
	}

	// 解析地址
	if u.Host != "" {
		c.Address = u.Host
	} else {
		c.Address = "localhost:6379" // 默认Redis地址
	}

	// 解析数据库编号
	if u.Path != "" && u.Path != "/" {
		dbStr := strings.TrimPrefix(u.Path, "/")
		db, err := strconv.Atoi(dbStr)
		if err != nil {
			return WrapError(ErrCodeValidationError, "invalid database number", err)
		}
		c.DB = db
	}

	return nil
}

// parseQueryParams 解析查询参数
func (c *CacheConfig) parseQueryParams(values url.Values) error {
	for key, vals := range values {
		if len(vals) == 0 {
			continue
		}
		val := vals[0]

		switch key {
		case "ttl":
			ttl, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return WrapError(ErrCodeValidationError, "invalid ttl value", err)
			}
			c.TTL = ttl

		case "prefix":
			c.Prefix = val

		case "tag_ttl":
			tagTTL, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return WrapError(ErrCodeValidationError, "invalid tag_ttl value", err)
			}
			c.TagTTL = tagTTL

		case "touch_ttl":
			touchTTL, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return WrapError(ErrCodeValidationError, "invalid touch_ttl value", err)
			}
			c.TouchTTL = touchTTL

		case "cleanup":
			if c.Store == "memory" {
				cleanup, err := time.ParseDuration(val)
				if err != nil {
					return WrapError(ErrCodeValidationError, "invalid cleanup duration", err)
				}
				c.Cleanup = cleanup
			}

		default:
			// 保存其他选项
			c.Options[key] = val
		}
	}

	// 设置默认值
	c.setDefaults()

	return nil
}

// setDefaults 设置默认值
func (c *CacheConfig) setDefaults() {
	if c.TTL == 0 {
		c.TTL = DefaultTTL
	}
	if c.TagTTL == 0 {
		c.TagTTL = DefaultTagTTL
	}
	if c.TouchTTL == 0 {
		c.TouchTTL = DefaultTouchTTL
	}
	if c.Prefix == "" {
		c.Prefix = DefaultPrefix
	}
	if c.Store == "memory" && c.Cleanup == 0 {
		c.Cleanup = time.Minute // 默认1分钟清理间隔
	}
}

// NewCache 从配置创建Cache实例
func (c *CacheConfig) NewCache() (*Cache, error) {
	var storeImpl store.Store

	switch c.Store {
	case "memory":
		if c.Cleanup > 0 {
			storeImpl = memory.NewMemoryWithCleanupInterval(c.Cleanup)
		} else {
			storeImpl = memory.NewMemory()
		}

	case "redis":
		// 为Redis DSN创建默认客户端
		redisOptions := &redis.Options{
			Addr: c.Address,
			DB:   c.DB,
		}

		// 如果有认证信息，添加到选项中
		if c.Username != "" || c.Password != "" {
			if c.Username != "" {
				redisOptions.Username = c.Username
			}
			redisOptions.Password = c.Password
		}

		// 如果没有指定地址，使用默认地址
		if redisOptions.Addr == "" {
			redisOptions.Addr = "localhost:6379"
		}

		storeImpl = redisStore.NewRedis(redis.NewClient(redisOptions))

	default:
		return nil, NewCacheError(ErrCodeValidationError, fmt.Sprintf("unsupported store type: %s", c.Store))
	}

	options := []Option{
		WithStore(storeImpl),
		WithTTL(c.TTL),
		WithPrefix(c.Prefix),
		WithTagTTL(c.TagTTL),
		WithTouchTTL(c.TouchTTL),
	}

	cache, err := New(options...)
	if err != nil {
		return nil, err
	}
	return cache, nil
}

// String 返回DSN字符串表示
func (c *CacheConfig) String() string {
	scheme := c.Store
	if !strings.HasSuffix(scheme, "://") {
		scheme += "://"
	}

	var result strings.Builder
	result.WriteString(scheme)

	if c.Store == "redis" {
		if c.Username != "" || c.Password != "" {
			result.WriteString(c.Username)
			if c.Password != "" {
				result.WriteByte(':')
				result.WriteString(c.Password)
			}
			result.WriteByte('@')
		}
		if c.Address != "" {
			result.WriteString(c.Address)
		}
		if c.DB != 0 {
			result.WriteByte('/')
			result.WriteString(strconv.Itoa(c.DB))
		}
	}

	query := url.Values{}

	if c.TTL != DefaultTTL {
		query.Set("ttl", strconv.FormatInt(c.TTL, 10))
	}
	if c.Prefix != DefaultPrefix {
		query.Set("prefix", c.Prefix)
	}
	if c.TagTTL != DefaultTagTTL {
		query.Set("tag_ttl", strconv.FormatInt(c.TagTTL, 10))
	}
	if c.TouchTTL != DefaultTouchTTL {
		query.Set("touch_ttl", strconv.FormatInt(c.TouchTTL, 10))
	}
	if c.Store == "memory" && c.Cleanup != time.Minute {
		query.Set("cleanup", c.Cleanup.String())
	}

	for k, v := range c.Options {
		query.Set(k, v)
	}

	if encoded := query.Encode(); encoded != "" {
		result.WriteByte('?')
		result.WriteString(encoded)
	}

	return result.String()
}
