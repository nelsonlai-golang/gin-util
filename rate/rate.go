package rate

import (
	"time"

	"github.com/gin-gonic/gin"
	limiter "github.com/ulule/limiter/v3"
	limiter_gin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

// This store is used to store the rate limiters
// It's not be exposed to the outside
var store = memory.NewStore()

// SetupGlobal setup the global rate limiter
func SetupGlobal(r *gin.Engine, d time.Duration, limit int64) {
	r.Use(_getLimiterMiddleware(d, limit))
}

// GET setup gin `GET` router with rate limiter
func GET(r *gin.Engine, path string, h gin.HandlerFunc, d time.Duration, limit int64) gin.IRoutes {
	return r.GET(path, _getLimiterMiddleware(d, limit), h)
}

// POST setup gin `POST` router with rate limiter
func POST(r *gin.Engine, path string, h gin.HandlerFunc, d time.Duration, limit int64) gin.IRoutes {
	return r.POST(path, _getLimiterMiddleware(d, limit), h)
}

// PUT setup gin `PUT` router with rate limiter
func PUT(r *gin.Engine, path string, h gin.HandlerFunc, d time.Duration, limit int64) gin.IRoutes {
	return r.PUT(path, _getLimiterMiddleware(d, limit), h)
}

// DELETE setup gin `DELETE` router with rate limiter
func DELETE(r *gin.Engine, path string, h gin.HandlerFunc, d time.Duration, limit int64) gin.IRoutes {
	return r.DELETE(path, _getLimiterMiddleware(d, limit), h)
}

func _getLimiterMiddleware(d time.Duration, limit int64) gin.HandlerFunc {
	return limiter_gin.NewMiddleware(limiter.New(store, limiter.Rate{Period: d, Limit: limit}))
}
