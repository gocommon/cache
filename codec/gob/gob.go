package gob

import (
	"bytes"
	"encoding/gob"
)

// Codec Codec
type Codec struct{}

// NewCodec NewCodec
func NewCodec() *Codec {
	return &Codec{}
}

// Encode Encode
func (c *Codec) Encode(v interface{}) (string, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(v)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// Decode Decode
func (c *Codec) Decode(data string, v interface{}) error {
	r := bytes.NewReader([]byte(data))
	dec := gob.NewDecoder(r)

	return dec.Decode(v)
}
