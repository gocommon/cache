package gob

import (
	"github.com/gocommon/cache/codec"
)

// NewCodecgst NewCodecgst
func NewCodec() codec.Codec {
	return codec.NewGobCodec()
}
