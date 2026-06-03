# BigLiuHome Golang API 接口文档

本文档为 **BigLiuHome** 后端的 API 接口定义，供前端开发人员参考对接。

*   **服务默认地址**: `http://localhost:80` (或部署后的实际域名)
*   **统一响应格式**: 所有 API（除主页静态 HTML 路由外）均返回 `application/json` 格式的数据，其中包含 `code`、`errorMsg`（可选）以及 `data` 字段。
    *   `code = 0` 表示请求处理成功（部分创建接口可能会返回其他正数，如 `1`）。
    *   `code = -1` 表示请求处理失败，可通过 `errorMsg` 获取具体错误提示。

---

## 目录

1. [静态页面路由](#1-静态页面路由-)
2. [获取计数器](#2-获取计数器-get-apicount)
3. [更新计数器](#3-更新计数器-post-apicount)
4. [用户登录](#4-用户登录-get-apilogin)
5. [查询三餐帖子列表](#5-查询三餐帖子列表-get-apimeals)
6. [发布三餐帖子](#6-发布三餐帖子-post-apimeals)
7. [查询三餐帖子详情](#7-查询三餐帖子详情-get-apimealsid)
8. [获取帖子评论列表](#8-获取帖子评论列表-get-apimealsidcomments)
9. [发表评论 / 回复](#9-发表评论--回复-post-apimealsidcomments)
10. [点赞 / 取消点赞](#10-点赞--取消点赞-post-apimealsidlike)

---

## 1. 静态页面路由 `/`

*   **请求路径**: `/`
*   **请求方式**: `GET`
*   **请求参数**: 无
*   **返回内容**: 返回 `index.html` 的 HTML 静态网页内容。

---

## 2. 获取计数器 `GET /api/count`

*   **请求路径**: `/api/count`
*   **请求方式**: `GET`
*   **请求参数**: 无
*   **返回内容示例**:
    *   **成功 (`code: 0`)**:
        ```json
        {
          "code": 0,
          "data": 12
        }
        ```
    *   **失败 (`code: -1`)**:
        ```json
        {
          "code": -1,
          "errorMsg": "数据库连接超时"
        }
        ```

---

## 3. 更新计数器 `POST /api/count`

*   **请求路径**: `/api/count`
*   **请求方式**: `POST`
*   **请求体格式**: `application/json`
*   **请求体参数**:
    | 字段名 | 类型 | 必填 | 说明 |
    | :--- | :--- | :--- | :--- |
    | `action` | string | 是 | 执行的动作。可选值：`"inc"` (自增)、`"clear"` (清零) |
*   **请求体示例**:
    ```json
    {
      "action": "inc"
    }
    ```
*   **返回内容示例**:
    *   **自增成功 (`code: 0`)**:
        ```json
        {
          "code": 0,
          "data": 13
        }
        ```
    *   **清零成功 (`code: 0`)**:
        ```json
        {
          "code": 0,
          "data": 0
        }
        ```
    *   **失败 (`code: -1`)**:
        ```json
        {
          "code": -1,
          "errorMsg": "参数 action : invalid_action 错误"
        }
        ```

---

## 4. 用户登录 `GET /api/login`

*   **请求路径**: `/api/login`
*   **请求方式**: `GET`
*   **请求参数 (Query)**:
    | 参数名 | 类型 | 必填 | 说明 |
    | :--- | :--- | :--- | :--- |
    | `user_key` | string | 是 | 用户的手机号或唯一标识 |
*   **请求示例**:
    `GET http://localhost/api/login?user_key=13800138000`
*   **返回内容示例**:
    *   **登录成功 (`code: 0`)**:
        ```json
        {
          "code": 0,
          "data": {
            "id": 1,
            "user_key": "13800138000",
            "user_name": "大刘",
            "role": "admin",
            "status": 1,
            "created_at": 1685785200,
            "updated_at": 1685785200
          }
        }
        ```
    *   **登录失败 - 用户不存在 (`code: -1`)**:
        ```json
        {
          "code": -1,
          "errorMsg": "用户不存在"
        }
        ```
    *   **登录失败 - 缺少参数 (`code: -1`)**:
        ```json
        {
          "code": -1,
          "errorMsg": "缺少 user_key 参数"
        }
        ```

---

## 5. 查询三餐帖子列表 `GET /api/meals`

*   **请求路径**: `/api/meals`
*   **请求方式**: `GET`
*   **请求参数 (Query)**:
    | 参数名 | 类型 | 必填 | 说明 |
    | :--- | :--- | :--- | :--- |
    | `user_id` | integer | 是 | 当前登录的用户 ID（用于判断用户是否点赞过帖子） |
    | `page` | integer | 是 | 页码，从 1 开始 |
    | `size` | integer | 是 | 每页限制条数 |
    | `date` | string | 否 | 按日期过滤，格式为 `YYYY-MM-DD`（如 `2026-06-03`） |
    | `meal_type` | string | 否 | 按餐点类型过滤，可选值：`breakfast` (早餐), `lunch` (午餐), `dinner` (晚餐), `afternoontea` (下午茶) |
*   **请求示例**:
    `GET http://localhost/api/meals?user_id=1&page=1&size=10&meal_type=breakfast`
*   **返回内容示例**:
    *   **成功 (`code: 0`)**:
        ```json
        {
          "code": 0,
          "data": {
            "list": [
              {
                "id": 1,
                "user_id": 1,
                "user_name": "大刘",
                "meal_type": "breakfast",
                "images": [
                  "https://example.com/1.jpg",
                  "https://example.com/2.jpg"
                ],
                "description": "今天的美味早餐！",
                "likes_count": 2,
                "comments_count": 0,
                "created_at": 1780482000,
                "is_liked": true
              }
            ],
            "total": 1,
            "page": 1,
            "size": 10
          }
        }
        ```
    *   **失败 (`code: -1`)**:
        ```json
        {
          "code": -1,
          "errorMsg": "user_id must be an integer, err=..."
        }
        ```

> [!NOTE]
> **关于 `images` 字段的说明**:
> 数据库中 `images` 字段已重构为 `json.RawMessage`，接口直接返回标准的图片 URL 数组（`array of string`），前端无需进行 Base64 解码，可直接循环渲染。

---

## 6. 发布三餐帖子 `POST /api/meals`

*   **请求路径**: `/api/meals`
*   **请求方式**: `POST`
*   **请求体格式**: `application/json`
*   **请求体参数**:
    | 字段名 | 类型 | 必填 | 说明 |
    | :--- | :--- | :--- | :--- |
    | `user_id` | integer | 是 | 发布用户的 ID |
    | `username` | string | 是 | 发布用户的用户名 |
    | `meal_type` | string | 是 | 餐点类型。可选值：`breakfast`, `lunch`, `dinner`, `afternoontea` |
    | `images` | array of string | 是 | 图片 URL 数组，空数组传 `[]` |
    | `description` | string | 否 | 帖子的文字描述描述 |
*   **请求体示例**:
    ```json
    {
      "user_id": 1,
      "username": "大刘",
      "meal_type": "breakfast",
      "images": [
        "https://example.com/meal1.jpg",
        "https://example.com/meal2.jpg"
      ],
      "description": "煎蛋和全麦面包，美味健康"
    }
    ```
*   **返回内容示例**:
    *   **发布成功 (`code: 1`)**:
        ```json
        {
          "code": 1
        }
        ```
    *   **发布失败 (`code: -1`)**:
        ```json
        {
          "code": -1,
          "errorMsg": "meal_type 必须是 breakfast, lunch, dinner, afternoontea 之一"
        }
        ```

---

## 7. 查询三餐帖子详情 `GET /api/meals/{id}`

*   **请求路径**: `/api/meals/{id}` (其中 `{id}` 为帖子 ID)
*   **请求方式**: `GET`
*   **请求参数 (Query)**:
    | 参数名 | 类型 | 必填 | 说明 |
    | :--- | :--- | :--- | :--- |
    | `user_id` | integer | 否 | 当前登录的用户 ID（用于计算当前用户是否点赞过该帖子） |
*   **请求示例**:
    `GET http://localhost/api/meals/1?user_id=1`
*   **返回内容示例**:
    *   **查询成功 (`code: 0`)**:
        ```json
        {
          "code": 0,
          "data": {
            "list": [
              {
                "id": 1,
                "user_id": 1,
                "user_name": "大刘",
                "meal_type": "breakfast",
                "images": [
                  "https://example.com/1.jpg",
                  "https://example.com/2.jpg"
                ],
                "description": "今天的美味早餐！",
                "likes_count": 2,
                "comments_count": 0,
                "created_at": 1780482000,
                "is_liked": true
              }
            ],
            "total": 0,
            "page": 0,
            "size": 0
          }
        }
        ```
    *   **查询失败 - 帖子不存在 (`code: -1`)**:
        ```json
        {
          "code": -1,
          "errorMsg": "record not found"
        }
        ```

---

## 8. 获取帖子评论列表 `GET /api/meals/{id}/comments`

*   **请求路径**: `/api/meals/{id}/comments` (其中 `{id}` 为帖子 ID)
*   **请求方式**: `GET`
*   **请求参数**: 无
*   **返回内容**: 接口在内存中会自动将评论按一级评论和二级回复组装为**两层树状盖楼结构**。
*   **返回内容示例**:
    *   **成功 (`code: 0`)**:
        ```json
        {
          "code": 0,
          "data": [
            {
              "id": 10,
              "meal_id": 1,
              "parent_id": 0,
              "user_id": 2,
              "user_name": "小明",
              "reply_to_id": 0,
              "reply_to_name": "",
              "content": "哇，看起来真不错！",
              "created_at": 1780482100,
              "replies": [
                {
                  "id": 11,
                  "meal_id": 1,
                  "parent_id": 10,
                  "user_id": 1,
                  "user_name": "大刘",
                  "reply_to_id": 2,
                  "reply_to_name": "小明",
                  "content": "多谢，哈哈！",
                  "created_at": 1780482200
                }
              ]
            }
          ]
        }
        ```

---

## 9. 发表评论 / 回复 `POST /api/meals/{id}/comments`

*   **请求路径**: `/api/meals/{id}/comments` (其中 `{id}` 为帖子 ID)
*   **请求方式**: `POST`
*   **请求体格式**: `application/json`
*   **请求体参数**:
    | 字段名 | 类型 | 必填 | 说明 |
    | :--- | :--- | :--- | :--- |
    | `user_id` | integer | 是 | 评论发布者的用户 ID |
    | `user_name` | string | 是 | 评论发布者的用户名 |
    | `parent_id` | integer | 是 | 所属一级评论 ID。如果是对帖子的直接评论，传 `0`；如果是回复别人的评论，传该回复所属的**最上层一级评论的 ID** |
    | `reply_to_id` | integer | 否 | 被回复用户的 ID。直接评论帖子传 `0` |
    | `reply_to_name` | string | 否 | 被回复用户的用户名。直接评论帖子传空字符串 `""` |
    | `content` | string | 是 | 评论内容 |
*   **请求体示例**:
    *   **对帖子直接评论 (一级评论)**:
        ```json
        {
          "user_id": 2,
          "user_name": "小明",
          "parent_id": 0,
          "reply_to_id": 0,
          "reply_to_name": "",
          "content": "哇，看起来真不错！"
        }
        ```
    *   **回复某个一级评论 (二级回复)**:
        ```json
        {
          "user_id": 1,
          "user_name": "大刘",
          "parent_id": 10,
          "reply_to_id": 2,
          "reply_to_name": "小明",
          "content": "多谢，哈哈！"
        }
        ```
*   **返回内容示例**:
    *   **发表成功 (`code: 0`)**: 返回包含当前新增评论内容的单元素数组。
        ```json
        {
          "code": 0,
          "data": [
            {
              "id": 12,
              "meal_id": 1,
              "parent_id": 0,
              "user_id": 2,
              "user_name": "小明",
              "reply_to_id": 0,
              "reply_to_name": "",
              "content": "哇，看起来真不错！",
              "created_at": 1780482300,
              "replies": null
            }
          ]
        }
        ```

> [!NOTE]
> 成功提交后，返回的一级评论结构中 `"replies"` 字段的值为 `null`。

---

## 10. 点赞 / 取消点赞 `POST /api/meals/{id}/like`

*   **请求路径**: `/api/meals/{id}/like` (其中 `{id}` 为帖子 ID)
*   **请求方式**: `POST`
*   **请求体格式**: `application/json`
*   **请求体参数**:
    | 字段名 | 类型 | 必填 | 说明 |
    | :--- | :--- | :--- | :--- |
    | `user_id` | integer | 是 | 当前操作用户的 ID |
    | `user_name` | string | 是 | 当前操作用户的用户名 |
*   **请求体示例**:
    ```json
    {
      "user_id": 1,
      "user_name": "大刘"
    }
    ```
*   **返回内容示例**:
    *   **成功 (`code: 0`)**:
        ```json
        {
          "code": 0,
          "data": {
            "is_liked": true,
            "likes_count": 3
          }
        }
        ```
