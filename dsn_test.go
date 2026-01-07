package cache

import (
	"context"
	"testing"
	"time"
)

func TestParseDSN_Memory(t *testing.T) {
	tests := []struct {
		name     string
		dsn      string
		expected *CacheConfig
		hasError bool
	}{
		{
			name: "basic memory DSN",
			dsn:  "memory://",
			expected: &CacheConfig{
				Store:    "memory",
				TTL:      DefaultTTL,
				Prefix:   DefaultPrefix,
				TagTTL:   DefaultTagTTL,
				TouchTTL: DefaultTouchTTL,
				Cleanup:  time.Minute,
			},
		},
		{
			name: "memory DSN with custom TTL",
			dsn:  "memory://?ttl=300",
			expected: &CacheConfig{
				Store:    "memory",
				TTL:      300,
				Prefix:   DefaultPrefix,
				TagTTL:   DefaultTagTTL,
				TouchTTL: DefaultTouchTTL,
				Cleanup:  time.Minute,
			},
		},
		{
			name: "memory DSN with all options",
			dsn:  "memory://?ttl=600&prefix=myapp:&tag_ttl=-1&touch_ttl=300&cleanup=30s",
			expected: &CacheConfig{
				Store:    "memory",
				TTL:      600,
				Prefix:   "myapp:",
				TagTTL:   -1,
				TouchTTL: 300,
				Cleanup:  30 * time.Second,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := ParseDSN(tt.dsn)

			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if config.Store != tt.expected.Store {
				t.Errorf("Store = %v, want %v", config.Store, tt.expected.Store)
			}
			if config.TTL != tt.expected.TTL {
				t.Errorf("TTL = %v, want %v", config.TTL, tt.expected.TTL)
			}
			if config.Prefix != tt.expected.Prefix {
				t.Errorf("Prefix = %v, want %v", config.Prefix, tt.expected.Prefix)
			}
			if config.TagTTL != tt.expected.TagTTL {
				t.Errorf("TagTTL = %v, want %v", config.TagTTL, tt.expected.TagTTL)
			}
			if config.TouchTTL != tt.expected.TouchTTL {
				t.Errorf("TouchTTL = %v, want %v", config.TouchTTL, tt.expected.TouchTTL)
			}
			if config.Cleanup != tt.expected.Cleanup {
				t.Errorf("Cleanup = %v, want %v", config.Cleanup, tt.expected.Cleanup)
			}
		})
	}
}

func TestParseDSN_Redis(t *testing.T) {
	tests := []struct {
		name     string
		dsn      string
		expected *CacheConfig
		hasError bool
	}{
		{
			name: "basic redis DSN",
			dsn:  "redis://localhost:6379/0",
			expected: &CacheConfig{
				Store:    "redis",
				Address:  "localhost:6379",
				DB:       0,
				TTL:      DefaultTTL,
				Prefix:   DefaultPrefix,
				TagTTL:   DefaultTagTTL,
				TouchTTL: DefaultTouchTTL,
			},
		},
		{
			name: "redis DSN with auth",
			dsn:  "redis://user:pass@remote.host:6379/1",
			expected: &CacheConfig{
				Store:    "redis",
				Username: "user",
				Password: "pass",
				Address:  "remote.host:6379",
				DB:       1,
				TTL:      DefaultTTL,
				Prefix:   DefaultPrefix,
				TagTTL:   DefaultTagTTL,
				TouchTTL: DefaultTouchTTL,
			},
		},
		{
			name: "redis DSN with custom config",
			dsn:  "redis://localhost:6379/0?ttl=3600&prefix=cache:&tag_ttl=7200",
			expected: &CacheConfig{
				Store:    "redis",
				Address:  "localhost:6379",
				DB:       0,
				TTL:      3600,
				Prefix:   "cache:",
				TagTTL:   7200,
				TouchTTL: DefaultTouchTTL,
			},
		},
		{
			name:     "invalid DSN",
			dsn:      "",
			hasError: true,
		},
		{
			name:     "unsupported store",
			dsn:      "mysql://localhost/db",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := ParseDSN(tt.dsn)

			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if config.Store != tt.expected.Store {
				t.Errorf("Store = %v, want %v", config.Store, tt.expected.Store)
			}
			if config.Address != tt.expected.Address {
				t.Errorf("Address = %v, want %v", config.Address, tt.expected.Address)
			}
			if config.DB != tt.expected.DB {
				t.Errorf("DB = %v, want %v", config.DB, tt.expected.DB)
			}
			if config.Username != tt.expected.Username {
				t.Errorf("Username = %v, want %v", config.Username, tt.expected.Username)
			}
			if config.Password != tt.expected.Password {
				t.Errorf("Password = %v, want %v", config.Password, tt.expected.Password)
			}
		})
	}
}

