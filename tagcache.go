package cache

var _ Cacher = &TagCache{}

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
	err := c.cache.Get(c.taggedItemKey(key), val)
	if err != nil {
		return err
	}

	return nil

}

// Forever Forever
func (c *TagCache) Forever(key string, val interface{}) error {
	err := c.cache.Forever(c.taggedItemKey(key), val)
	if err != nil {
		return err
	}

	return nil

}

// Del Del
func (c *TagCache) Del(key string) error {
	err := c.cache.Del(c.taggedItemKey(key))
	if err != nil {
		return err
	}

	return nil

}

// taggedItemKey real store key
func (c *TagCache) taggedItemKey(key string) string {
	return EncodeMD5(c.tagSet.GetNamespace() + key)
}
