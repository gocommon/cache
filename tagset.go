package cache

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

// TagSet TagSet
type TagSet struct {
	opts  *Options
	names []string
}

// NewTagSet NewTagSet
func NewTagSet(opts *Options, names []string) *TagSet {

	t := &TagSet{opts, names}
	return t
}

// Reset tag key
func (s *TagSet) Reset() error {
	for _, name := range s.names {
		s.ResetTag(name)
	}
	return nil
}

// TagID get tag id
func (s *TagSet) TagID(name string) string {
	id, _ := s.opts.Store.Get(s.TagKey(name))
	if len(id) == 0 {
		return s.ResetTag(name)
	}

	return string(id)
}

// TagIDs TagIDs get all tag ids
func (s *TagSet) TagIDs() []string {
	l := len(s.names)
	if l == 0 {
		return nil
	}

	//  排序
	sort.Strings(s.names)

	ids := make([]string, l)
	for i, name := range s.names {
		ids[i] = s.TagID(name)
	}

	return ids
}

// GetNamespace GetNamespace
func (s *TagSet) GetNamespace() string {
	ids := s.TagIDs()
	if len(ids) == 0 {
		return ""
	}
	return strings.Join(ids, "|")
}

// ResetTag ResetTag
func (s *TagSet) ResetTag(name string) string {
	ver := strconv.FormatInt(time.Now().UnixNano(), 36)
	if s.opts.TagTTL > 0 {
		s.opts.Store.Set(s.TagKey(name), []byte(ver), s.opts.TagTTL)
	} else {
		s.opts.Store.Forever(s.TagKey(name), []byte(ver))
	}

	return ver
}

// TagKey TagKey
func (s *TagSet) TagKey(name string) string {
	return fmt.Sprintf("%s:%s:key", s.opts.Prefix, name)
}
