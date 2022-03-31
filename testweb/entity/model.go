package entity

import "time"

type ModelPrefix struct {
	Id int64 `json:"id" gorm:"column:Id;primaryKey"`
}

// ModelSuffix 如果只需要创建与修改时间，则用这个
type ModelSuffix struct {
	CreatedAt *time.Time `json:"createdAt,omitempty" gorm:"column:CreateTime;index"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty" gorm:"column:UpdateTime"`
}
