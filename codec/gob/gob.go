package codec

import (
	"bytes"
	"encoding/gob"

	"github.com/gocommon/cache/codec/codec"
)

// GobCodec GobCodec
type GobCodec struct{}

// NewGobCodec NewGobCodec
func NewGobCodec() *GobCodec {
	return &GobCodec{}
}

// NewWithConf NewWithConf
func (c *GobCodec) NewWithConf(jsonconf string) error {
	return nil
}

// Encode Encode
func (c *GobCodec) Encode(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(v)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Decode Decode
func (c *GobCodec) Decode(data []byte, v interface{}) error {
	r := bytes.NewReader(data)
	dec := gob.NewDecoder(r)

	return dec.Decode(v)
}

func init() {
	codec.Register("gob", &GobCodec{})
}
