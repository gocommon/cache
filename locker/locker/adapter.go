package locker

import (
	"strings"
)

// Adapter Adapter
var Adapter = map[string]Locker{}

// Register Register
func Register(name string, s Locker) {
	name = strings.ToLower(name)
	Adapter[name] = s
}

// NewWithAdapter NewWithAdapter
func NewWithAdapter(name, jsonconf string) Locker {
	name = strings.ToLower(name)
	s, ok := Adapter[name]
	if !ok {
		return NewErrLocker(ErrLockerNotFound)
	}

	err := s.NewWithConf(jsonconf)
	if err != nil {
		return NewErrLocker(err)
	}

	return s
}
