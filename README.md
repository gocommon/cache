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

### Basic Usage

```go
// Create cache instance
c, err := cache.New(cache.WithMemoryStore())
if err != nil {
    log.Fatal(err)
}

// Create session with tags
session := c.Tags(context.Background(), "user", "profile")

// Set cache
err = session.Set("user_data", userData)

// Get cache
var data UserData
has, err := session.Get("user_data", &data)

// Flush by tags (invalidate related cache)
err = session.Flush()
```

### DSN Usage (One-Line Configuration)

Just like database connections, you can create cache instances using DSN strings:

```go
// Memory cache with custom settings
cache, err := cache.NewFromDSN("memory://?ttl=300&prefix=myapp:&cleanup=30s")

// Redis cache (automatically creates Redis client)
cache, err := cache.NewFromDSN("redis://localhost:6379/0?ttl=3600&prefix=cache:")

// Redis with authentication
cache, err := cache.NewFromDSN("redis://user:pass@remote.host:6379/1?ttl=7200")

// For advanced Redis configurations, use NewFromDSN with custom client creation
```

#### DSN Format

```
[store://][username[:password]@][address[:port]][/db][?query]
```

#### Supported Stores

- **Memory**: `memory://`
- **Redis**: `redis://`

#### Query Parameters

- `ttl`: Default cache TTL in seconds (default: 7200)
- `prefix`: Key prefix (default: "tc.")
- `tag_ttl`: Tag TTL in seconds, -1 for permanent (default: -1)
- `touch_ttl`: Auto-touch TTL in seconds (default: 600)
- `cleanup`: Memory cleanup interval (memory only, default: "1m")

#### DSN Examples

```go
// Simple memory cache
"memory://"

// Memory with custom TTL and prefix
"memory://?ttl=300&prefix=myapp:"

// Redis basic
"redis://localhost:6379/0"

// Redis with auth and custom settings
"redis://user:pass@remote.host:6379/1?ttl=3600&prefix=cache:&tag_ttl=7200"
```

## Tag System

- Tags are sorted and combined to generate unique cache keys
- Tag values are cached to improve performance
- Use `Flush()` to invalidate all cache entries with the same tags