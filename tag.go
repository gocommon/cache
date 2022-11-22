package cache

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

// encodeItemKey real store key
func (p *session) encodeItemKey(key string) (string, error) {
	space, err := p.getNamespace()
	if err != nil {
		return "", err
	}
	return EncodeMD5(space + key), nil
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

	vals, err := p.c.store.MGet(p.ctx, p.tags)
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
	if p.c.opts.TagTTL > 0 {
		err := p.c.store.SetEx(p.ctx, p.newTagKey(tag), []byte(ver), p.c.opts.TagTTL)
		if err != nil {
			return "", err
		}
	} else {
		err := p.c.store.Set(p.ctx, p.newTagKey(tag), []byte(ver))
		if err != nil {
			return "", err
		}
	}

	return ver, nil
}

// TagKey 拼接tagkey,添加前缀
func (p *session) newTagKey(tag string) string {
	return fmt.Sprintf("%s.tagid:%s", p.c.opts.Prefix, tag)
}