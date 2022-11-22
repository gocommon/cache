package cache

import (
	"fmt"
	"strings"
)

// Adapter Adapter
var Adapter = map[string]Storer{}

// Register Register
func Register(name string, s Storer) {
	name = strings.ToLower(name)

	if _, ok := Adapter[name]; ok {
		panic(fmt.Sprintf("store exists:%s", name))
	}
}

// InitStore InitStore
func InitStore(name, jsonconf string) Storer {
	name = strings.ToLower(name)
	s, ok := Adapter[name]
	if !ok {
		return NewErrStorer(ErrStorerNotFound)
	}

	err := s.Init(jsonconf)
	if err != nil {
		return NewErrStorer(err)
	}

	return s
}
