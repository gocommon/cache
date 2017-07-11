package cache

import (
	"crypto/md5"
	"encoding/hex"
)

// EncodeMD5 EncodeMD5
func EncodeMD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
