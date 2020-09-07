package routes

import (
	"github.com/gin-gonic/gin"
)

// Root zzz
func Root(r *gin.RouterGroup) *gin.RouterGroup {
	root := r.Group("")

	root.GET("", func(c *gin.Context) {
		c.Redirect(307, "ping")
	})

	root.GET("ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	root.GET("pong", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ping",
		})
	})

	return root
}
