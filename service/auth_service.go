package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"wxcloudrun-golang/db/auth"

	"gorm.io/gorm"
)

// AuthResponse 登录返回结构
type AuthResponse struct {
	Code     int         `json:"code"`
	ErrorMsg string      `json:"errorMsg,omitempty"`
	Data     interface{} `json:"data"`
}

// AuthHandler 用户登录接口
func AuthHandler(w http.ResponseWriter, r *http.Request) {
	res := &AuthResponse{}

	if r.Method == http.MethodGet {
		userKey := r.URL.Query().Get("user_key")
		if userKey == "" {
			res.Code = -1
			res.ErrorMsg = "缺少 user_key 参数"
		} else {
			user, err := auth.UserImp.Login(userKey)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					res.ErrorMsg = "用户不存在"
				} else {
					res.ErrorMsg = "系统繁忙，请稍后再试" // 屏蔽真实的数据库报错，防止暴露安全信息
				}
			} else {
				res.Data = user
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
