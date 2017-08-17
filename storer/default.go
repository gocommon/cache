package storer

import (
	_ "github.com/gocommon/cache/storer/redis"
	"github.com/gocommon/cache/storer/storer"
)

var (
	// DefaultStore DefaultStore
	DefaultStore = storer.NewWithAdapter("redis", "")
)
