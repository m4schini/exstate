package redis

import (
	"context"
	"github.com/m4schini/exstate"
	"time"
)

func (r *redisSrc) CacheString(expires time.Duration, path ...string) (exstate.GetString, exstate.Setter[string]) {
	var getter exstate.GetString
	var setter exstate.Setter[string]

	getter = func() string {
		ctx, cancel := context.WithTimeout(context.Background(), Timeout)
		defer cancel()

		res, err := r.db.Get(ctx, toKey(path...)).Result()
		if err != nil {
			res = ""
		}

		return res
	}

	setter = func(val string) {
		ctx, cancel := context.WithTimeout(context.Background(), Timeout)
		defer cancel()
		r.db.Set(ctx, toKey(path...), val, expires) // 0 expiration == no timeout
	}

	return getter, setter
}

func (r *redisSrc) CacheInt(expires time.Duration, path ...string) (exstate.GetInt, exstate.Setter[int]) {
	var getter exstate.GetInt
	var setter exstate.Setter[int]

	getter = func() int {
		ctx, cancel := context.WithTimeout(context.Background(), Timeout)
		defer cancel()

		res := r.db.Get(ctx, toKey(path...))
		err := res.Err()
		if err != nil {
			return -1
		}

		i, err := res.Int()
		if err != nil {
			return -1
		}

		return i
	}

	setter = func(val int) {
		ctx, cancel := context.WithTimeout(context.Background(), Timeout)
		defer cancel()
		r.db.Set(ctx, toKey(path...), val, expires) // 0 expiration == no timeout
	}

	return getter, setter
}

func (r *redisSrc) CacheFloat(expires time.Duration, path ...string) (exstate.GetFloat, exstate.Setter[float64]) {
	var getter exstate.GetFloat
	var setter exstate.Setter[float64]

	getter = func() float64 {
		ctx, cancel := context.WithTimeout(context.Background(), Timeout)
		defer cancel()

		res := r.db.Get(ctx, toKey(path...))
		if res.Err() != nil {
			return -1.0
		}

		f, err := res.Float64()
		if err != nil {
			return -1
		}

		return f
	}

	setter = func(val float64) {
		ctx, cancel := context.WithTimeout(context.Background(), Timeout)
		defer cancel()
		r.db.Set(ctx, toKey(path...), val, expires) // 0 expiration == no timeout
	}

	return getter, setter
}

func (r *redisSrc) CacheBool(expires time.Duration, path ...string) (exstate.GetBool, exstate.Setter[bool]) {
	var getter exstate.GetBool
	var setter exstate.Setter[bool]

	getter = func() bool {
		ctx, cancel := context.WithTimeout(context.Background(), Timeout)
		defer cancel()

		res := r.db.Get(ctx, toKey(path...))
		if res.Err() != nil {
			return false
		}

		f, err := res.Bool()
		if err != nil {
			return false
		}

		return f
	}

	setter = func(val bool) {
		ctx, cancel := context.WithTimeout(context.Background(), Timeout)
		defer cancel()
		r.db.Set(ctx, toKey(path...), val, expires) // 0 expiration == no timeout
	}

	return getter, setter
}
