package like

import (
	"errors"
	"time"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"

	"gorm.io/gorm"
)

const tableName = "meal_likes"

// PostLike 插入点赞或取消点赞
func (imp *MealLikeInterfaceImp) PostLike(mealId, userId int32, userName string) (bool, int32, error) {
	cli := db.Get()

	// 1. 查询有没有当前用户对当前帖子的点赞记录
	var like model.MealLikeModel
	// 2. 查询单条记录是否存在（只查询 id 字段提升效率）
	err := cli.Table(tableName).
		Where("meal_id = ? AND user_id = ?", mealId, userId).
		First(&like).Error
	// 这里不能用pluck，因为pluck找不到不会返回error

	// 1.1 如果没有则插入
	var change int32 // 点赞数变更：+1 或 -1
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 创建点赞记录对象
		newLike := &model.MealLikeModel{
			MealId:    mealId,
			UserId:    userId,
			UserName:  userName,
			CreatedAt: time.Now().Unix(),
		}
		err = cli.Table(tableName).Create(newLike).Error
		if err != nil {
			return false, -1, err
		}
		change = 1
	} else if err == nil {
		// 1.2 如果有则说明是取消，删除这一条记录
		err = cli.Table(tableName).
			Where("meal_id = ? AND user_id = ?", mealId, userId).
			Delete(&model.MealLikeModel{}).Error
		// GORM 在删除时，需要通过反射传入的结构体来获取当前表的主键字段（比如哪个是 id，字段类型是什么），以此来生成正确的 SQL
		if err != nil {
			return false, -1, err
		}
		change = -1
	} else {
		return false, -1, err
	}

	// 3. 更新 meals 表中的点赞总数（原子更新，使用 GREATEST 防止点赞数为负数）
	err = cli.Table("meals").
		Where("id = ?", mealId).
		Update("likes_count", gorm.Expr("GREATEST(likes_count + ?, 0)", change)).Error
	// gorm.Expr 是为了适配gorm的语法的，并发安全
	if err != nil {
		return change == 1, -1, err
	}

	// 4. 查询最新的点赞总数并返回
	var updatedCount int32
	err = cli.Table("meals").Where("id = ?", mealId).Pluck("likes_count", &updatedCount).Error // 查询一个字段的时候用Pluck
	if err != nil {
		return change == 1, -1, err
	}
	// First 只支持将结果扫描进「结构体（Struct）」或「图（Map）」中，不支持直接扫描进 Go 的基础数据类型（如 int32、string）指针中

	return change == 1, updatedCount, nil
}
