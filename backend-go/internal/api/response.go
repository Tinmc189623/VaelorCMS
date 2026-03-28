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
	"encoding/json"
	"net/http"
)

// Response 标准 API 响应结构
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// JSONResponse 发送 JSON 格式的响应
// 参数: w - HTTP 响应写入器, statusCode - HTTP 状态码, response - 响应数据
func JSONResponse(w http.ResponseWriter, statusCode int, response Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// WriteSuccess 发送成功响应（简化版本）
// 参数: w - HTTP 响应写入器, data - 响应数据
func WriteSuccess(w http.ResponseWriter, data interface{}) {
	JSONResponse(w, http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

// SuccessResponse 发送成功响应
// 参数: w - HTTP 响应写入器, data - 响应数据, message - 成功消息
func SuccessResponse(w http.ResponseWriter, data interface{}, message string) {
	JSONResponse(w, http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// WriteError 发送错误响应（简化版本）
// 参数: w - HTTP 响应写入器, statusCode - HTTP 状态码, errorMsg - 错误消息
func WriteError(w http.ResponseWriter, statusCode int, errorMsg string) {
	JSONResponse(w, statusCode, Response{
		Success: false,
		Error:   errorMsg,
	})
}

// CreatedResponse 发送创建成功响应
// 参数: w - HTTP 响应写入器, data - 创建的资源数据, message - 成功消息
func CreatedResponse(w http.ResponseWriter, data interface{}, message string) {
	JSONResponse(w, http.StatusCreated, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse 发送错误响应
// 参数: w - HTTP 响应写入器, statusCode - HTTP 状态码, errorMsg - 错误消息
func ErrorResponse(w http.ResponseWriter, statusCode int, errorMsg string) {
	JSONResponse(w, statusCode, Response{
		Success: false,
		Error:   errorMsg,
	})
}
