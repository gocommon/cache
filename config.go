package cache

import (
	"github.com/gocommon/cache/codec/codec"
	"github.com/gocommon/cache/locker/locker"
	"github.com/gocommon/cache/storer/storer"

	_ "github.com/gocommon/cache/locker/etcdv3"
	_ "github.com/gocommon/cache/locker/redis"
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

	// UseLocker           bool
	LockerAdapter       string
	LockerAdapterConfig string
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

	if len(conf.LockerAdapter) > 0 {
		opts.Locker = locker.NewWithAdapter(conf.LockerAdapter, conf.LockerAdapterConfig)
	}

	if len(conf.CodecAdapter) > 0 {
		opts.Codec = codec.NewWithAdapter(conf.CodecAdapter, conf.CodecAdapterConfig)
	}

	// opts.UseLocker = conf.UseLocker

	return New(opts)

}
