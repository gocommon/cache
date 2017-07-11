package cache

// TagCache Cache
type TagCache struct {
	tagSet *TagSet
	cache  Cacher
}

// Set Set
func (c *TagCache) Set(key string, val interface{}) error {
	return c.cache.Set(c.taggedItemKey(key), val)
}

// Get Get
func (c *TagCache) Get(key string, val interface{}) error {
	return nil

}

// Forever Forever
func (c *TagCache) Forever(key string, val interface{}) error {
	return nil

}

// Del Del
func (c *TagCache) Del(key string) error {
	return nil

}

// taggedItemKey real store key
func (c *TagCache) taggedItemKey(key string) string {
	return EncodeMD5(c.tagSet.GetNamespace() + key)
}
