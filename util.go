package cache

import (
	"crypto/md5"
	"encoding/hex"
	"reflect"
)

// EncodeMD5 EncodeMD5
func EncodeMD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// SetNil SetNil
func SetNil(i interface{}) {
	if IsNil(i) {
		return
	}
	v := reflect.ValueOf(i)
	v.Elem().Set(reflect.Zero(v.Elem().Type()))
}

// IsNil IsNil
func IsNil(i interface{}) bool {
	return reflect.ValueOf(i).IsNil()
}

// Uint64ToBytes Uint64ToBytes bigEndian
func Uint64ToBytes(v uint64) []byte {
	b := make([]byte, 8)

	b[0] = byte(v >> 56)
	b[1] = byte(v >> 48)
	b[2] = byte(v >> 40)
	b[3] = byte(v >> 32)
	b[4] = byte(v >> 24)
	b[5] = byte(v >> 16)
	b[6] = byte(v >> 8)
	b[7] = byte(v)

	return b
}

// BytesToUint64 BytesToUint64 bigEndian
func BytesToUint64(b []byte) uint64 {
	_ = b[7] // bounds check hint to compiler; see golang.org/issue/14808
	return uint64(b[7]) | uint64(b[6])<<8 | uint64(b[5])<<16 | uint64(b[4])<<24 |
		uint64(b[3])<<32 | uint64(b[2])<<40 | uint64(b[1])<<48 | uint64(b[0])<<56
}
