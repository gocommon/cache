package cache

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

var _ TagCacher = &TagCache{}

// TagCache Cache
type TagCache struct {
	cache Cacher
	names []string
}

// Set Set
func (c *TagCache) Set(key string, val interface{}) error {
	return c.cache.Set(c.taggedItemKey(key), val)
}

// Get Get
func (c *TagCache) Get(key string, val interface{}) (bool, error) {
	return c.cache.Get(c.taggedItemKey(key), val)

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

func (c *TagCache) Flush() error {
	for k := range c.names {
		c.ResetTag(c.names[k])
	}

	return nil

}

// taggedItemKey real store key
func (c *TagCache) taggedItemKey(key string) string {
	return EncodeMD5(c.GetNamespace() + key)
}

//////////// tagSet ///////////

// TagID get tag id
func (c *TagCache) TagID(name string) string {
	id, _ := c.cache.Options().Store.Get(c.TagKey(name))
	if len(id) == 0 {
		return c.ResetTag(name)
	}

	return string(id)
}

// TagIDs TagIDs get all tag ids
func (c *TagCache) TagIDs() []string {
	l := len(c.names)
	if l == 0 {
		return nil
	}

	//  排序
	sort.Strings(c.names)

	ids := make([]string, l)
	for i, name := range c.names {
		ids[i] = c.TagID(name)
	}

	return ids
}

// GetNamespace GetNamespace
func (c *TagCache) GetNamespace() string {
	ids := c.TagIDs()
	if len(ids) == 0 {
		return ""
	}
	return strings.Join(ids, "|")
}

// ResetTag ResetTag
func (c *TagCache) ResetTag(name string) string {
	ver := strconv.FormatInt(time.Now().UnixNano(), 36)
	if c.cache.Options().TagTTL > 0 {
		c.cache.Options().Store.Set(c.TagKey(name), []byte(ver), c.cache.Options().TagTTL)
	} else {
		c.cache.Options().Store.Forever(c.TagKey(name), []byte(ver))
	}

	return ver
}

// TagKey TagKey
func (c *TagCache) TagKey(name string) string {
	return fmt.Sprintf("%s.tagid:%s", c.cache.Options().Prefix, name)
}
