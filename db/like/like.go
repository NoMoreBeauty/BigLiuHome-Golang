package like

import (
	"errors"
	"time"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
	"wxcloudrun-golang/logger"

	"gorm.io/gorm"
)

const tableName = "meal_likes"

// PostLike 插入点赞或取消点赞
func (imp *MealLikeInterfaceImp) PostLike(mealId, userId int32, userName string) (bool, int32, error) {
	cli := db.Get()
	const mod = "点赞DB"

	logger.Info(mod, "开始查询点赞记录", "mealId", mealId, "userId", userId)

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
		// 未找到点赞记录，说明是首次点赞，属于正常业务流程
		logger.Warn(mod, "未找到点赞记录，执行首次点赞插入", "mealId", mealId, "userId", userId, "userName", userName)
		// 创建点赞记录对象
		newLike := &model.MealLikeModel{
			MealId:    mealId,
			UserId:    userId,
			UserName:  userName,
			CreatedAt: time.Now().Unix(),
		}
		err = cli.Table(tableName).Create(newLike).Error
		if err != nil {
			logger.Error(mod, "插入点赞记录失败", "mealId", mealId, "userId", userId, "err", err)
			return false, -1, err
		}
		logger.Info(mod, "点赞记录插入成功", "mealId", mealId, "userId", userId)
		change = 1
	} else if err == nil {
		// 1.2 找到点赞记录，说明是取消点赞，删除记录
		logger.Info(mod, "找到点赞记录，执行取消点赞", "mealId", mealId, "userId", userId)
		// GORM 在删除时，需要通过反射传入的结构体来获取当前表的主键字段（比如哪个是 id，字段类型是什么），以此来生成正确的 SQL
		err = cli.Table(tableName).
			Where("meal_id = ? AND user_id = ?", mealId, userId).
			Delete(&model.MealLikeModel{}).Error
		if err != nil {
			logger.Error(mod, "删除点赞记录失败", "mealId", mealId, "userId", userId, "err", err)
			return false, -1, err
		}
		logger.Info(mod, "取消点赞成功", "mealId", mealId, "userId", userId)
		change = -1
	} else {
		logger.Error(mod, "查询点赞记录时发生数据库错误", "mealId", mealId, "userId", userId, "err", err)
		return false, -1, err
	}

	// 3. 更新 meals 表中的点赞总数（原子更新，使用 GREATEST 防止点赞数为负数）
	// gorm.Expr 是为了适配gorm的语法的，并发安全
	err = cli.Table("meals").
		Where("id = ?", mealId).
		Update("likes_count", gorm.Expr("GREATEST(likes_count + ?, 0)", change)).Error
	if err != nil {
		logger.Error(mod, "更新帖子点赞总数失败", "mealId", mealId, "change", change, "err", err)
		return change == 1, -1, err
	}

	// 4. 查询最新的点赞总数并返回
	// First 只支持将结果扫描进「结构体（Struct）」或「图（Map）」中，不支持直接扫描进 Go 的基础数据类型（如 int32、string）指针中
	var updatedCount int32
	err = cli.Table("meals").Where("id = ?", mealId).Pluck("likes_count", &updatedCount).Error // 查询一个字段的时候用Pluck
	if err != nil {
		logger.Error(mod, "查询帖子最新点赞数失败", "mealId", mealId, "err", err)
		return change == 1, -1, err
	}

	logger.Info(mod, "点赞操作完成", "mealId", mealId, "userId", userId, "isLiked", change == 1, "latestCount", updatedCount)
	return change == 1, updatedCount, nil
}
