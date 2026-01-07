package cache_test

import (
	"context"
	"testing"
	"time"

	"github.com/gocommon/cache/v2"
)

// UserData 示例数据结构
type UserData struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestMemoryCache_Example(t *testing.T) {
	ctx := context.Background()

	// 创建使用内存存储的缓存实例
	c := cache.New(
		cache.WithMemoryStore(),
		cache.WithTTL(300), // 5分钟过期
	)

	// 创建带标签的会话
	session := c.Tags(ctx, "user", "123")

	// 准备测试数据
	user := UserData{
		ID:   123,
		Name: "John Doe",
		Age:  30,
	}

	// 设置缓存
	err := session.Set("profile", user)
	if err != nil {
		t.Fatalf("Failed to set cache: %v", err)
	}

	// 获取缓存
	var retrieved UserData
	has, err := session.Get("profile", &retrieved)
	if err != nil {
		t.Fatalf("Failed to get cache: %v", err)
	}

	if !has {
		t.Fatal("Cache should exist")
	}

	if retrieved != user {
		t.Errorf("Expected %+v, got %+v", user, retrieved)
	}

	t.Logf("Successfully cached and retrieved: %+v", retrieved)
}

func TestMemoryCache_WithExpiration(t *testing.T) {
	ctx := context.Background()

	// 创建内存缓存，设置短过期时间用于测试
	c := cache.New(
		cache.WithMemoryStore(),
		cache.WithTTL(1), // 1秒过期
	)

	session := c.Tags(ctx, "test", "expire")

	// 设置缓存
	err := session.Set("temp_data", "temporary value")
	if err != nil {
		t.Fatalf("Failed to set cache: %v", err)
	}

	// 立即获取，应该存在
	var value string
	has, err := session.Get("temp_data", &value)
	if err != nil {
		t.Fatalf("Failed to get cache: %v", err)
	}
	if !has || value != "temporary value" {
		t.Fatal("Cache should exist immediately after setting")
	}

	// 等待过期
	time.Sleep(1100 * time.Millisecond)

	// 再次获取，应该已过期
	has, err = session.Get("temp_data", &value)
	if err != nil {
		t.Fatalf("Failed to get cache: %v", err)
	}
	if has {
		t.Fatal("Cache should have expired")
	}

	t.Log("Cache expiration test passed")
}

func TestMemoryCache_TagFlush(t *testing.T) {
	ctx := context.Background()

	c := cache.New(cache.WithMemoryStore())

	// 创建两个不同标签的会话
	userSession := c.Tags(ctx, "user", "123")
	postSession := c.Tags(ctx, "post", "456")

	// 设置用户数据
	err := userSession.Set("profile", "user profile data")
	if err != nil {
		t.Fatalf("Failed to set user cache: %v", err)
	}

	// 设置帖子数据
	err = postSession.Set("content", "post content data")
	if err != nil {
		t.Fatalf("Failed to set post cache: %v", err)
	}

	// 验证两个缓存都存在
	var userData, postData string
	has, _ := userSession.Get("profile", &userData)
	if !has || userData != "user profile data" {
		t.Fatal("User cache should exist")
	}

	has, _ = postSession.Get("content", &postData)
	if !has || postData != "post content data" {
		t.Fatal("Post cache should exist")
	}

	// 刷新用户标签
	err = userSession.Flush()
	if err != nil {
		t.Fatalf("Failed to flush user cache: %v", err)
	}

	// 用户缓存应该被清除
	has, _ = userSession.Get("profile", &userData)
	if has {
		t.Fatal("User cache should be flushed")
	}

	// 帖子缓存应该仍然存在
	has, _ = postSession.Get("content", &postData)
	if !has || postData != "post content data" {
		t.Fatal("Post cache should still exist after user flush")
	}

	t.Log("Tag-based cache flushing test passed")
}
