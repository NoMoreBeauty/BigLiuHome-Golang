package auth

import (
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

const tableName = "users"

// Login 用户登录
func (imp *UserInterfaceImp) Login(userKey string) (*model.UserModel, error) {
	cli := db.Get()
	var user = new(model.UserModel)
	err := cli.Table(tableName).Where("user_key = ?", userKey).First(user).Error
	return user, err
}

// GetMembers 获取所有家庭成员
func (imp *UserInterfaceImp) GetMembers() ([]*model.UserModel, error) {
	cli := db.Get()
	var userList []*model.UserModel
	err := cli.Table(tableName).Select("id, user_name").Find(&userList).Error
	return userList, err
}
