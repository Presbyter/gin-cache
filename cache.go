package gin_cache

import (
	"bytes"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/Presbyter/gin-cache/cachestore"
	"github.com/gin-gonic/gin"
)

const (
	CACHEPREFIX = "gin-cache.page"
)

type responseCache struct {
	Status int
	Header http.Header
	Data   []byte
}

type cacheWriter struct {
	gin.ResponseWriter
	store  cachestore.CacheStore
	expire time.Duration
	key    string
}

func (w *cacheWriter) Write(data []byte) (int, error) {
	ret, err := w.ResponseWriter.Write(data)
	if err == nil {
		if w.ResponseWriter.Status() < 300 {
			cache := responseCache{
				Status: w.ResponseWriter.Status(),
				Header: w.ResponseWriter.Header(),
				Data:   data,
			}
			w.store.Set(w.key, cache, w.expire)
		}
	}
	return ret, err
}

func (w *cacheWriter) WriteString(data string) (int, error) {
	ret, err := w.ResponseWriter.WriteString(data)
	if err == nil && w.ResponseWriter.Status() < 300 {
		cache := responseCache{
			Status: w.ResponseWriter.Status(),
			Header: w.ResponseWriter.Header(),
			Data:   []byte(data),
		}
		w.store.Set(w.key, cache, w.expire)
	}
	return ret, err
}

func newCacheWriter(store cachestore.CacheStore, expire time.Duration, writer gin.ResponseWriter, key string) *cacheWriter {
	return &cacheWriter{
		ResponseWriter: writer,
		store:          store,
		expire:         expire,
		key:            key,
	}
}

func createKey(str ...string) string {
	var buffer bytes.Buffer
	buffer.WriteString(CACHEPREFIX)
	buffer.WriteString(":")
	for _, item := range str {
		buffer.WriteString(url.QueryEscape(item))
		buffer.WriteString(":")
	}
	return strings.Trim(buffer.String(), ":")
}

func CachePage(store cachestore.CacheStore, expire time.Duration, handle gin.HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var cache responseCache
		key := createKey(ctx.Request.URL.RequestURI())
		if err := store.Get(key, &cache); err != nil {
			writer := newCacheWriter(store, expire, ctx.Writer, key)
			ctx.Writer = writer
			handle(ctx)

			if ctx.IsAborted() {
				store.Delete(key)
			}
		} else {
			ctx.Writer.WriteHeader(cache.Status)
			for k, vals := range cache.Header {
				for _, v := range vals {
					ctx.Writer.Header().Set(k, v)
				}
			}
			ctx.Writer.Write(cache.Data)
		}
	}
}

func CachePageAtomic(store cachestore.CacheStore, expire time.Duration, handle gin.HandlerFunc) gin.HandlerFunc {
	var m sync.Mutex
	p := CachePage(store, expire, handle)
	return func(ctx *gin.Context) {
		m.Lock()
		defer m.Unlock()
		p(ctx)
	}
}
