package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/m4schini/exstate"
	"strings"
	"time"
)

const (
	Timeout = 5 * time.Second
)

type redisSrc struct {
	db *redis.Client
}

func New(addr, password string, db int) (*redisSrc, error) {
	src := new(redisSrc)
	src.db = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	status := src.db.Ping(ctx)
	cancel()

	return src, status.Err()
}

func (r *redisSrc) String(keys ...string) (exstate.GetString, exstate.Setter[string]) {
	var getter exstate.GetString
	var setter exstate.Setter[string]

	getter = func() string {
		ctx, cancel := context.WithTimeout(context.Background(), Timeout)
		defer cancel()
		return r.db.Get(ctx, toKey(keys...)).Val()
	}

	setter = func(val string) {
		ctx, cancel := context.WithTimeout(context.Background(), Timeout)
		defer cancel()
		r.db.Set(ctx, toKey(keys...), val, 0) // 0 expiration == no timeout
	}

	return getter, setter
}

func (r *redisSrc) Int(keys ...string) (exstate.GetInt, exstate.Setter[int]) {
	var getter exstate.GetInt
	var setter exstate.Setter[int]

	getter = func() int {
		ctx, cancel := context.WithTimeout(context.Background(), Timeout)
		defer cancel()
		i, err := r.db.Get(ctx, toKey(keys...)).Int()
		if err != nil {
			i = -1
		}
		return i
	}

	setter = func(val int) {
		ctx, cancel := context.WithTimeout(context.Background(), Timeout)
		defer cancel()
		r.db.Set(ctx, toKey(keys...), val, Timeout)
	}

	return getter, setter
}

func (r *redisSrc) Float(keys ...string) (exstate.GetFloat, exstate.Setter[float64]) {
	var getter exstate.GetFloat
	var setter exstate.Setter[float64]

	getter = func() float64 {
		ctx, cancel := context.WithTimeout(context.Background(), Timeout)
		defer cancel()
		i, err := r.db.Get(ctx, toKey(keys...)).Float64()
		if err != nil {
			i = -1
		}
		return i
	}

	setter = func(val float64) {
		ctx, cancel := context.WithTimeout(context.Background(), Timeout)
		defer cancel()
		r.db.Set(ctx, toKey(keys...), val, Timeout)
	}

	return getter, setter
}

func (r *redisSrc) Bool(keys ...string) (exstate.GetBool, exstate.Setter[bool]) {
	var getter exstate.GetBool
	var setter exstate.Setter[bool]

	getter = func() bool {
		ctx, cancel := context.WithTimeout(context.Background(), Timeout)
		defer cancel()
		b, err := r.db.Get(ctx, toKey(keys...)).Bool()
		if err != nil {
			b = false
		}
		return b
	}

	setter = func(val bool) {
		ctx, cancel := context.WithTimeout(context.Background(), Timeout)
		defer cancel()
		r.db.Set(ctx, toKey(keys...), val, Timeout)
	}

	return getter, setter
}

func (r *redisSrc) Set(keys ...string) (exstate.SetAdd, exstate.SetGet, exstate.SetRemove, exstate.SetContains) {
	var adder exstate.SetAdd
	var getter exstate.SetGet
	var container exstate.SetContains
	var remover exstate.SetRemove

	adder = func(value interface{}) {
		ctx, cancel := context.WithTimeout(context.Background(), Timeout)
		defer cancel()

		r.db.SAdd(ctx, toKey(keys...), value)
	}

	getter = func() []string {
		ctx, cancel := context.WithTimeout(context.Background(), Timeout)
		defer cancel()

		r, err := r.db.SMembers(ctx, toKey(keys...)).Result()
		if err != nil {
			r = make([]string, 0)
		}

		return r
	}

	container = func(value interface{}) bool {
		ctx, cancel := context.WithTimeout(context.Background(), Timeout)
		defer cancel()

		return r.db.SIsMember(ctx, toKey(keys...), value).Val()
	}

	remover = func(value interface{}) {
		ctx, cancel := context.WithTimeout(context.Background(), Timeout)
		defer cancel()

		r.db.SRem(ctx, toKey(keys...), value)
	}

	return adder, getter, remover, container
}

func (r *redisSrc) Close() {
	r.db.Close()
}

func toKey(keys ...string) string {
	return strings.Join(keys, ".")
}
