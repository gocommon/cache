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
