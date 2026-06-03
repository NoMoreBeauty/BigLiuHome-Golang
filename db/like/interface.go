package like

// MealLikeInterface 三餐帖子点赞数据模型接口
type MealLikeInterface interface {
	PostLike(mealId, userId int32, userName string) (bool, int32, error)
}

// MealLikeInterfaceImp MealLikeInterface的实现对象
type MealLikeInterfaceImp struct{}

// 实现实例
var MealLikeImp MealLikeInterface = &MealLikeInterfaceImp{}
