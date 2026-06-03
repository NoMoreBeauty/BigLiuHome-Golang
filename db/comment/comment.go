package comment

import (
	"time"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

const tableName = "meal_comments"

// GetMealCommentsHandler 根据id查询帖子评论
func (imp *MealCommentInterfaceImp) GetMealCommentsHandler(mealId int32) ([]*model.MealCommentModel, error) {
	cli := db.Get()
	var meals []*model.MealCommentModel

	// 构建基础查询，使用 Select 注入 EXISTS 子查询来动态计算 is_liked 字段
	query := cli.Table(tableName).
		Where("meal_id = ?", mealId).
		Order("created_at ASC")

	// 执行查询
	err := query.Find(&meals).Error

	return meals, err
}

// PostCommentHandler 插入评论
func (imp *MealCommentInterfaceImp) PostCommentHandler(mealId int32, userId int32, userName string, parentId int32, replyId int32, replyName string, content string) (*model.MealCommentModel, error) {
	cli := db.Get()

	var comment = &model.MealCommentModel{
		MealId:      mealId,
		ParentId:    parentId,
		UserId:      userId,
		UserName:    userName,
		ReplyToId:   replyId,
		ReplyToName: replyName,
		Content:     content,
		CreatedAt:   time.Now().Unix(),
	}
	err := cli.Table(tableName).Create(comment).Error
	return comment, err
}
