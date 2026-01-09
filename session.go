package cache

import (
	"bytes"
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	// EmptyValue 表示空值的标记
	EmptyValueStr = "##empty- -!##"

	// 时间戳分隔符长度 (1字节 '@' + 8字节 uint64)
	TimestampSuffixLength = 9

	// 时间戳分隔符
	TimestampSeparator = '@'

	// TagKey分隔符
	TagKeySeparator = ".tagid:"
)

var EmptyValue = []byte(EmptyValueStr)

var _ Session = (*session)(nil)

type session struct {
	ctx        context.Context
	tags       []string
	opts       *Options
	tagsSorted bool
}

func NewSession(ctx context.Context, tags []string, opts *Options) *session {

	tagsSorted := false
	if len(tags) > 0 {
		sort.Strings(tags)
		tagsSorted = true
	}

	return &session{ctx: ctx, tags: tags, opts: opts, tagsSorted: tagsSorted}
}

func (p *session) KeyVersion(key string) (string, string, error) {
	k, version, err := p.genKey(key)
	if err != nil {
		return "", "", err
	}
	return k, version, nil
}

// genKey 统一处理生成key, tag
func (p *session) genKey(key string) (enkey, version string, err error) {
	if len(p.tags) > 0 {
		k, version, err := p.encodeItemKey(key)
		if err != nil {
			return "", "", err
		}
		return p.keyWithPrefix(k), version, nil
	}

	return p.keyWithPrefix(key), "", nil

}

func (c *session) keyWithPrefix(key string) string {
	return c.opts.prefix + key
}

func (p *session) Get(key string, val interface{}) (has bool, err error) {
	has, _, err = p.GetWithVersion(key, val)
	return has, err
}

// GetWithVersion implements Session.
func (p *session) GetWithVersion(key string, val interface{}) (has bool, version string, err error) {

	var rk string

	rk, version, err = p.genKey(key)
	if err != nil {
		return false, version, err
	}

	src, err := p.opts.store.Get(p.ctx, rk)
	if err != nil {
		return false, version, err
	}

	if len(src) == 0 {
		return false, version, nil
	}

	d, unix := splitUnix(src)

	// near expire
	if unix > 0 && unix+p.opts.ttl-time.Now().Unix() < p.opts.touchTTL {
		unix := time.Now().Unix()
		d = joinUnix(d, unix)
		p.opts.store.SetEx(p.ctx, rk, d, p.opts.ttl)
	}

	if bytes.Contains(d, EmptyValue) {
		// SetNil(val)
		return true, version, nil
	}

	return true, version, p.opts.codec.Decode(d, val)
}

func (p *session) Set(key string, val interface{}) error {
	d := EmptyValue

	if !IsNil(val) {
		var err error
		d, err = p.opts.codec.Encode(val)
		if err != nil {
			return err
		}

	}

	// add unix to the end @
	unix := time.Now().Unix()
	d = joinUnix(d, unix)

	rk, _, err := p.genKey(key)
	if err != nil {
		return err
	}

	return p.opts.store.SetEx(p.ctx, rk, d, p.opts.ttl)
}

func (p *session) Del(key string) error {
	rk, _, err := p.genKey(key)
	if err != nil {
		return err
	}
	return p.opts.store.Del(p.ctx, rk)
}

func (p *session) Flush() error {
	if len(p.tags) == 0 {
		return nil
	}

	for k := range p.tags {
		p.setTag(p.tags[k])
	}

	return nil
}

func splitUnix(src []byte) (data []byte, unix int64) {
	idx := len(src) - TimestampSuffixLength

	flag := src[idx : idx+1]
	if idx < 0 || flag[0] != TimestampSeparator {
		return src, 0
	}

	return src[0:idx], int64(BytesToUint64(src[idx+1:]))

}

func joinUnix(data []byte, unix int64) []byte {
	buf := bytes.NewBuffer(data)
	buf.WriteByte(TimestampSeparator)
	buf.Write(Uint64ToBytes(uint64(unix)))

	return buf.Bytes()
}

// Version implements Session.
func (p *session) Version() (string, error) {
	space, err := p.getNamespace()
	if err != nil {
		return "", err
	}

	if len(space) == 0 {
		return "", nil
	}
	return EncodeHash(space), nil
}

// encodeItemKey real store key
func (p *session) encodeItemKey(key string) (enkey, version string, err error) {
	space, err := p.getNamespace()
	if err != nil {
		return "", "", err
	}

	hash := EncodeHash(space)
	return key + "." + hash, hash, nil
}

// getNamespace getNamespace
func (p *session) getNamespace() (string, error) {
	ids, err := p.getOrCreateTagIDs()
	if err != nil {
		return "", err
	}
	if len(ids) == 0 {
		return "", nil
	}

	namespace := strings.Join(ids, "|")
	if namespace == "" {
		return "", NewCacheError(ErrCodeTagError, ErrMsgNamespaceEmpty)
	}

	return namespace, nil
}

// getOrCreateTagIDs 取tag对应的值
func (p *session) getOrCreateTagIDs() ([]string, error) {
	l := len(p.tags)
	if l == 0 {
		return nil, nil
	}

	//  排序
	if !p.tagsSorted {
		sort.Strings(p.tags)
		p.tagsSorted = true
	}

	ids := make([]string, l)

	getTags := make([]string, len(p.tags))
	for k, v := range p.tags {
		getTags[k] = p.newTagKey(v)
	}

	vals, err := p.opts.store.MGet(p.ctx, getTags)
	if err != nil {
		return nil, err
	}

	if len(vals) != l {
		return nil, NewCacheError(ErrCodeStoreError, ErrMsgStoreMisaligned)
	}

	for i, val := range vals {
		if len(val) == 0 {
			tid, err := p.setTag(p.tags[i])
			if err != nil {
				return nil, err
			}
			ids[i] = tid
		} else {
			ids[i] = string(val)
		}
	}

	return ids, nil
}

// setTag 更新tag的值
func (p *session) setTag(tag string) (string, error) {
	ver := strconv.FormatInt(time.Now().UnixNano(), 36)
	if p.opts.tagTTL > 0 {
		err := p.opts.store.SetEx(p.ctx, p.newTagKey(tag), []byte(ver), p.opts.tagTTL)
		if err != nil {
			return "", err
		}
	} else {
		err := p.opts.store.Set(p.ctx, p.newTagKey(tag), []byte(ver))
		if err != nil {
			return "", err
		}
	}

	return ver, nil
}

// TagKey 拼接tagkey,添加前缀
func (p *session) newTagKey(tag string) string {
	return fmt.Sprintf("%s%s%s", p.opts.prefix, TagKeySeparator, tag)
}
