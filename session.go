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
	_, err := p.SetWithVersion(key, val)
	return err
}

// SetWithVersion implements Session.
func (p *session) SetWithVersion(key string, val interface{}) (version string, err error) {

	d := EmptyValue

	if !IsNil(val) {
		var err error
		d, err = p.opts.codec.Encode(val)
		if err != nil {
			return "", err
		}

	}

	// add unix to the end @
	unix := time.Now().Unix()
	d = joinUnix(d, unix)

	rk, version, err := p.genKey(key)
	if err != nil {
		return "", err
	}

	return version, p.opts.store.SetEx(p.ctx, rk, d, p.opts.ttl)
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
