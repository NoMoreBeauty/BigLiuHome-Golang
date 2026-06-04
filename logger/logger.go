// Package logger 提供统一的中文结构化日志能力。
// 基于标准库 log 封装，不引入任何第三方依赖，输出格式统一、易于检索。
//
// 日志格式：[级别] 时间 - 模块 - 消息 [key=value ...]
// 示例：
//   [INFO] 2026/06/04 12:00:00 - 点赞服务 - 收到点赞请求 mealId=4 userId=2
//   [WARN] 2026/06/04 12:00:00 - 点赞DB - 未找到点赞记录（首次点赞）mealId=4 userId=2
//   [ERROR] 2026/06/04 12:00:00 - 点赞DB - 数据库操作失败 err=...
package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// 全局 logger，使用标准输出，带日期时间前缀
var std = log.New(os.Stdout, "", log.LstdFlags)

// Info 打印普通信息日志（正常业务流程）
// module: 模块名，如 "点赞服务"、"点赞DB"
// msg: 简短描述
// kvPairs: 可选的 key=value 对，如 "mealId", 4, "userId", 2
func Info(module, msg string, kvPairs ...interface{}) {
	std.Printf("[INFO] %s - %s%s", module, msg, formatKV(kvPairs...))
}

// Warn 打印警告日志（非正常但可接受的情况，如记录不存在）
func Warn(module, msg string, kvPairs ...interface{}) {
	std.Printf("[WARN] %s - %s%s", module, msg, formatKV(kvPairs...))
}

// Error 打印错误日志（需要关注的异常）
func Error(module, msg string, kvPairs ...interface{}) {
	std.Printf("[ERROR] %s - %s%s", module, msg, formatKV(kvPairs...))
}

// formatKV 将 key-value 对格式化为 " key=value key=value" 形式的字符串
// 传入参数应为偶数个，奇数位为 key（string），偶数位为 value（任意类型）
func formatKV(kvPairs ...interface{}) string {
	if len(kvPairs) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString("  ")
	for i := 0; i+1 < len(kvPairs); i += 2 {
		if i > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(fmt.Sprintf("%v=%v", kvPairs[i], kvPairs[i+1]))
	}
	return sb.String()
}
