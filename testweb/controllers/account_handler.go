package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jiyuwu/gotest/testweb/common"
	"github.com/jiyuwu/gotest/testweb/dao"
	"github.com/jiyuwu/gotest/testweb/vo"
	log "github.com/sirupsen/logrus"
)

type AccountHandler struct {
}

func (ah *AccountHandler) Init(server *gin.Engine) {
	accGroup := server.Group("/account")
	accGroup.POST("/login", ah.login)
}

func (ah *AccountHandler) login(c *gin.Context) {
	// 第一步，验证参数
	req := new(vo.LoginReq)
	err := c.ShouldBindJSON(req)
	if err != nil {
		log.Info("LoginReq req:%+v err:%s", req, err)
		return
	}
	var accDao dao.AccountDAO
	t, err := accDao.GetAccount(req.UserName, common.MD5(req.Password))
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, t)
}
