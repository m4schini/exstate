package redis

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/m4schini/exstate"
	"time"
)

func (r *redisSrc) New(path string, expires time.Duration, refresh func() interface{}) (exstate.GetCacheFunc, exstate.SetCacheFunc, error) {
	if r.db == nil {
		return nil, nil, errors.New("redis db not available")
	}

	var getter exstate.GetCacheFunc
	var setter exstate.SetCacheFunc

	setter = func(value interface{}) error {
		ctx, cancel := context.WithTimeout(context.Background(), Timeout)
		defer cancel()
		return r.db.Set(ctx, path, value, expires).Err() // 0 expiration == no timeout
	}

	getter = func() (interface{}, error) {
		var result interface{}

		ctx, cancel := context.WithTimeout(context.Background(), Timeout)
		defer cancel()

		v, err := r.db.Get(ctx, path).Result()
		if err == redis.Nil {
			result = refresh()
			setter(result)
			return result, nil
		}
		if err != nil {
			return nil, err
		}

		return v, nil
	}

	return getter, setter, nil
}
