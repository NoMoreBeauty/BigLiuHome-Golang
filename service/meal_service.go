package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"wxcloudrun-golang/db/meal"
	"wxcloudrun-golang/db/model"

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

// GetMealsHandler 查询三餐帖子列表接口
func GetMealsHandler(w http.ResponseWriter, r *http.Request) {
	res := &MealsResponse{}
	if r.Method == http.MethodGet {
		req, err := parseGetMealsRequest(r)
		if err != nil {
			res.Code = -1
			res.ErrorMsg = err.Error()
		} else {
			// 查询帖子
			meals, err := meal.MealImp.GetMeals(req.Page, req.Size, req.UserId, req.Date, req.MealType)
			if err != nil {
				res.Code = -1
				res.ErrorMsg = err.Error()
			} else {
				res.Data.List = meals
				res.Data.Total = int32(len(meals))
				res.Data.Page = req.Page
				res.Data.Size = req.Size
			}
		}
	} else if r.Method == http.MethodPost {
		req, err := parsePostMealRequest(r)
		if err != nil {
			res.Code = -1
			res.ErrorMsg = err.Error()
		} else {
			// 要做图片类型的转化来适配mysql的json字段
			imagesByte, err := json.Marshal(req.Images)
			if err != nil {
				res.Code = -1
				res.ErrorMsg = err.Error()
			} else {
				err = meal.MealImp.PostMeals(req.UserId, req.UserName, req.MealType, imagesByte, req.Description)
				if err != nil {
					res.Code = -1
					res.ErrorMsg = err.Error()
				} else {
					res.Code = 1
				}
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

	// 获取动态路由中的id
	vars := mux.Vars(r)
	idStr := vars["id"] // 拿到 "id" 字符串
	id, _ := strconv.Atoi(idStr)

	res := &MealsResponse{}
	if r.Method == http.MethodGet {
		userIdStr := r.URL.Query().Get("user_id")
		userId, _ := strconv.Atoi(userIdStr)

		meal, err := meal.MealImp.GetMealById(int32(id), int32(userId))
		if err != nil {
			res.Code = -1
			res.ErrorMsg = err.Error()
		} else {
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

	// 解析餐点类型
	mealType := r.URL.Query().Get("meal_type")
	if mealType != "" {
		switch mealType {
		case "", "breakfast", "lunch", "dinner", "afternoontea":
		default:
			return nil, fmt.Errorf("meal_type must be one of: breakfast, lunch, dinner")
		}
	}
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
