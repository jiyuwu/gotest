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
		//执行注销之前需要检查注销原token，保证一个用户一个持久化连接

		// 查询到用户生成并更新token
		token := uuid.New()
		ret := databaseConn.Table(acc.TableName()).Where("Id = ?", acc.Id).Updates(map[string]interface{}{"Token": token, "UpdateTime": time.Now().UTC()}).Scan(&acc)
		if ret.Error == nil {
			return &acc, nil
		}
	}

	return nil, e
}
