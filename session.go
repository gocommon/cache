package cache

import (
	"bytes"
	"context"
	"time"
)

var _ Session = (*session)(nil)

type session struct {
	ctx  context.Context
	tags []string
	c    *Cache
}

// genKey 统一处理生成key, tag
func (p *session) genKey(key string) (string, error) {
	if len(p.tags) > 0 {
		k, err := p.encodeItemKey(key)
		if err != nil {
			return "", err
		}
		return p.c.keyWithPrefix(k), nil
	}

	return p.c.keyWithPrefix(key), nil

}

func (p *session) Get(key string, val interface{}) (has bool, err error) {

	rk, err := p.genKey(key)
	if err != nil {
		return false, err
	}

	src, err := p.c.store.Get(p.ctx, rk)
	if err != nil {
		return false, err
	}

	if len(src) == 0 {
		return false, nil
	}

	d, unix := p.c.splitUnix(src)

	// near expire
	if unix > 0 && unix+p.c.opts.TTL-time.Now().Unix() < p.c.opts.TouchTTL {
		unix := time.Now().Unix()
		d = p.c.joinUnix(d, unix)
		p.c.store.SetEx(p.ctx, p.c.keyWithPrefix(key), d, p.c.opts.TTL)
	}

	if bytes.Contains(d, EmptyValue) {
		// SetNil(val)
		return true, nil
	}

	return true, p.c.codec.Decode(d, val)
}

func (p *session) Set(key string, val interface{}) error {
	d := EmptyValue

	if !IsNil(val) {
		var err error
		d, err = p.c.codec.Encode(val)
		if err != nil {
			return err
		}

	}

	// add unix to the end @
	unix := time.Now().Unix()
	d = p.c.joinUnix(d, unix)

	rk, err := p.genKey(key)
	if err != nil {
		return err
	}

	return p.c.store.SetEx(p.ctx, rk, d, p.c.opts.TTL)
}

func (p *session) Del(key string) error {
	rk, err := p.genKey(key)
	if err != nil {
		return err
	}
	return p.c.store.Del(p.ctx, rk)
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
