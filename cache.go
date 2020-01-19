package gin_cache

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type responseCache struct {
	Status int
	Header http.Header
	Data   []byte
}

type cacheWriter struct {
	gin.ResponseWriter
	store  interface{}
	expire time.Duration
	key    string
}
