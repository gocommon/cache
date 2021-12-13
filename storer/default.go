package storer

import (
	"errors"

	_ "github.com/gocommon/cache/storer/redis"
	"github.com/gocommon/cache/storer/storer"
)

var (
	// DefaultStore DefaultStore
	DefaultStore = storer.NewErrStorer(errors.New("cache store muster be set"))
)
