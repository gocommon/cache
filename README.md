# cache
cache with tags write in GO

go get github.com/gocommon/cache/v2

## Features

- **Tag-based caching**: Support for cache invalidation by tags
- **Multiple storage backends**: Redis, Memory
- **Auto-expiration**: Automatic cache expiration and renewal
- **Serialization**: Support for custom codecs (Gob, JSON, etc.)

## Storage Backends

### Memory Storage

For development and testing, you can use in-memory storage:

```go
import "github.com/gocommon/cache/v2"

// Create cache with memory storage
cache := cache.New(
    cache.WithMemoryStore(),
    cache.WithTTL(3600), // 1 hour
)

// Or with custom cleanup interval
cache := cache.New(
    cache.WithMemoryStoreWithCleanup(30*time.Minute),
    cache.WithTTL(3600),
)
```

### Redis Storage

For production use, Redis is recommended:

```go
import (
    "github.com/go-redis/redis/v8"
    "github.com/gocommon/cache/v2"
    "github.com/gocommon/cache/v2/store/redis"
)

// Create Redis client
rdb := redis.NewClient(&redis.Options{
    Addr: "localhost:6379",
})

// Create cache with Redis storage
cache := cache.New(
    cache.WithStore(redis.NewRedis(rdb)),
    cache.WithTTL(7200), // 2 hours
)
```

## Usage

```go
// Create cache instance
c := cache.New(cache.WithMemoryStore())

// Create session with tags
session := c.Tags(context.Background(), "user", "profile")

// Set cache
err := session.Set("user_data", userData)

// Get cache
var data UserData
has, err := session.Get("user_data", &data)

// Flush by tags (invalidate related cache)
err = session.Flush()
```

## Tag System

- Tags are sorted and combined to generate unique cache keys
- Tag values are cached to improve performance
- Use `Flush()` to invalidate all cache entries with the same tags