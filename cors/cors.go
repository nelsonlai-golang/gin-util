package cors

import (
	c "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Config = c.Config

// Setup cors policy
func Setup(r *gin.Engine, config Config) {
	r.Use(c.New(c.Config(config)))
}

// QuickSetup cors policy (allow all origins)
func QuickSetup(r *gin.Engine) {
	Setup(r, c.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	})
}
