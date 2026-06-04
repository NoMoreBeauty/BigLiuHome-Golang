package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"wxcloudrun-golang/db/meal"
	"wxcloudrun-golang/db/model"
	"wxcloudrun-golang/logger"

	"github.com/gorilla/mux"
)

// GetMealsRequest 查询三餐帖子参数
type GetMealsRequest struct {
	UserId   int32  `json:"user_id"`
	Page     int32  `json:"page"`
	Size     int32  `json:"size"`
	Date     string `json:"date"`
	MealType string `json:"meal_type"`
}

// PostMealRequest 上传三餐帖子参数
type PostMealRequest struct {
	UserId      int32    `json:"user_id"`
	UserName    string   `json:"username"`
	MealType    string   `json:"meal_type"`
	Images      []string `json:"images"`
	Description string   `json:"description"`
}

// MealsResponse 三餐帖子返回结构(Get/Post)
type MealsResponse struct {
	Code     int      `json:"code"` // -1表示失败、0表示查询成功、1表示插入成果
	ErrorMsg string   `json:"errorMsg,omitempty"`
	Data     MealData `json:"data"`
}

// MealData 帖子具体内容
type MealData struct {
	List  []*model.MealModel `json:"list"`
	Total int32              `json:"total"`
	Page  int32              `json:"page"`
	Size  int32              `json:"size"`
}

// MealsCalendarResponse
type MealsCalendarResponse struct {
	Code     int    `json:"code"` // -1表示失败、0表示查询成功、1表示插入成果
	ErrorMsg string `json:"errorMsg,omitempty"`
	Data     struct {
		Dates []string `json:"dates"`
	} `json:"data"`
}

