package main

import (
	gin_cache "github.com/Presbyter/gin-cache"
	"github.com/Presbyter/gin-cache/cachestore"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(os.Stdout)

	g := gin.Default()

	r := g.Use(TestMw1(), TestMw2())
	r.GET("/timestamp", do())
	store := cachestore.NewRedisStore([]string{"127.0.0.1:6379"}, "", 5*time.Second)
	r.GET("/cache_timestamp", gin_cache.CachePageAtomic(store, 2*time.Second, do()))

	g.Run(":3000")
}

func TestMw1() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		now := time.Now()
		logrus.Debugf("DEBUG1 \ttimestamp: %d", now.Unix())
		ctx.Set("test1", &now)
		ctx.Next()
	}
}

func TestMw2() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		now := time.Now()
		logrus.Debugf("DEBUG2 \ttimestamp: %d", now.Unix())
		ctx.Set("test2", &now)
		ctx.Next()
	}
}

func do() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var str1, str2 string
		if v, ok := ctx.Get("test1"); ok {
			str1 = v.(*time.Time).Format(time.RFC3339)
		}
		if v, ok := ctx.Get("test2"); ok {
			str2 = v.(*time.Time).Format(time.RFC3339)
		}

		ctx.JSON(200, gin.H{
			"msg":   "ok",
			"time1": str1,
			"time2": str2,
			"time3": time.Now().Format(time.RFC3339),
		})
	}
}
