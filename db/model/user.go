package model

// User 用户模型
type UserModel struct {
	Id        int32  `gorm:"column:id" json:"id"`
	UserKey   string `gorm:"column:user_key" json:"user_key"`
	UserName  string `gorm:"column:user_name" json:"user_name"`
	Role      string `gorm:"column:role" json:"role"`
	Status    int32  `gorm:"column:status" json:"status"`
	CreatedAt int64  `gorm:"column:created_at" json:"created_at"` // 数据表建的时候建成unix时间戳了
	UpdatedAt int64  `gorm:"column:updated_at" json:"updated_at"`
}
