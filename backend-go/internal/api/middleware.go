/*
 * VaelorCMS - 现代化内容管理系统
 *
 * Copyright © 2025-2026 Nexsteaduser. All rights reserved.
 *
 * 作者: Tinmc189623
 * 团队: Nexsteaduser
 *
 * 本程序是自由软件: 你可以重新分发和/或修改
 * 它在 GNU Affero 通用公共许可证的条款下发布,
 * 版本 3 或 (根据你的选择) 任何更高版本。
 *
 * 本程序是希望它有用,
 * 但没有任何保证; 甚至没有适销性或
 * 特定用途的默示保证。见
 * GNU Affero 通用公共许可证获取更多细节。
 *
 * 你应该收到 GNU Affero 通用公共许可证的副本
 * 与此程序一起。如果没有, 请见 &lt;https://www.gnu.org/licenses/&gt;.
 */

package api

import (
	"net/http"
	"time"
)

// LoggerMiddleware 日志中间件，记录请求信息
// 参数: next - 下一个 HTTP 处理函数
// 返回: 带有日志功能的 HTTP 处理函数
func LoggerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		next(w, r)
		
		duration := time.Since(start)
	}
}

// RecoveryMiddleware 恢复中间件，处理 panic
// 参数: next - 下一个 HTTP 处理函数
// 返回: 带有 panic 恢复功能的 HTTP 处理函数
func RecoveryMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				WriteError(w, http.StatusInternalServerError, "服务器内部错误")
			}
		}()
		
		next(w, r)
	}
}

// CORSHandler CORS 中间件，处理跨域请求
// 参数: next - 下一个 HTTP 处理函数
// 返回: 带有 CORS 支持的 HTTP 处理函数
func CORSHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next(w, r)
	}
}

// ContentTypeMiddleware 内容类型中间件，设置默认 Content-Type
// 参数: next - 下一个 HTTP 处理函数
// 返回: 带有内容类型设置的 HTTP 处理函数
func ContentTypeMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}
