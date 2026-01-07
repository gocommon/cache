package gob

import (
	"testing"
)

// TestData 测试用数据结构
type TestData struct {
	Name  string
	Value int
	List  []string
}

func TestGobCodec_EncodeDecode(t *testing.T) {
	codec := NewGobCodec()

	// 测试数据
	original := TestData{
		Name:  "test",
		Value: 42,
		List:  []string{"a", "b", "c"},
	}

	// 编码
	data, err := codec.Encode(original)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	if len(data) == 0 {
		t.Fatal("Encoded data is empty")
	}

	// 解码
	var decoded TestData
	err = codec.Decode(data, &decoded)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	// 验证结果
	if decoded.Name != original.Name || decoded.Value != original.Value || len(decoded.List) != len(original.List) {
		t.Errorf("Expected %+v, got %+v", original, decoded)
	}
	for i, v := range original.List {
		if decoded.List[i] != v {
			t.Errorf("List mismatch at index %d: expected %s, got %s", i, v, decoded.List[i])
		}
	}
}

func TestGobCodec_EncodeNil(t *testing.T) {
	codec := NewGobCodec()

	// Gob不能编码nil值，这是一个已知限制
	// 我们测试编码nil时的错误处理
	_, err := codec.Encode(nil)
	if err == nil {
		t.Fatal("Expected error when encoding nil value")
	}

	t.Log("Gob codec correctly rejects nil values")
}

func TestGobCodec_DecodeInvalidData(t *testing.T) {
	codec := NewGobCodec()

	invalidData := []byte("invalid gob data")

	var decoded TestData
	err := codec.Decode(invalidData, &decoded)
	if err == nil {
		t.Fatal("Expected decode error for invalid data")
	}
}

func TestGobCodec_ComplexTypes(t *testing.T) {
	codec := NewGobCodec()

	// 测试简单复杂类型 - Gob对interface{}有限制
	original := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	data, err := codec.Encode(original)
	if err != nil {
		t.Fatalf("Encode complex type failed: %v", err)
	}

	var decoded map[string]string
	err = codec.Decode(data, &decoded)
	if err != nil {
		t.Fatalf("Decode complex type failed: %v", err)
	}

	// 验证字段
	if decoded["key1"] != "value1" || decoded["key2"] != "value2" {
		t.Errorf("Map decode mismatch: %+v", decoded)
	}
}
