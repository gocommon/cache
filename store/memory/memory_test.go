package memory

import (
	"context"
	"testing"
	"time"
)

func TestMemoryStore_BasicOperations(t *testing.T) {
	ctx := context.Background()
	store := NewMemory()
	defer store.Close()

	// Test Set and Get
	key := "test_key"
	value := []byte("test_value")

	// Set value
	err := store.Set(ctx, key, value)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Get value
	retrieved, err := store.Get(ctx, key)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if string(retrieved) != string(value) {
		t.Errorf("Expected %s, got %s", string(value), string(retrieved))
	}
}

func TestMemoryStore_SetEx(t *testing.T) {
	ctx := context.Background()
	store := NewMemory()
	defer store.Close()

	key := "expire_key"
	value := []byte("expire_value")
	ttl := int64(1) // 1 second

	// Set with expiration
	err := store.SetEx(ctx, key, value, ttl)
	if err != nil {
		t.Fatalf("SetEx failed: %v", err)
	}

	// Should exist immediately
	retrieved, err := store.Get(ctx, key)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if string(retrieved) != string(value) {
		t.Errorf("Expected %s, got %s", string(value), string(retrieved))
	}

	// Wait for expiration
	time.Sleep(1100 * time.Millisecond)

	// Should be expired
	retrieved, err = store.Get(ctx, key)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if retrieved != nil {
		t.Errorf("Expected nil (expired), got %s", string(retrieved))
	}
}

func TestMemoryStore_MGet(t *testing.T) {
	ctx := context.Background()
	store := NewMemory()
	defer store.Close()

	// Set multiple values
	data := map[string][]byte{
		"key1": []byte("value1"),
		"key2": []byte("value2"),
		"key3": []byte("value3"),
	}

	for k, v := range data {
		err := store.Set(ctx, k, v)
		if err != nil {
			t.Fatalf("Set failed for %s: %v", k, err)
		}
	}

	// Test MGet
	keys := []string{"key1", "key2", "key3", "nonexistent"}
	results, err := store.MGet(ctx, keys)
	if err != nil {
		t.Fatalf("MGet failed: %v", err)
	}

	if len(results) != len(keys) {
		t.Fatalf("Expected %d results, got %d", len(keys), len(results))
	}

	// Check existing keys
	for i, key := range keys[:3] {
		expected := data[key]
		if string(results[i]) != string(expected) {
			t.Errorf("Key %s: expected %s, got %s", key, string(expected), string(results[i]))
		}
	}

	// Check nonexistent key
	if results[3] != nil {
		t.Errorf("Expected nil for nonexistent key, got %s", string(results[3]))
	}
}

func TestMemoryStore_Del(t *testing.T) {
	ctx := context.Background()
	store := NewMemory()
	defer store.Close()

	key := "delete_key"
	value := []byte("delete_value")

	// Set value
	err := store.Set(ctx, key, value)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Verify exists
	retrieved, err := store.Get(ctx, key)
	if err != nil || retrieved == nil {
		t.Fatalf("Key should exist before deletion")
	}

	// Delete
	err = store.Del(ctx, key)
	if err != nil {
		t.Fatalf("Del failed: %v", err)
	}

	// Verify deleted
	retrieved, err = store.Get(ctx, key)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if retrieved != nil {
		t.Errorf("Expected nil after deletion, got %s", string(retrieved))
	}
}

func TestMemoryStore_NonExistentKey(t *testing.T) {
	ctx := context.Background()
	store := NewMemory()
	defer store.Close()

	// Get non-existent key
	retrieved, err := store.Get(ctx, "nonexistent")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if retrieved != nil {
		t.Errorf("Expected nil for non-existent key, got %s", string(retrieved))
	}
}

func TestMemoryStore_DataIsolation(t *testing.T) {
	ctx := context.Background()
	store := NewMemory()
	defer store.Close()

	key := "isolation_test"
	originalValue := []byte("original")

	// Set value
	err := store.Set(ctx, key, originalValue)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Get value and modify the returned slice
	retrieved, err := store.Get(ctx, key)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	// Modify the retrieved slice (this should not affect stored data)
	retrieved[0] = 'X'

	// Get again - should still be original value
	retrievedAgain, err := store.Get(ctx, key)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if string(retrievedAgain) != string(originalValue) {
		t.Errorf("Data isolation failed: expected %s, got %s", string(originalValue), string(retrievedAgain))
	}
}
