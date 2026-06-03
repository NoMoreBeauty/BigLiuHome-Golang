package meal

import "wxcloudrun-golang/db/model"

// MealInterface 三餐帖子数据模型接口
type MealInterface interface {
	GetMeals(page, size, userId int32, date, mealType string) ([]*model.MealModel, error)
	GetMealById(id, userId int32) (*model.MealModel, error)
	PostMeals(userId int32, userName, mealType string, images []byte, description string) error
}

// MealLikeInterface 三餐帖子点赞数据模型接口
type MealLikeInterface interface{}

// MealInterfaceImp MealInterface的实现对象
type MealInterfaceImp struct{}

// MealLikeInterfaceImp MealLikeInterface的实现对象
type MealLikeInterfaceImp struct{}

// 实现实例
var MealImp MealInterface = &MealInterfaceImp{}
var MealLikeImp MealLikeInterface = &MealLikeInterfaceImp{}
