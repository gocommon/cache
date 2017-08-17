package cache

import (
	"github.com/gocommon/cache/storer/storer"
)

// Conf Conf
type Conf struct {
	Prefix string
	TTL    int64
	TagTTL int64

	StoreAdapter       string
	StoreAdapterConfig string

	// CodecAdapter       string
	// CodecAdapterConfig string
}

// NewCacheWithConf NewCacheWithConf
func NewCacheWithConf(conf Conf) Cacher {

	opts := &Options{}
	opts.Prefix = conf.Prefix
	opts.TagTTL = conf.TagTTL
	opts.TTL = conf.TTL

	if len(conf.StoreAdapter) > 0 {
		opts.Store = storer.NewWithAdapter(conf.StoreAdapter, conf.StoreAdapterConfig)
	}

	// if len(conf.CodecAdapter) > 0 {
	// 	opts.Codec = codec.NewWithAdapter(conf.CodecAdapter, conf.CodecAdapterConfig)
	// }

	return New(opts)

}
