package storer

// Storer store
type Storer interface {
	Set(key string, val []byte, ttl int64) error
	Get(key string) ([]byte, error)
	Expire(key string, ttl int64) error
	Forever(key string, val []byte) error
	Del(key string) error
	NewWithConf(jsonconf string) error
}
