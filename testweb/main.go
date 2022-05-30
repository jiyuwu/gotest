package main

import (
	"fmt"
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jiyuwu/gotest/testweb/cache"
	"github.com/jiyuwu/gotest/testweb/controllers"
	"github.com/jiyuwu/gotest/testweb/dao"
	"github.com/jiyuwu/gotest/testweb/middleware"
	"github.com/jiyuwu/gotest/testweb/task"
	"github.com/spf13/viper"
)

func main() {
	gin.SetMode(gin.ReleaseMode) //线上环境
	initConfig()
	initFile()
	initRedis()
	go controllers.Start()

	//健康检查  1.没有token的直接失效 2.没有超时的失效

	r := gin.Default()
	r.Use(middleware.Cors())
	middleware.Router(r)       //普通接口
	middleware.WebsocketInit() // websocket接口

	// 定时任务
	task.Init()

	r.GET("/ws", controllers.WsHandler)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	// 初始化数据库
	dao.InitDatabaseConn()

	httpPort := viper.GetString("app.httpPort")
	r.Run(":" + httpPort)
}

// 初始化日志
func initFile() {
	// Disable Console Color, you don't need console color when writing the logs to file.
	gin.DisableConsoleColor()

	// Logging to a file.
	logFile := viper.GetString("app.logFile")
	f, _ := os.Create(logFile)
	gin.DefaultWriter = io.MultiWriter(f)
}

func initConfig() {
	viper.SetConfigName("config/app")
	viper.AddConfigPath(".") // 添加搜索路径

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	fmt.Println("config app:", viper.Get("app"))
	fmt.Println("config redis:", viper.Get("redis"))

}

func initRedis() {
	cache.NewClient()
}
