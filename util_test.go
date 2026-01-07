package cache

import (
	"testing"
)

// BenchmarkEncodeMD5 MD5哈希性能测试
func BenchmarkEncodeMD5(b *testing.B) {
	testData := []string{
		"short",
		"this_is_a_medium_length_string_for_testing",
		"this_is_a_very_long_string_that_should_be_used_for_performance_testing_of_hash_functions_and_should_be_long_enough_to_see_any_performance_differences_between_different_hashing_algorithms",
		"user:123:profile",
		"cache:key:with:multiple:segments",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, data := range testData {
			EncodeMD5(data)
		}
	}
}

// BenchmarkEncodeHash xxHash性能测试
func BenchmarkEncodeHash(b *testing.B) {
	testData := []string{
		"short",
		"this_is_a_medium_length_string_for_testing",
		"this_is_a_very_long_string_that_should_be_used_for_performance_testing_of_hash_functions_and_should_be_long_enough_to_see_any_performance_differences_between_different_hashing_algorithms",
		"user:123:profile",
		"cache:key:with:multiple:segments",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, data := range testData {
			EncodeHash(data)
		}
	}
}

// TestEncodeHashCorrectness 验证哈希函数的正确性
func TestEncodeHashCorrectness(t *testing.T) {
	testCases := []string{
		"hello",
		"world",
		"cache",
		"test_string",
		"user:123",
		"",
		"special_chars_!@#$%^&*()",
	}

	for _, input := range testCases {
		md5Result := EncodeMD5(input)
		hashResult := EncodeHash(input)

		// 结果不应该为空
		if md5Result == "" {
			t.Errorf("EncodeMD5(%q) returned empty string", input)
		}
		if hashResult == "" {
			t.Errorf("EncodeHash(%q) returned empty string", input)
		}

		// 相同输入应该总是产生相同输出
		if EncodeMD5(input) != md5Result {
			t.Errorf("EncodeMD5(%q) is not deterministic", input)
		}
		if EncodeHash(input) != hashResult {
			t.Errorf("EncodeHash(%q) is not deterministic", input)
		}

		t.Logf("Input: %q, MD5: %s, xxHash: %s", input, md5Result, hashResult)
	}
}
