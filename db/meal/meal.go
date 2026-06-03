package meal

import (
	"time"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

const tableName = "meals"

// GetMeals 查询三餐帖子信息
func (imp *MealInterfaceImp) GetMeals(page, size, userId int32, date, mealType string) ([]*model.MealModel, error) {
	cli := db.Get()
	var meals []*model.MealModel

	// 构建基础查询，使用 Select 注入 EXISTS 子查询来动态计算 is_liked 字段
	query := cli.Table(tableName).
		Select("meals.*, EXISTS(SELECT 1 FROM meal_likes l WHERE l.meal_id = meals.id AND l.user_id = ?) AS is_liked", userId)

	// 如果传入了 date 或 mealType，说明进入“回忆模块”进行筛选
	if date != "" {
		// 假设传入格式为 "YYYY-MM-DD" (例如 "2026-06-03")
		// 将其解析为当天的 Unix 时间戳范围（秒级戳）
		t, err := time.ParseInLocation("2006-01-02", date, time.Local)
		if err == nil {
			startTime := t.Unix()
			endTime := startTime + 86400 // 加一天的秒数
			query = query.Where("created_at >= ? AND created_at < ?", startTime, endTime)
		}
	}

	if mealType != "" {
		query = query.Where("meal_type = ?", mealType)
	}

	// 计算分页偏移量偏移
	offset := (page - 1) * size

	// 执行查询
	err := query.Order("created_at DESC").
		Limit(int(size)).
		Offset(int(offset)).
		Find(&meals).Error

	return meals, err
}

// PostMeals 上传三餐帖子
func (imp *MealInterfaceImp) PostMeals(userId int32, userName, mealType string, images []string, description string) error {
	cli := db.Get()

	var meal = &model.MealModel{
		UserId:      userId,
		UserName:    userName,
		MealType:    mealType,
		Images:      images,
		Description: description,
		CreatedAt:   time.Now().Unix(),
	}
	err := cli.Table(tableName).Create(meal).Error
	return err
}

// GetMealById 根据id查询帖子详情
func (imp *MealInterfaceImp) GetMealById(id, userId int32) (*model.MealModel, error) {
	cli := db.Get()
	var meal model.MealModel // 1. 声明一个具体的结构体变量（不是指针）

	// 2. 将 Where 链式操作连起来，并使用 First 查询单条记录
	err := cli.Table(tableName).
		Select("meals.*, EXISTS(SELECT 1 FROM meal_likes l WHERE l.meal_id = meals.id AND l.user_id = ?) AS is_liked", userId).
		Where("meals.id = ?", id). // 明确指定 meals.id 防止多表字段歧义
		First(&meal).Error         // 使用 First 查询单条

	if err != nil {
		return nil, err // 查询出错（例如没找到，会返回 gorm.ErrRecordNotFound）
	}

	return &meal, nil // 返回指针
}
