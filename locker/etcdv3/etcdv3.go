package etcdv3

import (
	"encoding/json"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"github.com/gocommon/cache/locker/locker"
	"golang.org/x/net/context"
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
	conf.Endpoints = []string{"127.0.0.1:2379"}

	if len(jsonconf) > 0 {
		err := json.Unmarshal([]byte(jsonconf), &conf)
		if err != nil {
			return err
		}
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

	m, err := NewMutex(l.cli, key)
	if err != nil {
		return locker.NewErrFuncer(err)
	}

	return m
}

// Mutex Mutex
type Mutex struct {
	m   *concurrency.Mutex
	s   *concurrency.Session
	ctx context.Context
}

// NewMutex NewMutex
func NewMutex(cli *clientv3.Client, key string) (*Mutex, error) {
	s, err := concurrency.NewSession(cli)
	if err != nil {
		return nil, err
	}

	m := concurrency.NewMutex(s, "/cache-locker/"+key)

	return &Mutex{m: m, s: s, ctx: cli.Ctx()}, nil
}

// Lock Lock
func (m *Mutex) Lock() error {
	return m.m.Lock(m.ctx)
}

// Unlock Unlock
func (m *Mutex) Unlock() error {
	err := m.m.Unlock(m.ctx)
	m.s.Close()
	return err

}

func init() {
	locker.Register("etcdv3", &Locker{})
}
