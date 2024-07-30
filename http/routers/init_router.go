package routers

import "github.com/gin-gonic/gin"

var (
	RoutersNoCheck = make([]func(*gin.RouterGroup), 0)
	RoutersCheck   = make([]func(*gin.RouterGroup, gin.HandlerFunc), 0)
)
