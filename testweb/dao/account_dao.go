package dao

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jiyuwu/gotest/testweb/common"
	"github.com/jiyuwu/gotest/testweb/entity"
	"gorm.io/gorm"
)

// AccountDAO dao type
type AccountDAO int

// GetAccount  根据用户名和密码查询
func (dao AccountDAO) GetAccount(username, password string) (*entity.AccountEntity, error) {
	var acc entity.AccountEntity

	e := databaseConn.Table(acc.TableName()).Where("UserName = ? and Password = ?", username, password).First(&acc).Error
	if e == gorm.ErrRecordNotFound {
		return nil, e
	}
	if e == nil {
		// 查询到用户生成并更新token,数据库存储的token其实没什么意义。因为要模仿微信多端登录，考虑效率问题可以移除更新。
		token := uuid.New()
		ret := databaseConn.Table(acc.TableName()).Where("Id = ?", acc.Id).Updates(map[string]interface{}{"Token": token, "UpdateTime": time.Now().UTC()}).Scan(&acc)
		if ret.Error == nil {
			return &acc, nil
		}
	}

	return nil, e
}

// CreateAccount 创建账号
func (dao AccountDAO) CreateAccount(username, password string) (*entity.AccountEntity, error) {
	var acc entity.AccountEntity

	e := databaseConn.Table(acc.TableName()).Where("UserName = ? and Password = ?", username, password).First(&acc).Error
	if e == gorm.ErrRecordNotFound {
		acc.UserName = username
		acc.Password = password
		acc.Token = uuid.New().String()
		acc.LoginTime = time.Now().Format("2006-01-02 15:04:05")
		acc.DeviceType = 1
		acc.State = 2
		acc.Ip = common.GetServerIp()
		e = databaseConn.Table(acc.TableName()).Create(&acc).First(&acc).Error
		if e != nil {
			return nil, e
		} else {
			return &acc, nil
		}
	}
	return nil, errors.New("register err!")
}
