package redis

// Options Options
type Options struct {
	Addr        string
	Passwd      string
	SelectDB    int
	MaxIdle     int
	IdleTimeout int
}

func defaultOptions(opts Options) Options {
	if len(opts.Addr) == 0 {
		opts.Addr = "127.0.0.1:6379"
	}

	if opts.MaxIdle == 0 {
		opts.MaxIdle = 100
	}

	if opts.IdleTimeout == 0 {
		opts.IdleTimeout = 30
	}

	return opts
}
