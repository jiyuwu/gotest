package entity

// AccountEntity 账号实体
type AccountEntity struct {
	ModelPrefix
	Ip         string `json:"ip,omitempty" gorm:"column:Ip;type:varchar(150); index:idx_account_ip;"`
	UserName   string `json:"userName,omitempty" gorm:"column:UserName;type:varchar(30); index:idx_account_name;"`
	Password   string `json:"password,omitempty" gorm:"column:Password;type:varchar(150);"`
	DeviceType int32  `json:"deviceType,omitempty" gorm:"column:DeviceType;"`
	State      int32  `json:"state,omitempty" gorm:"column:State;"`
	LoginTime  string `json:"loginTime,omitempty" gorm:"column:LoginTime;"`
	Token      string `json:"token,omitempty" gorm:"column:Token;type:varchar(255);"`
	ModelSuffix
}

//TableName 映射表名
func (acc *AccountEntity) TableName() string {
	return "User"
}
