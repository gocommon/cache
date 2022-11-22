package codec

import (
	"github.com/gocommon/cache/v2/codec/codec"
	_ "github.com/gocommon/cache/v2/codec/gob"
)

var (
	// DefaultCodec DefaultCodec
	DefaultCodec = codec.NewWithAdapter("gob", "")
)
