package cache

import (
	"bytes"
	"context"
	"time"
)

// EmptyValue EmptyValue
var EmptyValue = []byte("##empty- -!##")

var _ Session = (*session)(nil)

type session struct {
	ctx  context.Context
	tags []string
	opts *Options
}

// genKey 统一处理生成key, tag
func (p *session) genKey(key string) (string, error) {
	if len(p.tags) > 0 {
		k, err := p.encodeItemKey(key)
		if err != nil {
			return "", err
		}
		return p.keyWithPrefix(k), nil
	}

	return p.keyWithPrefix(key), nil

}

func (c *session) keyWithPrefix(key string) string {
	return c.opts.prefix + key
}

func (p *session) Get(key string, val interface{}) (has bool, err error) {

	rk, err := p.genKey(key)
	if err != nil {
		return false, err
	}

	src, err := p.opts.store.Get(p.ctx, rk)
	if err != nil {
		return false, err
	}

	if len(src) == 0 {
		return false, nil
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
		return true, nil
	}

	return true, p.opts.codec.Decode(d, val)
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

	rk, err := p.genKey(key)
	if err != nil {
		return err
	}

	return p.opts.store.SetEx(p.ctx, rk, d, p.opts.ttl)
}

func (p *session) Del(key string) error {
	rk, err := p.genKey(key)
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
	idx := len(src) - 9

	flag := src[idx : idx+1]
	if idx < 0 || flag[0] != '@' {
		return src, 0
	}

	return src[0:idx], int64(BytesToUint64(src[idx+1:]))

}

func joinUnix(data []byte, unix int64) []byte {
	buf := bytes.NewBuffer(data)
	buf.WriteByte(byte('@'))
	buf.Write(Uint64ToBytes(uint64(unix)))

	return buf.Bytes()
}
