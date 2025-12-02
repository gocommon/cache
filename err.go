package cache

import "context"

// ErrCache 实现 Cacher，用于在初始化失败或不希望真正访问存储时，统一返回错误。
//
// 使用场景：
//
//	c := cache.NewErrCache(err)
//	s := c.Tags(ctx, "user", "1")
//	// 后续对 s 的任何调用都会返回同一个 err
type ErrCache struct {
	err error
}

var _ Cacher = (*ErrCache)(nil)

// NewErrCache 创建一个带固定错误的 Cacher。
func NewErrCache(err error) *ErrCache {
	return &ErrCache{err: err}
}

// Tags 实现 Cacher 接口，返回一个始终返回错误的 Session。
func (e *ErrCache) Tags(ctx context.Context, tags ...string) Session {
	return &errSession{err: e.err}
}

// errSession 实现 Session，所有方法都返回同一个错误。
type errSession struct {
	err error
}

var _ Session = (*errSession)(nil)

// Get 始终返回 false 和错误。
func (s *errSession) Get(key string, val interface{}) (bool, error) {
	return false, s.err
}

// GetWithVersion 始终返回 false、空版本和错误。
func (s *errSession) GetWithVersion(key string, val interface{}) (bool, string, error) {
	return false, "", s.err
}

// Set 始终返回错误。
func (s *errSession) Set(key string, val interface{}) error {
	return s.err
}

// Del 始终返回错误。
func (s *errSession) Del(key string) error {
	return s.err
}

// Version 始终返回空字符串和错误。
func (s *errSession) Version() (string, error) {
	return "", s.err
}

// Flush 始终返回错误。
func (s *errSession) Flush() error {
	return s.err
}
