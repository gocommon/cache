package storer

import (
	"github.com/gocommon/cache/storer/redis"
)

var (
	// DefaultStore DefaultStore
	DefaultStore = redis.NewStore()
)
