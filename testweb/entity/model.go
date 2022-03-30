package entity

import "time"

type ModelPrefix struct {
	ID int64 `gorm:"column:Id;primaryKey"`
}

// ModelSuffix 如果只需要创建与修改时间，则用这个
type ModelSuffix struct {
	CreatedAt *time.Time `gorm:"column:CreateTime;index"`
	UpdatedAt *time.Time `gorm:"column:UpdateTime"`
}
