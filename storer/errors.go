package storer

import "errors"

var (
	_ Storer = &ErrStorer{}
	// ErrStorerNotFound ErrStorerNotFound
	ErrStorerNotFound = errors.New("storer not found")
)

// ErrStorer ErrStorer
type ErrStorer struct {
	err error
}

func NewErrStorer(err error) *ErrStorer {
	return &ErrStorer{err}
}

// Set Set
func (e *ErrStorer) Set(key string, val string, ttl int64) error {
	return e.err

}

// Get Get
func (e *ErrStorer) Get(key string) (string, error) {
	return "", e.err

}

// Expire Expire
func (e *ErrStorer) Expire(key string, ttl int64) error {
	return e.err

}

// Forever Forever
func (e *ErrStorer) Forever(key string, val string) error {
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
