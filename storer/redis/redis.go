package redis

import (
	"github.com/gocommon/cache/storer"
)

// NewStock NewStock
func NewStock(opts ...storer.RedisOptions) {
	return storer.NewRedisStore(opts...)
}
