package codec

import (
	"strings"
)

// Adapter Adapter
var Adapter = map[string]Codec{}

// Register Register
func Register(name string, s Codec) {
	name = strings.ToLower(name)
	Adapter[name] = s
}

// NewWithAdapter NewWithAdapter
func NewWithAdapter(name, jsonconf string) Codec {
	name = strings.ToLower(name)
	s, ok := Adapter[name]
	if !ok {
		return NewErrStorer(ErrCodeNotFound)
	}

	err := s.NewWithConf(jsonconf)
	if err != nil {
		return NewErrStorer(err)
	}

	return s
}
