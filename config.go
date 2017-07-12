package cache

import (
	"strings"

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

func NewCacheWithConf(conf Conf) TagCacher {
	opts := &Options{}
	opts.Prefix = opts.Prefix
	opts.TagTTL = opts.TagTTL
	opts.TTL = opts.TTL

}

func StoreAdapter(name, jsonconf string) storer.Storer {
	switch strings.ToLower(name) {

	}
}

func LockerAdapter(name, jsonconf string) locker.Locker {

}

func CodecAdapter(name, jsonconf string) codec.Codec {

}
