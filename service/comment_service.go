package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"wxcloudrun-golang/db/comment"
	"wxcloudrun-golang/db/model"
	"wxcloudrun-golang/logger"

	"github.com/gorilla/mux"
)

// PostCommentRequest 插入评论数据结构
type PostCommentRequest struct {
	UserId      int32  `json:"user_id"`
	UserName    string `json:"user_name"`
	ParentId    int32  `json:"parent_id"` // 一级评论传0
	ReplyToId   int32  `json:"reply_to_id"`
	ReplyToName string `json:"reply_to_name"`
	Content     string `json:"content"`
}

// CommentDto 用于接口返回的嵌套评论结构（二级盖楼）
type CommentDto struct {
	*model.MealCommentModel                           // 组合原始的评论模型
	Replies                 []*model.MealCommentModel `json:"replies"` // 存放该一级评论下的所有二级回复
}

// CommentsResponse 评论相关接口的统一返回包体
type CommentsResponse struct {
	Code     int           `json:"code"`               // 0表示成功，-1表示失败
	ErrorMsg string        `json:"errorMsg,omitempty"` // 错误信息描述
	Data     []*CommentDto `json:"data"`               // 嵌套评论列表
}

// GetMealCommentsHandler 查询某篇三餐帖子的所有评论（支持两层盖楼）
func GetMealCommentsHandler(w http.ResponseWriter, r *http.Request) {
	/**
	 * 业务流程：
	 *  1. 从路由参数解析出帖子 id（meal_id）
	 *  2. 调用 db 层获取该帖子下的全部评论记录（按时间正序排列）
	 *  3. 在内存中将平铺的记录组装成「两层盖楼」的嵌套树结构：
	 *     - 过滤出 parent_id 为 0 的记录作为「一级评论」（根节点）
	 *     - 将 parent_id > 0 的二级回复按 parent_id 归入对应的一级评论 replies 数组中
	 *  4. 将组装好的嵌套数组返回给前端
	 */
	const mod = "评论服务"
	start := time.Now()
	res := &CommentsResponse{
		Code: 0,
		Data: make([]*CommentDto, 0),
	}

	if r.Method == http.MethodGet {
		// 1. 解析路由参数 meals/{id}/comments 中的帖子 id
		vars := mux.Vars(r)
		mealIdStr := vars["id"]
		mealId, _ := strconv.Atoi(mealIdStr)

		logger.Info(mod, "查询帖子评论列表", "mealId", mealId)
		// 2. 调用数据访问层查询所有的平铺评论数据
		flatComments, err := comment.MealCommentImp.GetMealCommentsHandler(int32(mealId))
		if err != nil {
			logger.Error(mod, "查询评论列表失败", "mealId", mealId, "err", err)
			res.Code = -1
			res.ErrorMsg = fmt.Sprintf("获取评论失败：%s", err.Error())
		} else {
			// 3. 在内存中执行两层嵌套组装
			res.Data = assembleComments(flatComments)
			logger.Info(mod, "查询评论列表成功", "mealId", mealId, "flatCount", len(flatComments), "rootCount", len(res.Data), "耗时", time.Since(start))
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

// PostCommentHandler 插入评论
func PostCommentHandler(w http.ResponseWriter, r *http.Request) {
	const mod = "评论服务"
	start := time.Now()
	res := &CommentsResponse{
		Code: 0,
		Data: make([]*CommentDto, 0),
	}
	if r.Method == http.MethodPost {
		// 1. 解析路由参数 meals/{id}/comments 中的帖子 id
		vars := mux.Vars(r)
		mealIdStr := vars["id"]
		mealId, _ := strconv.Atoi(mealIdStr)

		// 2. 解析请求体中的参数
		req, err := parsePostCommentRequest(r)
		if err != nil {
			logger.Warn(mod, "解析发布评论请求失败", "mealId", mealId, "err", err)
			res.Code = -1
			res.ErrorMsg = fmt.Sprintf("解析评论失败：%s", err.Error())
			return
		}

		logger.Info(mod, "收到发布评论请求", "mealId", mealId, "userId", req.UserId, "userName", req.UserName, "parentId", req.ParentId, "contentLen", len(req.Content))
		// 3. 调用数据访问层插入评论数据
		newComment, err := comment.MealCommentImp.PostCommentHandler(int32(mealId), req.UserId, req.UserName, req.ParentId, req.ReplyToId, req.ReplyToName, req.Content)
		if err != nil {
			logger.Error(mod, "发布评论失败", "mealId", mealId, "userId", req.UserId, "err", err)
			res.Code = -1
			res.ErrorMsg = fmt.Sprintf("获取评论失败：%s", err.Error())
		} else {
			logger.Info(mod, "发布评论成功", "mealId", mealId, "userId", req.UserId, "commentId", newComment.Id, "耗时", time.Since(start))
			res.Data = append(res.Data, &CommentDto{MealCommentModel: newComment})
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

// assembleComments 嵌套组装算法
func assembleComments(flatComments []*model.MealCommentModel) []*CommentDto {
	/**
	 * 将一维的平铺评论数组组装成含有 replies 列表的二级评论树结构
	 * @param flatComments 数据库查出来的平铺评论数组（按时间正序）
	 * @returns 组装后的一级评论列表
	 */
	var roots []*CommentDto
	// 建立 map 用来快速检索一级评论节点，key 为评论 ID
	rootMap := make(map[int32]*CommentDto)

	// 第一次遍历：筛选并创建所有一级评论（ParentId 等于 0 的评论）
	for _, c := range flatComments {
		if c.ParentId == 0 {
			dto := &CommentDto{
				MealCommentModel: c,
				// 显式初始化为非 nil 切片，确保序列化出来的是 [] 而不是 null
				Replies: make([]*model.MealCommentModel, 0),
			}
			roots = append(roots, dto)
			rootMap[c.Id] = dto
		}
	}

	// 第二次遍历：将所有二级回复（ParentId 大于 0）分配到对应的一级评论下
	for _, c := range flatComments {
		if c.ParentId > 0 {
			// 根据 ParentId 检索对应的一级评论根节点
			if parent, exists := rootMap[c.ParentId]; exists {
				parent.Replies = append(parent.Replies, c)
			}
		}
	}

	return roots
}

// parsePostCommentRequest 解析插入评论请求
func parsePostCommentRequest(r *http.Request) (*PostCommentRequest, error) {
	req := &PostCommentRequest{}

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