// GetMealsHandler 查询三餐帖子列表接口
func GetMealsHandler(w http.ResponseWriter, r *http.Request) {
	const mod = "帖子服务"
	start := time.Now()
	res := &MealsResponse{}
	if r.Method == http.MethodGet {
		req, err := parseGetMealsRequest(r)
		if err != nil {
			logger.Warn(mod, "解析帖子列表请求失败", "err", err)
			res.Code = -1
			res.ErrorMsg = err.Error()
		} else {
			logger.Info(mod, "查询帖子列表", "userId", req.UserId, "page", req.Page, "size", req.Size, "date", req.Date, "mealType", req.MealType)
			// 查询帖子
			meals, err := meal.MealImp.GetMeals(req.Page, req.Size, req.UserId, req.Date, req.MealType)
			if err != nil {
				logger.Error(mod, "查询帖子列表失败", "userId", req.UserId, "err", err)
				res.Code = -1
				res.ErrorMsg = err.Error()
			} else {
				logger.Info(mod, "查询帖子列表成功", "userId", req.UserId, "count", len(meals), "耗时", time.Since(start))
				res.Data.List = meals
				res.Data.Total = int32(len(meals))
				res.Data.Page = req.Page
				res.Data.Size = req.Size
			}
		}
	} else if r.Method == http.MethodPost {
		req, err := parsePostMealRequest(r)
		if err != nil {
			logger.Warn(mod, "解析发布帖子请求失败", "err", err)
			res.Code = -1
			res.ErrorMsg = err.Error()
		} else {
			logger.Info(mod, "收到发布帖子请求", "userId", req.UserId, "userName", req.UserName, "mealType", req.MealType, "imageCount", len(req.Images))
			err = meal.MealImp.PostMeals(req.UserId, req.UserName, req.MealType, req.Images, req.Description)
			if err != nil {
				logger.Error(mod, "发布帖子失败", "userId", req.UserId, "err", err)
				res.Code = -1
				res.ErrorMsg = err.Error()
			} else {
				logger.Info(mod, "发布帖子成功", "userId", req.UserId, "mealType", req.MealType, "耗时", time.Since(start))
				res.Code = 1
			}
		}
	} else {
		res.Code = -1
		res.ErrorMsg = fmt.Sprintf("请求方法 %s 不支持", r.Method)
	}

	msg, err := json.Marshal(res)
	if err != nil {
		fmt.Fprint(w, "内部错误")
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(msg)
}

// GetMealDetailHandler 查询三餐帖子详情接口
func GetMealDetailHandler(w http.ResponseWriter, r *http.Request) {
	const mod = "帖子详情服务"
	start := time.Now()

	// 获取动态路由中的id
	vars := mux.Vars(r)
	idStr := vars["id"] // 拿到 "id" 字符串
	id, _ := strconv.Atoi(idStr)

	res := &MealsResponse{}
	if r.Method == http.MethodGet {
		userIdStr := r.URL.Query().Get("user_id")
		userId, _ := strconv.Atoi(userIdStr)

		logger.Info(mod, "查询帖子详情", "mealId", id, "userId", userId)
		meal, err := meal.MealImp.GetMealById(int32(id), int32(userId))
		if err != nil {
			logger.Error(mod, "查询帖子详情失败", "mealId", id, "userId", userId, "err", err)
			res.Code = -1
			res.ErrorMsg = err.Error()
		} else {
			logger.Info(mod, "查询帖子详情成功", "mealId", id, "userId", userId, "耗时", time.Since(start))
			res.Data.List = []*model.MealModel{meal}
		}

	} else {
		res.Code = -1
		res.ErrorMsg = fmt.Sprintf("请求方法 %s 不支持", r.Method)
	}

	msg, err := json.Marshal(res)
	if err != nil {
		fmt.Fprint(w, "内部错误")
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(msg)
}

// GetMealsCalendarHandler 用于日历高亮，查询一个月那几天有帖子
func GetMealsCalendarHandler(w http.ResponseWriter, r *http.Request) {
	const mod = "日历服务"
	start := time.Now()
	res := &MealsCalendarResponse{}
	if r.Method == http.MethodGet {
		year, err := strconv.Atoi(r.URL.Query().Get("year"))
		month, err := strconv.Atoi(r.URL.Query().Get("month"))
		if err != nil || month <= 0 || month >= 13 || year <= 0 || year >= 2100 {
			logger.Warn(mod, "日历查询参数格式错误", "year", year, "month", month)
			res.Code = -1
			res.ErrorMsg = "日期格式错误"
		} else {
			logger.Info(mod, "查询月度日历", "year", year, "month", month)
			dates, err := meal.MealImp.GetMealsCalendar(year, month)
			if err != nil {
				logger.Error(mod, "查询月度日历失败", "year", year, "month", month, "err", err)
				res.Code = -1
				res.ErrorMsg = err.Error()
			} else {
				logger.Info(mod, "查询月度日历成功", "year", year, "month", month, "dayCount", len(dates), "耗时", time.Since(start))
				res.Data.Dates = dates
			}
		}
	} else {
		res.Code = -1
		res.ErrorMsg = fmt.Sprintf("请求方法 %s 不支持", r.Method)
	}

	msg, err := json.Marshal(res)
	if err != nil {
		fmt.Fprint(w, "内部错误")
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(msg)
}

// parseGetMealsRequest 解析查询三餐帖子请求
func parseGetMealsRequest(r *http.Request) (*GetMealsRequest, error) {
	req := &GetMealsRequest{}

	// 解析用户ID
	userId, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		return nil, fmt.Errorf("user_id must be an integer, err=%s", err.Error())
	}
	req.UserId = int32(userId)

	// 解析页码
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		return nil, fmt.Errorf("page must be an integer, err=%s", err.Error())
	}
	req.Page = int32(page)

	// 解析数量
	size, err := strconv.Atoi(r.URL.Query().Get("size"))
	if err != nil {
		return nil, fmt.Errorf("size must be an integer, err=%s", err.Error())
	}
	req.Size = int32(size)

	// 解析日期，这里只检查是不是unix格式，具体转换后面查表的时候做
	date := r.URL.Query().Get("date")
	if date != "" {
		_, err = time.ParseInLocation("2006-01-02", date, time.Local)
		if err != nil {
			return nil, fmt.Errorf("date must be in format YYYY-MM-DD, err=%s", err.Error())
		}
	}
	req.Date = date // 将解析结果赋值给 req，否则日期过滤永远不生效

	// 解析餐点类型
	mealType := r.URL.Query().Get("meal_type")
	if mealType != "" {
		switch mealType {
		case "", "breakfast", "lunch", "dinner", "afternoontea":
		default:
			return nil, fmt.Errorf("meal_type must be one of: breakfast, lunch, dinner")
		}
	}
	req.MealType = mealType // 将解析结果赋值给 req，否则餐类过滤永远不生效
	return req, nil

}

// parseGPostMealRequest 解析插入三餐帖子请求
func parsePostMealRequest(r *http.Request) (*PostMealRequest, error) {
	req := &PostMealRequest{}

	// 1. 从 HTTP Body 中解码 JSON 到结构体中
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		return nil, fmt.Errorf("解析请求体 JSON 失败, err=%s", err.Error())
	}
	// 2. 字段基本校验
	if req.UserId <= 0 {
		return nil, fmt.Errorf("user_id 必须大于 0")
	}
	if req.UserName == "" {
		return nil, fmt.Errorf("username 不能为空")
	}
	// 3. 校验餐点类型
	switch req.MealType {
	case "breakfast", "lunch", "dinner", "afternoontea":
	default:
		return nil, fmt.Errorf("meal_type 必须是 breakfast, lunch, dinner, afternoontea 之一")
	}
	return req, nil
}
