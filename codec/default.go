package codec

import (
	"github.com/gocommon/cache/codec/gob"
)

var (
	DefaultCodec = gob.NewCodec()
)
