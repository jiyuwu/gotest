package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jiyuwu/gotest/testweb/controllers"
	"github.com/jiyuwu/gotest/testweb/dao"
	"github.com/jiyuwu/gotest/testweb/middleware"
)

func main() {
	gin.SetMode(gin.ReleaseMode) //线上环境

	go controllers.Manager.Start()

	//健康检查  1.没有token的直接失效 2.没有超时的失效

	//websocket测试
	r := gin.Default()
	r.Use(middleware.Cors())
	middleware.Router(r)
	r.GET("/ws", controllers.WsHandler)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	// 初始化数据库
	dao.InitDatabaseConn()
	r.Run(":9962")
}
