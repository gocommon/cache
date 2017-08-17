package etcdv3

import (
	"context"
	"encoding/json"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"github.com/gocommon/cache/locker"
)

var _ locker.Locker = &Locker{}

// Locker Locker
type Locker struct {
	cli *clientv3.Client
}

// NewLocker NewLocker
func NewLocker(config clientv3.Config) (*Locker, error) {
	c, err := clientv3.New(config)
	if err != nil {
		return nil, err
	}

	return &Locker{c}, nil
}

// NewWithConf NewWithConf json Config
func (l *Locker) NewWithConf(jsonconf string) error {
	// clientv3.Config{Endpoints: endpoints}

	var conf clientv3.Config

	err := json.Unmarshal([]byte(jsonconf), &conf)
	if err != nil {
		return err
	}

	c, err := clientv3.New(conf)
	if err != nil {
		return err
	}

	l.cli = c

	return nil

}

// NewLocker NewLocker
func (l *Locker) NewLocker(key string) locker.Funcer {
	s, err := concurrency.NewSession(l.cli)
	if err != nil {
		return locker.NewErrFuncer(err)
	}

	return NewMutex(concurrency.NewMutex(s, key))
}

// Mutex Mutex
type Mutex struct {
	m *concurrency.Mutex
}

// NewMutex NewMutex
func NewMutex(s *concurrency.Mutex) *Mutex {
	return &Mutex{s}
}

// Lock Lock
func (m *Mutex) Lock() error {
	return m.m.Lock(context.TODO())
}

// Unlock Unlock
func (m *Mutex) Unlock() error {
	return m.m.Unlock(context.TODO())

}

func init() {
	locker.Register("etcdv3", &Locker{})
}
