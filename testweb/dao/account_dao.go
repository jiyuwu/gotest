package dao

import (
	"github.com/jiyuwu/gotest/testweb/entity"
	"gorm.io/gorm"
)

// AccountDAO dao type
type AccountDAO int

// GetAccount  根据用户名和密码查询
func (dao AccountDAO) GetAccount(username, password string) (*entity.AccountEntity, error) {
	var acc entity.AccountEntity

	e := databaseConn.Where("UserName = ? and Password = ?", username, password).First(&acc).Error
	if e == nil {
		return &acc, nil
	}
	if e == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return nil, e
}
