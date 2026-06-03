package model

import "time"

// User 用户模型
type UserModel struct {
	Id        int32     `gorm:"column:id" json:"id"`
	UserKey   string    `gorm:"column:user_key" json:"user_key"`
	UserName  string    `gorm:"column:user_name" json:"user_name"`
	Role      string    `gorm:"column:role" json:"role"`
	Status    int32     `gorm:"column:status" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}
