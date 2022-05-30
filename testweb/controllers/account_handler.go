package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jiyuwu/gotest/testweb/cache"
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
	// 第二步登录验证
	var accDao dao.AccountDAO
	t, err := accDao.GetAccount(req.UserName, common.MD5(req.Password))
	if err != nil {
		result.Code = 2
		result.Msg = "login failure"
		c.JSON(http.StatusOK, result)
		return
	}
	// 存储数据
	err = cache.SetUserTokenInfo(vo.GetUserKey(req.AppId, strconv.FormatInt(t.Id, 10)), t.Token)
	if err != nil {
		result.Code = 3
		result.Msg = "save token failure"
		c.JSON(http.StatusOK, result)
		return
	}

	result.Token = t.Token
	result.Id = t.Id
	result.Code = 1
	result.Msg = "success"
	c.JSON(http.StatusOK, result)
}
