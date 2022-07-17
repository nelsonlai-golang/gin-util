package cache

import (
	"time"

	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
)

// cacheStore is a in-memory cache store
// It is not supposed to be exposed outside of this package
var cacheStore = persistence.NewInMemoryStore(time.Minute)

// CacheInTime cache a gin handler function in a given time
func CacheInTime(expire time.Duration, handler gin.HandlerFunc) gin.HandlerFunc {
	return cache.CachePage(cacheStore, expire, handler)
}
