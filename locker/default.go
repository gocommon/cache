package locker

import (
	_ "github.com/gocommon/cache/locker/etcdv3"
	"github.com/gocommon/cache/locker/locker"
)

var (
	// DefaultLocker DefaultLocker
	DefaultLocker = locker.NewWithAdapter("etcdv3", "")
)
