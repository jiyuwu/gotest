package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jiyuwu/gotest/testweb/controllers"
)

func main() {
	fmt.Println("1111111")

	//websocket测试
	router := gin.Default()
	router.GET("/wx", controllers.WsHandle)
}
