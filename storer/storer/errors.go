package storer

import "errors"

var (
	_ Storer = &ErrStorer{}
	// ErrStorerNotFound ErrStorerNotFound
	ErrStorerNotFound = errors.New("storer not found, forgot to init register?")
)

// ErrStorer ErrStorer
type ErrStorer struct {
	err error
}

// NewErrStorer NewErrStorer
func NewErrStorer(err error) *ErrStorer {
	return &ErrStorer{err}
}

// Set Set
func (e *ErrStorer) Set(key string, val []byte, ttl int64) error {
	return e.err

}

// Get Get
func (e *ErrStorer) Get(key string) ([]byte, error) {
	return nil, e.err

}

// Expire Expire
func (e *ErrStorer) Expire(key string, ttl int64) error {
	return e.err

}

// Forever Forever
func (e *ErrStorer) Forever(key string, val []byte) error {
	return e.err

}

// Del Del
func (e *ErrStorer) Del(key string) error {
	return e.err

}

// NewWithConf NewWithConf
func (e *ErrStorer) NewWithConf(jsonconf string) error {
	return e.err

}

func (e *ErrStorer) String() string {
	return e.err.Error()
}

func (e *ErrStorer) Error() string {
	return e.err.Error()
}