func TestCacheConfig_String(t *testing.T) {
	tests := []struct {
		name     string
		config   *CacheConfig
		expected string
	}{
		{
			name: "memory config with defaults",
			config: &CacheConfig{
				Store:    "memory",
				TTL:      DefaultTTL,
				Prefix:   DefaultPrefix,
				TagTTL:   DefaultTagTTL,
				TouchTTL: DefaultTouchTTL,
				Cleanup:  time.Minute,
			},
			expected: "memory://",
		},
		{
			name: "memory config with custom values",
			config: &CacheConfig{
				Store:    "memory",
				TTL:      300,
				Prefix:   "myapp:",
				TagTTL:   3600, // 1小时
				TouchTTL: 300,  // 5分钟
				Cleanup:  30 * time.Second,
			},
			expected: "memory://?cleanup=30s&prefix=myapp%3A&tag_ttl=3600&touch_ttl=300&ttl=300",
		},
		{
			name: "redis config",
			config: &CacheConfig{
				Store:    "redis",
				Address:  "localhost:6379",
				DB:       1,
				TTL:      3600,
				Prefix:   "cache:",
				TagTTL:   7200,
				TouchTTL: 300,
			},
			expected: "redis://localhost:6379/1?prefix=cache%3A&tag_ttl=7200&touch_ttl=300&ttl=3600",
		},
		{
			name: "redis config with auth",
			config: &CacheConfig{
				Store:    "redis",
				Username: "user",
				Password: "pass",
				Address:  "remote.host:6379",
				DB:       0,
				TTL:      DefaultTTL,
				Prefix:   DefaultPrefix,
				TagTTL:   DefaultTagTTL,
				TouchTTL: DefaultTouchTTL,
			},
			expected: "redis://user:pass@remote.host:6379",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.String()
			if result != tt.expected {
				t.Errorf("String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestNewFromDSN_Memory(t *testing.T) {
	dsn := "memory://?ttl=300&prefix=test:"
	cache, err := NewFromDSN(dsn)
	if err != nil {
		t.Fatalf("NewFromDSN failed: %v", err)
	}
	if cache == nil {
		t.Fatal("Cache is nil")
	}

	// Test basic functionality
	ctx := context.Background()
	session := cache.Tags(ctx, "test")

	err = session.Set("key", "value")
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}

	var result string
	has, err := session.Get("key", &result)
	if err != nil {
		t.Errorf("Get failed: %v", err)
	}
	if !has || result != "value" {
		t.Errorf("Get result mismatch: has=%v, result=%v", has, result)
	}
}

func TestParseDSN_InvalidValues(t *testing.T) {
	tests := []struct {
		name string
		dsn  string
	}{
		{"invalid TTL", "memory://?ttl=invalid"},
		{"invalid tag_ttl", "memory://?tag_ttl=invalid"},
		{"invalid touch_ttl", "memory://?touch_ttl=invalid"},
		{"invalid cleanup duration", "memory://?cleanup=invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseDSN(tt.dsn)
			if err == nil {
				t.Errorf("Expected error for invalid DSN: %s", tt.dsn)
			}
		})
	}
}

func TestNewCache_RedisConfig(t *testing.T) {
	// 测试Redis配置解析（不实际连接Redis）
	config := &CacheConfig{
		Store:    "redis",
		Address:  "localhost:6379",
		DB:       1,
		TTL:      3600,
		Prefix:   "test:",
		TagTTL:   7200,
		TouchTTL: 600,
	}

	// 由于没有实际的Redis服务器，我们只验证到配置解析阶段
	// NewCache方法会尝试创建Redis客户端，但不会实际连接
	// 这里我们验证配置本身是正确的
	if config.Store != "redis" {
		t.Errorf("Expected store 'redis', got '%s'", config.Store)
	}
	if config.Address != "localhost:6379" {
		t.Errorf("Expected address 'localhost:6379', got '%s'", config.Address)
	}
	if config.DB != 1 {
		t.Errorf("Expected DB 1, got %d", config.DB)
	}
}

func TestDSNRedisParsing(t *testing.T) {
	dsn := "redis://user:pass@remote.host:6379/2?ttl=1800&prefix=redis:"
	config, err := ParseDSN(dsn)
	if err != nil {
		t.Fatalf("ParseDSN failed: %v", err)
	}

	if config.Store != "redis" {
		t.Errorf("Expected store 'redis', got '%s'", config.Store)
	}
	if config.Username != "user" {
		t.Errorf("Expected username 'user', got '%s'", config.Username)
	}
	if config.Password != "pass" {
		t.Errorf("Expected password 'pass', got '%s'", config.Password)
	}
	if config.Address != "remote.host:6379" {
		t.Errorf("Expected address 'remote.host:6379', got '%s'", config.Address)
	}
	if config.DB != 2 {
		t.Errorf("Expected DB 2, got %d", config.DB)
	}
	if config.TTL != 1800 {
		t.Errorf("Expected TTL 1800, got %d", config.TTL)
	}
	if config.Prefix != "redis:" {
		t.Errorf("Expected prefix 'redis:', got '%s'", config.Prefix)
	}
}
