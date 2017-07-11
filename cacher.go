package cache

// Cache Cache
type Cacher interface {
	Set(key string, val interface{}) error
	Get(key string, val interface{}) error
	Forever(key string, val interface{}) error
	Del(key string) error
}

// TagCacher TagCacher
type TagCacher interface {
	Cacher
	Tags(tags []string) Cacher
	TagID(tag string) string
}
