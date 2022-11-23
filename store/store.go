package store

import (
	"context"
	"errors"
)

// Store store
type Store interface {
	Get(ctx context.Context, key string) ([]byte, error)
	MGet(ctx context.Context, keys []string) ([][]byte, error)
	Set(ctx context.Context, key string, val []byte) error
	SetEx(ctx context.Context, key string, val []byte, ttl int64) error
	Del(ctx context.Context, key string) error
	// network:[//[username[:password]@]address[:port][,address[:port]]][/path][?query][#fragment]
	// Init(dsn string) error
}

var (
	// ErrStoreNotFound ErrStoreNotFound
	ErrStoreNotFound = errors.New("store not found, forgot to init register?")
	DefaultErrStore  = NewErrStore(ErrStoreNotFound)
)

var _ Store = (*ErrStore)(nil)

// ErrStore ErrStore
type ErrStore struct {
	err error
}

// NewErrStore NewErrStore
func NewErrStore(err error) *ErrStore {
	return &ErrStore{err}
}

// Set Set
func (e *ErrStore) Set(ctx context.Context, key string, val []byte) error {
	return e.err

}

// Set Set˜
func (e *ErrStore) SetEx(ctx context.Context, key string, val []byte, ttl int64) error {
	return e.err

}

// Get Get
func (e *ErrStore) Get(ctx context.Context, key string) ([]byte, error) {
	return nil, e.err

}

// Get Get
func (e *ErrStore) MGet(ctx context.Context, keys []string) ([][]byte, error) {
	return nil, e.err

}

// Del Del
func (e *ErrStore) Del(ctx context.Context, key string) error {
	return e.err

}

// Infi Infi
func (e *ErrStore) Init(dsn string) error {
	return e.err

}

func (e *ErrStore) String() string {
	return e.err.Error()
}

func (e *ErrStore) Error() string {
	return e.err.Error()
}
