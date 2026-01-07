// +build ignore

package main

import (
	"fmt"
	"log"

	"github.com/gocommon/cache/v2"
)

func main() {
	fmt.Println("=== Cache DSN Examples ===")

	// 1. Memory cache with DSN
	fmt.Println("\n1. Memory Cache DSN:")
	memDSN := "memory://?ttl=300&prefix=myapp:&cleanup=30s"
	fmt.Printf("DSN: %s\n", memDSN)

	memCache, err := cache.NewFromDSN(memDSN)
	if err != nil {
		log.Printf("Failed to create memory cache: %v", err)
	} else {
		fmt.Println("✓ Memory cache created successfully")
		_ = memCache // Use the cache...
	}

	// 2. Redis cache with DSN (auto-creates Redis client)
	fmt.Println("\n2. Redis Cache DSN (auto-create client):")
	redisDSN := "redis://localhost:6379/0?ttl=3600&prefix=cache:"
	fmt.Printf("DSN: %s\n", redisDSN)

	redisCache, err := cache.NewFromDSN(redisDSN)
	if err != nil {
		log.Printf("Failed to create Redis cache: %v", err)
		log.Printf("Note: This is expected if Redis server is not running")
	} else {
		fmt.Println("✓ Redis cache created successfully")
		_ = redisCache // Use the cache...
	}

	// 3. Parse DSN to see configuration
	fmt.Println("\n3. DSN Configuration Parsing:")
	testDSN := "redis://user:pass@remote.host:6379/5?ttl=1800&prefix=api:"
	fmt.Printf("DSN: %s\n", testDSN)

	config, err := cache.ParseDSN(testDSN)
	if err != nil {
		log.Printf("Failed to parse DSN: %v", err)
	} else {
		fmt.Printf("Parsed Config:\n")
		fmt.Printf("  Store: %s\n", config.Store)
		fmt.Printf("  Username: %s\n", config.Username)
		fmt.Printf("  Password: %s\n", config.Password)
		fmt.Printf("  Address: %s\n", config.Address)
		fmt.Printf("  DB: %d\n", config.DB)
		fmt.Printf("  TTL: %d seconds\n", config.TTL)
		fmt.Printf("  Prefix: %s\n", config.Prefix)

		// Convert back to string
		reconstructed := config.String()
		fmt.Printf("Reconstructed DSN: %s\n", reconstructed)
	}

	// 4. Show different DSN formats
	fmt.Println("\n4. DSN Format Examples:")
	examples := []string{
		"memory://",
		"memory://?ttl=600&cleanup=1m",
		"redis://localhost:6379/0",
		"redis://user:pass@host:6379/1?ttl=3600&prefix=app:",
	}

	for _, ex := range examples {
		fmt.Printf("  %s\n", ex)
	}
}