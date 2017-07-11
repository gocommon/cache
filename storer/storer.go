package storer

// Storer store
type Storer interface {
	Set(key string, val string, ttl int64) error
	Get(key string) (string, error)
	Expire(key string, ttl int64) error
	Forever(key string, val string) error
	Del(key string) error
}
