package main

import (
	"fmt"
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jiyuwu/gotest/testweb/cache"
	"github.com/jiyuwu/gotest/testweb/controllers"
	"github.com/jiyuwu/gotest/testweb/dao"
	"github.com/jiyuwu/gotest/testweb/logs"
	"github.com/jiyuwu/gotest/testweb/middleware"
	"github.com/jiyuwu/gotest/testweb/rpcinterface/rpc_server"
	"github.com/jiyuwu/gotest/testweb/task"
	"github.com/spf13/viper"
)

func main() {
	initConfig()
	initFile()
	initRedis()
	// 全局日志，非针对请求
	logs.InitLogger()
	go controllers.Start()

	r := gin.Default()
	r.Use(middleware.Cors())
	middleware.Router(r)       //普通接口
	middleware.WebsocketInit() // websocket接口

	// 定时任务
	task.Init()
	//服务器监控初始化
	task.ServerInit()

	r.GET("/ws", controllers.WsHandler)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	// 初始化数据库
	dao.InitDatabaseConn()
	//开启rpc时，存ip+rpc_port到cache，方便集群发消息。
	controllers.SetServer()
	go rpc_server.Init()
	httpPort := viper.GetString("app.httpPort")
	r.Run(":" + httpPort)
}

// 初始化日志
func initFile() {
	// Disable Console Color, you don't need console color when writing the logs to file.
	gin.DisableConsoleColor()

	// Logging to a file.
	logFile := viper.GetString("app.ginLogFile")
	f, _ := os.Create(logFile)
	gin.DefaultWriter = io.MultiWriter(f)
}

func initConfig() {
	viper.SetConfigName("config/app")
	viper.AddConfigPath(".") // 添加搜索路径

	err := viper.ReadInConfig()
	if err != nil {
		logs.Error(fmt.Sprintf("Fatal error config file: %s \n", err))
	}

	// fmt.Println("config app:", viper.Get("app"))
	// fmt.Println("config redis:", viper.Get("redis"))
}

func initRedis() {
	cache.NewClient()
}
