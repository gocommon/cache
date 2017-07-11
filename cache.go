package cache

var _ TagCacher = &Cache{}

// Cache Cache
type Cache struct {
	opts *Options
}

// NewCache NewCache
func NewCache(opts ...Option) Cacher {
	options := &Options{}
	for i := range opts {
		opts[i](options)
	}

	options = defaultOptions(options)

	return &Cache{
		opts: options,
	}

}

func (c *Cache) keyWithPrefix(key string) string {
	return c.opts.Prefix + key
}

// Set Set
func (c *Cache) Set(key string, val interface{}) error {
	d := EmptyValue
	if val != nil {
		var err error
		d, err = c.opts.Codec.Encode(val)
		if err != nil {
			return err
		}
	}

	return c.opts.Store.Set(keyWithPrefix(key), d, c.opts.TTL)
}

// Get Get
func (c *Cache) Get(key string, val interface{}) error {
	d, err := c.opts.Store.Get(keyWithPrefix(key))
	if err != nil {
		return err
	}

	if d == EmptyValue {
		return ErrNil
	}

	return c.opts.Codec.Decode(d, val)

}

// Forever Forever
func (c *Cache) Forever(key string, val interface{}) error {
	d := EmptyValue
	if val != nil {
		var err error
		d, err = c.opts.Codec.Encode(val)
		if err != nil {
			return err
		}
	}
	return c.opts.Store.Forever(keyWithPrefix(key), d)

}

// Del Del
func (c *Cache) Del(key string) error {
	return c.opts.Store.Del(keyWithPrefix(key))

}

// Tags Tags
func (c *Cache) Tags(tags []string) Cacher {
	return &TagCache{
		tagSet: &TagSet{names: tags, opts: c.opts},
		cache:  c,
	}
}

// TagID TagID
func (c *Cache) TagID(tag string) string {
	return (&TagSet{names: []string{}, opts: c.opts}).TagID(tag)
}
