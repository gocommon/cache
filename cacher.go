package cache

import (
	"github.com/gocommon/cache/locker"
)

// Cacher Cacher
type Cacher interface {
	Set(key string, val interface{}) error
	Get(key string, val interface{}) (has bool, err error)
	Forever(key string, val interface{}) error
	Del(key string) error
}

// TagCacher TagCacher
type TagCacher interface {
	Cacher
	Tags(tags []string) Cacher
	TagID(tag string) string
	Flush(tags []string) error
	NewLocker(key string) locker.Funcer
}
