package memory

import (
	"context"
	"sync"
	"time"

	"github.com/gocommon/cache/v2/store"
)

var _ store.Store = (*Memory)(nil)

// item 表示内存中的缓存项
type item struct {
	value   []byte
	expires int64 // Unix timestamp, 0表示永不过期
}

// Memory 内存存储实现
type Memory struct {
	mu    sync.RWMutex
	data  map[string]*item
	stop  chan struct{}
	clean time.Duration // 清理间隔
}

// NewMemory 创建新的内存存储实例
func NewMemory() *Memory {
	m := &Memory{
		data:  make(map[string]*item),
		stop:  make(chan struct{}),
		clean: time.Minute, // 每分钟清理一次过期数据
	}

	// 启动清理协程
	go m.cleanup()

	return m
}

// NewMemoryWithCleanupInterval 创建内存存储并指定清理间隔
func NewMemoryWithCleanupInterval(cleanupInterval time.Duration) *Memory {
	m := &Memory{
		data:  make(map[string]*item),
		stop:  make(chan struct{}),
		clean: cleanupInterval,
	}

	// 启动清理协程
	go m.cleanup()

	return m
}

// Close 关闭内存存储，停止清理协程
func (m *Memory) Close() {
	close(m.stop)
}

// cleanup 定期清理过期数据
func (m *Memory) cleanup() {
	ticker := time.NewTicker(m.clean)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.deleteExpired()
		case <-m.stop:
			return
		}
	}
}

// deleteExpired 删除所有过期的数据
func (m *Memory) deleteExpired() {
	now := time.Now().Unix()

	m.mu.Lock()
	defer m.mu.Unlock()

	for key, item := range m.data {
		if item.expires > 0 && item.expires <= now {
			delete(m.data, key)
		}
	}
}

// isExpired 检查item是否过期
func (i *item) isExpired() bool {
	if i.expires == 0 {
		return false
	}
	return i.expires <= time.Now().Unix()
}

// Get 获取指定key的值
func (m *Memory) Get(ctx context.Context, key string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	item, exists := m.data[key]
	if !exists {
		return nil, nil // 不存在返回nil而不是错误
	}

	if item.isExpired() {
		// 虽然有清理协程，但这里也要检查过期
		return nil, nil
	}

	// 返回值的副本，避免外部修改
	value := make([]byte, len(item.value))
	copy(value, item.value)
	return value, nil
}

// MGet 批量获取多个key的值
func (m *Memory) MGet(ctx context.Context, keys []string) ([][]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	results := make([][]byte, len(keys))

	for i, key := range keys {
		item, exists := m.data[key]
		if exists && !item.isExpired() {
			// 返回值的副本
			value := make([]byte, len(item.value))
			copy(value, item.value)
			results[i] = value
		}
		// 如果不存在或过期，results[i]保持为nil
	}

	return results, nil
}

// Set 设置key-value（永久存储）
func (m *Memory) Set(ctx context.Context, key string, val []byte) error {
	return m.set(key, val, 0)
}

// SetEx 设置key-value并指定过期时间
func (m *Memory) SetEx(ctx context.Context, key string, val []byte, ttl int64) error {
	var expires int64
	if ttl > 0 {
		expires = time.Now().Unix() + ttl
	}
	return m.set(key, val, expires)
}

// set 内部设置方法
func (m *Memory) set(key string, val []byte, expires int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 存储值的副本，避免外部修改
	value := make([]byte, len(val))
	copy(value, val)

	m.data[key] = &item{
		value:   value,
		expires: expires,
	}

	return nil
}

// Del 删除指定的key
func (m *Memory) Del(ctx context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.data, key)
	return nil
}
