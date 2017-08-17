package codec

import (
	"github.com/gocommon/cache/codec/codec"
	_ "github.com/gocommon/cache/codec/gob"
)

var (
	// DefaultCodec DefaultCodec
	DefaultCodec = codec.NewWithAdapter("gob", "")
)
