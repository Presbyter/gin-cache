package cachestore

import (
	"github.com/Presbyter/gin-cache/utils"
	"github.com/patrickmn/go-cache"
	"time"
)

type memoryCache struct {
	cache.Cache
}

func (m *memoryCache) Get(key string, value interface{}) error {
	v, ok := m.Cache.Get(key)
	if !ok {
		return ErrCacheMiss
	}

	return utils.Deserializer(&utils.GobEncode{}, v.([]byte), value)
}

func (m *memoryCache) Set(key string, value interface{}, expire time.Duration) error {
	data, err := utils.Serializer(&utils.GobEncode{}, value)
	if err != nil {
		return err
	}
	m.Cache.Set(key, data, expire)
	return nil
}

func (m *memoryCache) Add(key string, value interface{}, expire time.Duration) error {
	data, err := utils.Serializer(&utils.GobEncode{}, value)
	if err != nil {
		return err
	}
	return m.Cache.Add(key, data, expire)
}

func (m *memoryCache) Delete(key string) error {
	m.Cache.Delete(key)
	return nil
}

func NewMemoryCache(expire time.Duration) CacheStore {
	c := cache.New(expire, 10*time.Minute)
	return &memoryCache{
		Cache: *c,
	}
}
