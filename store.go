package cache

import (
	"context"
	"errors"
)

// Storer store
type Storer interface {
	Get(ctx context.Context, key string) ([]byte, error)
	MGet(ctx context.Context, keys []string) ([][]byte, error)
	Set(ctx context.Context, key string, val []byte) error
	SetEx(ctx context.Context, key string, val []byte, ttl int64) error
	Del(ctx context.Context, key string) error
	// network:[//[username[:password]@]address[:port][,address[:port]]][/path][?query][#fragment]
	Init(dsn string) error
}

var (
	// ErrStorerNotFound ErrStorerNotFound
	ErrStorerNotFound = errors.New("storer not found, forgot to init register?")
)

var _ Storer = (*ErrStorer)(nil)

// ErrStorer ErrStorer
type ErrStorer struct {
	err error
}

// NewErrStorer NewErrStorer
func NewErrStorer(err error) *ErrStorer {
	return &ErrStorer{err}
}

// Set Set
func (e *ErrStorer) Set(ctx context.Context, key string, val []byte) error {
	return e.err

}

// Set Set˜
func (e *ErrStorer) SetEx(ctx context.Context, key string, val []byte, ttl int64) error {
	return e.err

}

// Get Get
func (e *ErrStorer) Get(ctx context.Context, key string) ([]byte, error) {
	return nil, e.err

}

// Get Get
func (e *ErrStorer) MGet(ctx context.Context, keys []string) ([][]byte, error) {
	return nil, e.err

}

// Del Del
func (e *ErrStorer) Del(ctx context.Context, key string) error {
	return e.err

}

// Infi Infi
func (e *ErrStorer) Init(dsn string) error {
	return e.err

}

func (e *ErrStorer) String() string {
	return e.err.Error()
}

func (e *ErrStorer) Error() string {
	return e.err.Error()
}
