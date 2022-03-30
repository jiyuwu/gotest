package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/jiyuwu/gotest/testweb/controllers"
)

func Router(server *gin.Engine) {
	new(controllers.AccountHandler).Init(server)
}
