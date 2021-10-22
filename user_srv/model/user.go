package model

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID        int       `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"column:add_time"`
	UpdatedAt time.Time `gorm:"column:update_time"`
	DeletedAt gorm.DeletedAt
	IsDeleted bool
}

type User struct {
	BaseModel
	Mobile   string `gorm:"uniqueIndex:uniq_mobile;type:varchar(11);not null"`
	Password string `gorm:"type:varchar(100);not null"`
	NickName string `gorm:"type:varchar(20)"`
	//防止保存的时候出现错误，这里采用指针类型
	Birthday *time.Time `gorm:"type:datetime"`
	Gender   string     `gorm:"default:male;type:varchar(6) comment 'female女，male男'"`
	Role     int        `gorm:"default:1;type:int comment '1表示普通用户，2表示管理员'"`
}
