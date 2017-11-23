package cache

// Options Options
type Options struct {
	Prefix   string
	TTL      int64 // key 有效期
	TouchTTL int64 // 多少秒内访问，自动续期
	TagTTL   int64 // tagkey 有效期，默认-1，永久，如果想省内容空间，可以设置值

	// CodecAdapter string

	StoreAdapter       string
	StoreAdapterConfig string
}

func defaultOptions(opts Options) Options {
	if opts.TTL == 0 {
		opts.TTL = 7200
	}

	if opts.TagTTL == 0 {
		opts.TagTTL = -1
	}

	if opts.TouchTTL == 0 {
		opts.TouchTTL = 600
	}

	if len(opts.Prefix) == 0 {
		opts.Prefix = "tagcache."
	}

	return opts
}
