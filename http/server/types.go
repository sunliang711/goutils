package server

import (
	"github.com/gin-gonic/gin"
)

type Middleware struct {
	Name    string
	Handler gin.HandlerFunc
}

type Handler struct {
	Name   string
	Method string
	Path   string

	Middlewares []gin.HandlerFunc
	Handler     gin.HandlerFunc
}

type Routes struct {
	GroupPath        string
	GroupMiddlewares []gin.HandlerFunc
	Handlers         []Handler
}

type CustomFunc func(*gin.Engine)
