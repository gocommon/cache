package cache

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Version implements Session.
func (p *session) Version() (string, error) {
	space, err := p.getNamespace()
	if err != nil {
		return "", err
	}

	if len(space) == 0 {
		return "", nil
	}
	return EncodeMD5(space), nil
}

// encodeItemKey real store key
func (p *session) encodeItemKey(key string) (enkey, version string, err error) {
	space, err := p.getNamespace()
	if err != nil {
		return "", "", err
	}

	return EncodeMD5(space + key), EncodeMD5(space), nil
}

// getNamespace getNamespace
func (p *session) getNamespace() (string, error) {
	ids, err := p.getTagIDs()
	if err != nil {
		return "", err
	}
	if len(ids) == 0 {
		return "", nil
	}

	return strings.Join(ids, "|"), nil
}

// getTagIDs 取tag对应的值
func (p *session) getTagIDs() ([]string, error) {
	l := len(p.tags)
	if l == 0 {
		return nil, nil
	}

	//  排序
	sort.Strings(p.tags)

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
		return nil, errors.New("store.MGet not align")
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
	return fmt.Sprintf("%s.tagid:%s", p.opts.prefix, tag)
}
