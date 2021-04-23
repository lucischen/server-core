package redis

import (
	goredis "github.com/go-redis/redis"
)

type IRedis interface {
	Use(id int64) (cache *goredis.Client, ok bool)
}

type RedisMap map[int64]*goredis.Client

var _ IRedis = new(RedisMap)

func New() RedisMap {
	return make(RedisMap)
}

type Options struct {
	goredis.Options
}

func (r RedisMap) Add(id int64, o *Options) error {
	if _, ok := r[id]; !ok {
		redisClient := goredis.NewClient(&goredis.Options{
			Addr: o.Addr,
			DB:   o.DB,
		})

		err := redisClient.Ping().Err()
		if err != nil {
			return err
		}

		r[id] = redisClient
	}

	return nil
}

func (r RedisMap) Use(id int64) (*goredis.Client, bool) {
	cache, ok := r[id]
	return cache, ok
}
