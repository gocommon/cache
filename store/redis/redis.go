package redis

import (
	"context"
	"time"

	redisv8 "github.com/go-redis/redis/v8"
	xtime "github.com/gocommon/cache/v2/pkg/time"
	"github.com/gocommon/cache/v2/store"
)

var _ store.Store = (*Redis)(nil)

type Config struct {
	Username string `dsn:"username"`
	Password string `dsn:"password"`
	Network  string `dsn:"network"`
	Address  string `dsn:"address"`
	DB       int    `dsn:"db"`

	DialTimeout  xtime.Duration `dsn:"query.dial_timeout"`
	WriteTimeout xtime.Duration `dsn:"query.write_timeout"`
	ReadTimeout  xtime.Duration `dsn:"query.read_timeout"`
}

type Redis struct {
	rdb *redisv8.Client
}

func NewRedis(rdb *redisv8.Client) *Redis {
	return &Redis{rdb: rdb}
}

// func (p *Redis) Init(dsn string) error {

// 	d, err := parser.Parse(dsn)
// 	if err != nil {
// 		return err
// 	}

// 	cnf := &Config{}

// 	d.Bind(&cnf)

// 	rdb := redisv8.NewClient(&redisv8.Options{
// 		Network:      cnf.Network,
// 		Addr:         cnf.Address,
// 		Password:     cnf.Password,
// 		DB:           cnf.DB,
// 		DialTimeout:  cnf.DialTimeout.AsDuration(),
// 		WriteTimeout: cnf.WriteTimeout.AsDuration(),
// 		ReadTimeout:  cnf.ReadTimeout.AsDuration(),
// 	})
// 	rdb.AddHook(redisotel.TracingHook{})

// 	err = rdb.Ping(context.TODO()).Err()
// 	if err != nil {
// 		return err
// 	}

// 	p.rdb = rdb

// 	return nil
// }

func (p *Redis) Get(ctx context.Context, key string) ([]byte, error) {
	return p.rdb.Get(ctx, key).Bytes()
}
func (p *Redis) MGet(ctx context.Context, keys []string) ([][]byte, error) {
	res, err := p.rdb.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	list := make([][]byte, len(res))
	for i, v := range res {
		if v == nil {
			list[i] = []byte("")
		} else {
			list[i] = []byte(v.(string))
		}

	}

	return list, nil
}
func (p *Redis) Set(ctx context.Context, key string, val []byte) error {
	_, err := p.rdb.Set(ctx, key, val, redisv8.KeepTTL).Result()
	if err != nil {
		return err
	}
	return nil
}
func (p *Redis) SetEx(ctx context.Context, key string, val []byte, ttl int64) error {
	_, err := p.rdb.SetEX(ctx, key, val, time.Second*time.Duration(ttl)).Result()
	if err != nil {
		return err
	}
	return nil
}
func (p *Redis) Del(ctx context.Context, key string) error {
	_, err := p.rdb.Del(ctx, key).Result()
	return err
}
