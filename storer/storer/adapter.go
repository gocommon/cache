package storer

import (
	"strings"
)

// Adapter Adapter
var Adapter = map[string]Storer{}

// Register Register
func Register(name string, s Storer) {
	name = strings.ToLower(name)
	Adapter[name] = s
}

// NewWithAdapter NewWithAdapter
func NewWithAdapter(name, jsonconf string) Storer {
	name = strings.ToLower(name)
	s, ok := Adapter[name]
	if !ok {
		return NewErrStorer(ErrStorerNotFound)
	}

	err := s.NewWithConf(jsonconf)
	if err != nil {
		return NewErrStorer(err)
	}

	return s
}
