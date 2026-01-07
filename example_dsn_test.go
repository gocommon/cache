package cache_test

import (
	"context"
	"log"
	"testing"

	"github.com/gocommon/cache/v2"
)

// UserData 示例数据结构
type UserProfile struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func ExampleNewFromDSN() {
	// 使用DSN创建内存缓存，就像连接数据库一样简单
	cache, err := cache.NewFromDSN("memory://?ttl=300&prefix=myapp:")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	session := cache.Tags(ctx, "user", "profile")

	user := UserProfile{
		ID:       123,
		Username: "john_doe",
		Email:    "john@example.com",
	}

	// 存储数据
	if err := session.Set("user_123", user); err != nil {
		log.Fatal(err)
	}

	// 检索数据
	var retrieved UserProfile
	has, err := session.Get("user_123", &retrieved)
	if err != nil {
		log.Fatal(err)
	}

	if has {
		log.Printf("Retrieved user: %+v", retrieved)
	} else {
		log.Println("User not found in cache")
	}
}

func TestDSNExamples(t *testing.T) {
	tests := []struct {
		name string
		dsn  string
		test func(t *testing.T, cache *cache.Cache)
	}{
		{
			name: "memory DSN",
			dsn:  "memory://?ttl=60&prefix=test:",
			test: func(t *testing.T, c *cache.Cache) {
				ctx := context.Background()
				session := c.Tags(ctx, "example")

				// Test basic operations
				err := session.Set("key", "value")
				if err != nil {
					t.Errorf("Set failed: %v", err)
				}

				var result string
				has, err := session.Get("key", &result)
				if err != nil {
					t.Errorf("Get failed: %v", err)
				}
				if !has || result != "value" {
					t.Errorf("Expected 'value', got '%s' (has=%v)", result, has)
				}
			},
		},
		{
			name: "memory DSN with cleanup",
			dsn:  "memory://?ttl=30&cleanup=10s",
			test: func(t *testing.T, c *cache.Cache) {
				ctx := context.Background()
				session := c.Tags(ctx, "cleanup_test")

				err := session.Set("temp", "data")
				if err != nil {
					t.Errorf("Set failed: %v", err)
				}

				// Verify data exists
				var result string
				has, err := session.Get("temp", &result)
				if err != nil || !has || result != "data" {
					t.Errorf("Data should exist immediately after setting")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache, err := cache.NewFromDSN(tt.dsn)
			if err != nil {
				t.Fatalf("NewFromDSN failed for %s: %v", tt.dsn, err)
			}

			tt.test(t, cache)
		})
	}
}

func TestDSNConfigParsing(t *testing.T) {
	config, err := cache.ParseDSN("memory://?ttl=300&prefix=myapp:&tag_ttl=3600&touch_ttl=120&cleanup=30s")
	if err != nil {
		t.Fatalf("ParseDSN failed: %v", err)
	}

	if config.Store != "memory" {
		t.Errorf("Expected store 'memory', got '%s'", config.Store)
	}
	if config.TTL != 300 {
		t.Errorf("Expected TTL 300, got %d", config.TTL)
	}
	if config.Prefix != "myapp:" {
		t.Errorf("Expected prefix 'myapp:', got '%s'", config.Prefix)
	}
	if config.TagTTL != 3600 {
		t.Errorf("Expected TagTTL 3600, got %d", config.TagTTL)
	}
	if config.TouchTTL != 120 {
		t.Errorf("Expected TouchTTL 120, got %d", config.TouchTTL)
	}
}

func TestDSNToString(t *testing.T) {
	original := "memory://?ttl=300&prefix=myapp:&tag_ttl=3600&touch_ttl=120&cleanup=30s"

	config, err := cache.ParseDSN(original)
	if err != nil {
		t.Fatalf("ParseDSN failed: %v", err)
	}

	// 转换回字符串（可能顺序不同，但内容应该相同）
	result := config.String()

	// 解析结果并比较关键配置
	resultConfig, err := cache.ParseDSN(result)
	if err != nil {
		t.Fatalf("ParseDSN result failed: %v", err)
	}

	if resultConfig.TTL != 300 || resultConfig.Prefix != "myapp:" {
		t.Errorf("DSN round-trip failed: original=%s, result=%s", original, result)
	}
}

func TestNewFromDSN_RedisConfigOnly(t *testing.T) {
	// 测试Redis DSN解析和配置创建（不实际连接）
	dsn := "redis://testuser:testpass@test.host:6380/5?ttl=1800&prefix=redis:"

	config, err := cache.ParseDSN(dsn)
	if err != nil {
		t.Fatalf("ParseDSN failed: %v", err)
	}

	// 验证解析结果
	if config.Store != "redis" {
		t.Errorf("Expected store 'redis', got '%s'", config.Store)
	}
	if config.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", config.Username)
	}
	if config.Password != "testpass" {
		t.Errorf("Expected password 'testpass', got '%s'", config.Password)
	}
	if config.Address != "test.host:6380" {
		t.Errorf("Expected address 'test.host:6380', got '%s'", config.Address)
	}
	if config.DB != 5 {
		t.Errorf("Expected DB 5, got %d", config.DB)
	}
	if config.TTL != 1800 {
		t.Errorf("Expected TTL 1800, got %d", config.TTL)
	}
	if config.Prefix != "redis:" {
		t.Errorf("Expected prefix 'redis:', got '%s'", config.Prefix)
	}

	// 注意：这里我们不调用NewCache来创建实际的Redis连接
	// 因为我们没有Redis服务器可用
	// 在实际使用中，你可以这样做：
	// cache, err := cache.NewFromDSN(dsn) // 这会创建Redis客户端并连接
	// 或者专门用于Redis:
	// cache, err := cache.NewRedisCacheFromDSN(dsn)
}