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
	Tags(tags ...string) TagCacher
	Locker(key string) locker.Funcer
	Options() *Options
}

// TagCacher TagCacher
type TagCacher interface {
	Set(key string, val interface{}) error
	Get(key string, val interface{}) (has bool, err error)
	Forever(key string, val interface{}) error
	Del(key string) error
	Flush() error
	TagID(tag string) string
	SetTags(tags ...string)
}
