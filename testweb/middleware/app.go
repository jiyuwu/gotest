package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/jiyuwu/gotest/testweb/controllers"
)

func Router(server *gin.Engine) {
	new(controllers.AccountHandler).Init(server)
}

// Websocket 路由
func WebsocketInit() {
	controllers.Register("login", controllers.LoginController)
	controllers.Register("heartbeat", controllers.HeartbeatController)
	controllers.Register("ping", controllers.PingController)
}
