package cachestore

import (
	"errors"
	"time"
)

const (
	DEFAULT = time.Duration(0)
	FORVER  = time.Duration(-1)
)

var (
	ErrCacheMiss  = errors.New("cache: key not found.")
	ErrNotStored  = errors.New("cache: not stored.")
	ErrNotSupport = errors.New("cache: not support.")
)

type CacheStore interface {
	Get(key string, value interface{}) error
	Set(key string, value interface{}, expire time.Duration) error
	Add(key string, value interface{}, expire time.Duration) error
	Delete(key string) error
}
