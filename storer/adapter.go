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

// StoreAdapter StoreAdapter
func StoreAdapter(name, jsonconf string) Storer {
	name = strings.ToLower(name)
	if s, ok := Adapter[name]; ok {
		// @todo
	}
}
