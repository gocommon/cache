package codec

import "errors"

var (
	_ Codec = &ErrCodec{}

	// ErrCodeNotFound ErrCodeNotFound
	ErrCodeNotFound = errors.New("codec not found")
)

// ErrCodec ErrCodec
type ErrCodec struct {
	err error
}

// NewErrStorer NewErrStorer
func NewErrStorer(err error) *ErrCodec {
	return &ErrCodec{err}
}

// Encode Encode
func (e *ErrCodec) Encode(v interface{}) (string, error) {
	return "", e.err

}

// Decode Decode
func (e *ErrCodec) Decode(data string, v interface{}) error {
	return e.err

}

// NewWithConf NewWithConf
func (e *ErrCodec) NewWithConf(jsonconf string) error {
	return e.err

}
