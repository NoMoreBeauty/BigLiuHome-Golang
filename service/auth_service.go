package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"wxcloudrun-golang/db/auth"
	"wxcloudrun-golang/logger"

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
	const mod = "登录服务"
	start := time.Now()
	res := &AuthResponse{}

	if r.Method == http.MethodGet {
		userKey := r.URL.Query().Get("user_key")
		if userKey == "" {
			logger.Warn(mod, "登录请求缺少 user_key 参数")
			res.Code = -1
			res.ErrorMsg = "缺少 user_key 参数"
		} else {
			logger.Info(mod, "收到登录请求", "userKey", userKey)
			user, err := auth.UserImp.Login(userKey)
			if err != nil {
				res.Code = -1
				if errors.Is(err, gorm.ErrRecordNotFound) {
					logger.Warn(mod, "用户不存在", "userKey", userKey)
					res.ErrorMsg = "用户不存在"
				} else {
					logger.Error(mod, "登录查询数据库失败", "userKey", userKey, "err", err)
					res.ErrorMsg = "系统繁忙，请稍后再试" // 屏蔽真实的数据库报错，防止暴露安全信息
				}
			} else {
				logger.Info(mod, "登录成功", "userKey", userKey, "userId", user.Id, "耗时", time.Since(start))
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

// GetMembersHandler 返回家庭成员列表
func GetMembersHandler(w http.ResponseWriter, r *http.Request) {
	const mod = "拉取用户成员信息"
	start := time.Now()
	res := &AuthResponse{}

	if r.Method == http.MethodGet {

		logger.Info(mod, "收到拉取用户成员信息请求")
		userList, err := auth.UserImp.GetMembers()
		if err != nil {
			res.Code = -1
			logger.Error(mod, "拉取用户成员信息失败", "err", err)
			res.ErrorMsg = "系统繁忙，请稍后再试" // 屏蔽真实的数据库报错，防止暴露安全信息
		} else {
			logger.Info(mod, "拉取用户成员信息成功", "耗时", time.Since(start))
			res.Data = userList
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
