package codec

import (
	"bytes"
	"encoding/gob"
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
func (c *GobCodec) Encode(v interface{}) (string, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(v)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// Decode Decode
func (c *GobCodec) Decode(data string, v interface{}) error {
	r := bytes.NewReader([]byte(data))
	dec := gob.NewDecoder(r)

	return dec.Decode(v)
}
