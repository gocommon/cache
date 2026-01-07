// +build ignore

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gocommon/cache/v2"
)

func main() {
	fmt.Println("=== Cache DSN Simplified API Demo ===")

	// 1. Memory cache - 简单
	fmt.Println("\n1. Memory Cache (Simple):")
	memDSN := "memory://?ttl=300&prefix=app:"
	fmt.Printf("DSN: %s\n", memDSN)

	memCache, err := cache.NewFromDSN(memDSN)
	if err != nil {
		log.Printf("Memory cache error: %v", err)
	} else {
		fmt.Println("✓ Memory cache created successfully")
		testCacheOperations(memCache, "memory")
	}

	// 2. Redis cache - 自动创建客户端
	fmt.Println("\n2. Redis Cache (Auto-create client):")
	redisDSN := "redis://localhost:6379/0?ttl=600&prefix=redis:"
	fmt.Printf("DSN: %s\n", redisDSN)

	redisCache, err := cache.NewFromDSN(redisDSN)
	if err != nil {
		log.Printf("Redis cache error (expected if no Redis server): %v", err)
	} else {
		fmt.Println("✓ Redis cache created successfully")
		testCacheOperations(redisCache, "redis")
	}

	// 3. Redis专用函数 - 更明确
	fmt.Println("\n3. Redis Cache (Dedicated function):")
	redisDSN2 := "redis://user:pass@remote.host:6379/1?ttl=1800&prefix=api:"
	fmt.Printf("DSN: %s\n", redisDSN2)

	redisCache2, err := cache.NewRedisCacheFromDSN(redisDSN2)
	if err != nil {
		log.Printf("Redis cache error (expected if no Redis server): %v", err)
	} else {
		fmt.Println("✓ Redis cache created with dedicated function")
		testCacheOperations(redisCache2, "redis-dedicated")
	}

	fmt.Println("\n=== DSN API Comparison ===")
	fmt.Println("Before: Multiple lines for Redis setup")
	fmt.Println("  rdb := redis.NewClient(&redis.Options{Addr: \"localhost:6379\"})")
	fmt.Println("  cache := cache.New(cache.WithStore(redis.NewRedis(rdb)), ...)")
	fmt.Println("")
	fmt.Println("After: One line DSN")
	fmt.Println("  cache, err := cache.NewFromDSN(\"redis://localhost:6379/0?ttl=600\")")
	fmt.Println("")
	fmt.Println("Even simpler for Redis:")
	fmt.Println("  cache, err := cache.NewRedisCacheFromDSN(\"redis://localhost:6379/0?ttl=600\")")
}

func testCacheOperations(c *cache.Cache, cacheType string) {
	ctx := context.Background()
	session := c.Tags(ctx, "demo", cacheType)

	// Test set
	err := session.Set("demo_key", fmt.Sprintf("Hello from %s cache!", cacheType))
	if err != nil {
		log.Printf("Set error: %v", err)
		return
	}

	// Test get
	var result string
	has, err := session.Get("demo_key", &result)
	if err != nil {
		log.Printf("Get error: %v", err)
		return
	}

	if has {
		fmt.Printf("✓ Cache operation successful: %s\n", result)
	} else {
		fmt.Printf("✗ Cache miss for %s\n", cacheType)
	}
}