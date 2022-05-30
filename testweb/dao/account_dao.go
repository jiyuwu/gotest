package dao

import (
	"time"

	"github.com/google/uuid"
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
