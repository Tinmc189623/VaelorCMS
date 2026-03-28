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
	"errors"
	"net/http"
	"strconv"
	"vaelorcms/pkg/vaelorcore"
)

// ParseJSONRequest 解析 JSON 请求体
// 参数: r - HTTP 请求, target - 目标结构体指针
// 返回: 错误信息，如果解析失败
func ParseJSONRequest(r *http.Request, target interface{}) error {
	if r.Body == nil {
		return errors.New("请求体为空")
	}
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(target); err != nil {
		return err
	}

	return nil
}

// GetQueryParam 获取查询参数
// 参数: r - HTTP 请求, key - 参数名, defaultValue - 默认值
// 返回: 参数值
func GetQueryParam(r *http.Request, key, defaultValue string) string {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetIntQueryParam 获取整数类型的查询参数
// 参数: r - HTTP 请求, key - 参数名, defaultValue - 默认值
// 返回: 参数值
func GetIntQueryParam(r *http.Request, key string, defaultValue int) int {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	}
	result, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return result
}

// GetIntParam 获取路径参数为整数
// 参数: r - HTTP 请求, key - 参数名
// 返回: 参数值，错误信息
func GetIntParam(r *http.Request, key string) (int64, error) {
	value := vaelorcore.GetParam(r, key)
	if value == "" {
		return 0, errors.New("参数不存在")
	}
	return strconv.ParseInt(value, 10, 64)
}
