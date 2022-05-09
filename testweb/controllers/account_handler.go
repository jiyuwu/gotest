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
	result := new(vo.LoginResponse)
	err := c.ShouldBindJSON(req)
	if err != nil {
		result.Code = 0
		result.Msg = "login req error"
		c.JSON(http.StatusOK, result)
		log.Info("login req:%+v err:%s", req, err)
		return
	}
	var accDao dao.AccountDAO
	t, err := accDao.GetAccount(req.UserName, common.MD5(req.Password))
	// 第二步登录成功后生成token并更新
	if err != nil {
		result.Code = 2
		result.Msg = "login failure"
		c.JSON(http.StatusOK, result)
		return
	}
	result.Token = t.Token
	result.Id = t.Id
	result.Code = 1
	result.Msg = "success"
	c.JSON(http.StatusOK, result)
}
