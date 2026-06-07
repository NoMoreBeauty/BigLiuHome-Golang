package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

// StringArray 自定义类型，实现与 MySQL JSON 字段的自动转换
type StringArray []string

// Scan 实现 sql.Scanner 接口，用于从数据库读出 JSON 时解析回 []string
func (a *StringArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSON value:", value))
	}
	return json.Unmarshal(bytes, a)
}

// Value 实现 driver.Valuer 接口，用于存入数据库时将 []string 转为 JSON 字节
func (a StringArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return "[]", nil
	}
	return json.Marshal(a)
}

// MealModel 三餐帖子数据模型
type MealModel struct {
	Id             int32           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId         int32           `gorm:"column:user_id" json:"user_id"`
	UserName       string          `gorm:"column:user_name" json:"user_name"`
	MealType       string          `gorm:"column:meal_type" json:"meal_type"`
	Images         StringArray     `gorm:"column:images" json:"images"`
	Description    string          `gorm:"column:description" json:"description"`
	LikesCount     int32           `gorm:"column:likes_count" json:"likes_count"`
	CommentsCount  int32           `gorm:"column:comments_count" json:"comments_count"`
	CreatedAt      int64           `gorm:"column:created_at" json:"created_at"`
	Companions     StringArray     `gorm:"column:companions" json:"-"` // 共同用餐的人
	CompanionInfos []CompanionInfo `gorm:"-" json:"companions"`
	IsLiked        bool            `gorm:"column:is_liked;->" json:"is_liked"` // 添加只读的 is_liked 动态计算字段，Insert和Update不受影响
}

// CompanionInfo 因为查询的时候要带上用户的name
type CompanionInfo struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
}

// MealCommentModel 三餐帖子评论模型
type MealCommentModel struct {
	Id          int32  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	MealId      int32  `gorm:"column:meal_id" json:"meal_id"`
	ParentId    int32  `gorm:"column:parent_id" json:"parent_id"`
	UserId      int32  `gorm:"column:user_id" json:"user_id"`
	UserName    string `gorm:"column:user_name" json:"user_name"`
	ReplyToId   int32  `gorm:"column:reply_to_id" json:"reply_to_id"`
	ReplyToName string `gorm:"column:reply_to_name" json:"reply_to_name"`
	Content     string `gorm:"column:content" json:"content"`
	CreatedAt   int64  `gorm:"column:created_at" json:"created_at"`
}

// MealLikeModel 三餐帖子点赞模型
type MealLikeModel struct {
	Id        int32  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	MealId    int32  `gorm:"column:meal_id" json:"meal_id"`
	UserId    int32  `gorm:"column:user_id" json:"user_id"`
	UserName  string `gorm:"column:user_name" json:"user_name"`
	CreatedAt int64  `gorm:"column:created_at" json:"created_at"`
}
