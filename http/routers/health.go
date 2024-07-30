package routers

import "github.com/gin-gonic/gin"

func init() {
	RoutersNoCheck = append(RoutersNoCheck, healthRouter)
}

func healthRouter(group *gin.RouterGroup) {
	group.GET("/health", func(c *gin.Context) {
		c.String(0, "ok")
	})
}
