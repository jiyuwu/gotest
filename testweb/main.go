package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jiyuwu/gotest/testweb/controllers"
	"github.com/jiyuwu/gotest/testweb/middleware"
)

func main() {
	gin.SetMode(gin.ReleaseMode) //线上环境

	go controllers.Manager.Start()
	//websocket测试
	r := gin.Default()
	r.Use(middleware.Cors())
	r.GET("/ws", controllers.WsHandler)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run(":80")
}
