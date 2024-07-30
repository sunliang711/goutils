package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/sunliang711/goutils/http/server"
)

func main() {
	httpServer := server.NewHttpServer(server.WithPort(9001))

	httpServer.AddMiddlewares(nil)
	httpServer.AddHealthHandler()

	err := httpServer.AddRoutes([]server.Routes{
		{
			GroupPath:        "/blockchain",
			GroupMiddlewares: []gin.HandlerFunc{},
			Handlers: []server.Handler{
				{
					Name:        "blockchain",
					Method:      "GET",
					Path:        "/list",
					Middlewares: []gin.HandlerFunc{},
					Handler:     ListBlockchain,
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}

	httpServer.AddCustomFunc(func(e *gin.Engine) {
		v2 := e.Group("/api/v2")
		v2.GET("/health", func(c *gin.Context) {
			c.JSON(200, "ok_v2")
		})
	})

	httpServer.Start()

	// wait signal
	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt, syscall.SIGTERM)
	<-osSignal

	httpServer.Stop()
}

func ListBlockchain(c *gin.Context) {
	blockchains := []Blockchain{
		{
			BlockchainName: "Ethereum",
			ID:             0,
		},
		{
			BlockchainName: "Bitcoin",
			ID:             1,
		},
	}
	c.JSON(0, blockchains)
}

type Blockchain struct {
	BlockchainName string `json:"blockchain_name"`
	ID             int    `json:"id"`
}
