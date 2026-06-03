package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"wxcloudrun-golang/db/like"

	"github.com/gorilla/mux"
)

// PostLikeRequest 点赞变更数据结构
type PostLikeRequest struct {
	UserId   int32  `json:"user_id"`
	UserName string `json:"user_name"`
}

// LikesResponse 点赞相关接口的统一返回包体
type LikesResponse struct {
	Code     int      `json:"code"`               // 0表示成功，-1表示失败
	ErrorMsg string   `json:"errorMsg,omitempty"` // 错误信息描述
	Data     *LikeDto `json:"data"`               // 点赞变更后返回内容
}

// LikeDto 点赞变更后返回内容
type LikeDto struct {
	IsLiked    bool  `json:"is_liked"`
	LikesCount int32 `json:"likes_count"`
}

// PostLikeHandler 点赞变更（包括点赞和取消）
func PostLikeHandler(w http.ResponseWriter, r *http.Request) {
	res := &LikesResponse{}
	if r.Method == http.MethodPost {
		// 1. 解析路由参数 meals/{id}/likes 中的帖子 id
		vars := mux.Vars(r)
		mealIdStr := vars["id"]
		mealId, _ := strconv.Atoi(mealIdStr)

		// 2. 解析body中的参数
		req, err := parsePostLikeRequest(r)
		if err != nil {
			res.Code = -1
			res.ErrorMsg = fmt.Sprintf("解析点赞失败：%s", err.Error())
		} else {
			// 3. 调用数据访问层查询所有的平铺评论数据
			isLiked, likesCount, err := like.MealLikeImp.PostLike(int32(mealId), req.UserId, req.UserName)
			if err != nil {
				res.Code = -1
				res.ErrorMsg = fmt.Sprintf("点赞失败：%s", err.Error())
			} else {
				res.Data = &LikeDto{
					IsLiked:    isLiked,
					LikesCount: likesCount,
				}
			}
		}
	} else {
		res.Code = -1
		res.ErrorMsg = fmt.Sprintf("请求方法 %s 不支持", r.Method)
	}

	// 4. 返回 JSON 数据
	msg, err := json.Marshal(res)
	if err != nil {
		w.Header().Set("content-type", "application/json")
		w.Write([]byte(`{"code":-1,"errorMsg":"序列化返回结果失败"}`))
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(msg)
}

// parsePostLikeRequest 解析点赞变更请求
func parsePostLikeRequest(r *http.Request) (*PostLikeRequest, error) {
	req := &PostLikeRequest{}

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
	return req, nil
}
