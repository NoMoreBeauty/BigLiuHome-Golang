package comment

import "wxcloudrun-golang/db/model"

// MealCommentInterface 三餐帖子评论数据模型接口
type MealCommentInterface interface {
	GetMealCommentsHandler(mealId int32) ([]*model.MealCommentModel, error)
	PostCommentHandler(mealId int32, userId int32, userName string, parentId int32, replyId int32, replyName string, content string) (*model.MealCommentModel, error)
}

// MealCommentInterfaceImp MealCommentInterface的实现对象
type MealCommentInterfaceImp struct{}

// 实现实例
var MealCommentImp MealCommentInterface = &MealCommentInterfaceImp{}
