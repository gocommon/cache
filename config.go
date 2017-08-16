package cache

import (
	"github.com/gocommon/cache/codec"
	"github.com/gocommon/cache/locker"
	"github.com/gocommon/cache/storer"
)

// Conf Conf
type Conf struct {
	Prefix string
	TTL    int64
	TagTTL int64

	StoreAdapter       string
	StoreAdapterConfig string

	CodecAdapter       string
	CodecAdapterConfig string

	UseLocker           bool
	LockerAdapter       string
	LockerAdapterConfig string
}

// NewCacheWithConf NewCacheWithConf
func NewCacheWithConf(conf Conf) Cacher {
	opts := Options{}
	opts.Prefix = conf.Prefix
	opts.TagTTL = conf.TagTTL
	opts.TTL = conf.TTL

	opts.Store = storer.NewWithAdapter(conf.StoreAdapter, conf.StoreAdapterConfig)
	opts.Locker = locker.NewWithAdapter(conf.LockerAdapter, conf.LockerAdapterConfig)
	opts.Codec = codec.NewWithAdapter(conf.CodecAdapter, conf.CodecAdapterConfig)

	opts.UseLocker = conf.UseLocker

	return NewWithOptions(opts)

}
