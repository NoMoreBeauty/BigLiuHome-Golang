package auth

import (
	"wxcloudrun-golang/db/model"
)

// UserInterface 用户数据模型接口
type UserInterface interface {
	Login(userKey string) (*model.UserModel, error) // 只用手机号登录
}

// UserInterfaceImp UserInterface的实现对象
type UserInterfaceImp struct{}

// 实现实例
var UserImp UserInterface = &UserInterfaceImp{}
