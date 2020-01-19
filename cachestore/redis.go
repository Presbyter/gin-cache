package cachestore

import (
	"time"

	"github.com/Presbyter/gin-cache/utils"
	"github.com/go-redis/redis"
)

type redisStore struct {
	rds           redis.UniversalClient
	defaultExpire time.Duration
}

func (r *redisStore) Get(key string, value interface{}) error {
	data, err := r.rds.Get(key).Bytes()
	if err != nil {
		return err
	}
	return utils.Deserializer(&utils.GobEncode{}, data, value)
}

func (r *redisStore) Set(key string, value interface{}, expire time.Duration) error {
	data, err := utils.Serializer(&utils.GobEncode{}, value)
	if err != nil {
		return err
	}
	return r.rds.Set(key, data, r.getExpire(expire)).Err()
}

func (r *redisStore) Add(key string, value interface{}, expire time.Duration) error {
	if n, _ := r.rds.Exists(key).Result(); n > 0 {
		return ErrNotStored
	}

	data, err := utils.Serializer(&utils.GobEncode{}, value)
	if err != nil {
		return err
	}
	return r.rds.Set(key, data, r.getExpire(expire)).Err()
}

func (r *redisStore) Delete(key string) error {
	return r.rds.Del(key).Err()
}

func (r *redisStore) getExpire(expire time.Duration) time.Duration {
	switch expire {
	case DEFAULT:
		return r.defaultExpire
	case FORVER:
		return time.Duration(0)
	default:
		return expire
	}
}

func NewRedisStore(addrs []string, passwd string, expire time.Duration) CacheStore {
	uc := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    addrs,
		Password: passwd,
	})
	if err := uc.Ping().Err(); err != nil {
		panic(err)
	}
	return &redisStore{rds: uc, defaultExpire: expire}
}

func NewRedisStoreWithClient(uc redis.UniversalClient, expire time.Duration) CacheStore {
	if err := uc.Ping().Err(); err != nil {
		panic(err)
	}
	return &redisStore{rds: uc, defaultExpire: expire}
}
