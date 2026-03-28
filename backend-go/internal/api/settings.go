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
	"database/sql"
	"net/http"
	"vaelorcms/internal/services"
	"vaelorcms/pkg/vaelorcore"
)

var settingService = services.NewSettingService()

// RegisterSettingRoutes 注册设置相关路由
// 参数: router - Vaelor Core 路由器
func RegisterSettingRoutes(router *vaelorcore.Router) {
	router.GET("/api/settings", ListSettingsHandler)
	router.GET("/api/settings/:id", GetSettingHandler)
	router.GET("/api/settings/key/:key", GetSettingByKeyHandler)
	router.POST("/api/settings", CreateSettingHandler)
	router.PUT("/api/settings/:id", UpdateSettingHandler)
	router.PUT("/api/settings/bulk", BulkSetSettingsHandler)
	router.DELETE("/api/settings/:id", DeleteSettingHandler)
}

// ListSettingsHandler 处理获取设置列表请求
// 参数: w - HTTP 响应写入器, r - HTTP 请求
func ListSettingsHandler(w http.ResponseWriter, r *http.Request) {
	groupStr := GetQueryParam(r, "group", "")

	var group *string
	if groupStr != "" {
		group = &groupStr
	}

	settings, err := settingService.GetSettings(group)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "获取设置列表失败")
		return
	}

	WriteSuccess(w, map[string]interface{}{
		"settings": settings,
	})
}

// GetSettingHandler 处理根据ID获取设置请求
// 参数: w - HTTP 响应写入器, r - HTTP 请求
func GetSettingHandler(w http.ResponseWriter, r *http.Request) {
	settingID, err := GetIntParam(r, "id")
	if err != nil {
		WriteError(w, http.StatusBadRequest, "无效的设置ID")
		return
	}

	setting, err := settingService.GetSettingByID(settingID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "获取设置失败")
		return
	}

	if setting == nil {
		WriteError(w, http.StatusNotFound, "设置不存在")
		return
	}

	WriteSuccess(w, setting)
}

// GetSettingByKeyHandler 处理根据Key获取设置请求
// 参数: w - HTTP 响应写入器, r - HTTP 请求
func GetSettingByKeyHandler(w http.ResponseWriter, r *http.Request) {
	key := vaelorcore.GetParam(r, "key")
	if key == "" {
		WriteError(w, http.StatusBadRequest, "无效的设置Key")
		return
	}

	setting, err := settingService.GetSettingByKey(key)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "获取设置失败")
		return
	}

	if setting == nil {
		WriteError(w, http.StatusNotFound, "设置不存在")
		return
	}

	WriteSuccess(w, setting)
}

// CreateSettingHandler 处理创建设置请求
// 参数: w - HTTP 响应写入器, r - HTTP 请求
func CreateSettingHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Key         string  `json:"key"`
		Value       *string `json:"value"`
		Description *string `json:"description"`
		Group       *string `json:"group"`
	}

	if err := ParseJSONRequest(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "无效的请求格式")
		return
	}

	if req.Key == "" {
		WriteError(w, http.StatusBadRequest, "设置Key不能为空")
		return
	}

	var value sql.NullString
	if req.Value != nil {
		value = sql.NullString{String: *req.Value, Valid: true}
	}

	var description sql.NullString
	if req.Description != nil {
		description = sql.NullString{String: *req.Description, Valid: true}
	}

	var group sql.NullString
	if req.Group != nil {
		group = sql.NullString{String: *req.Group, Valid: true}
	}

	settingData := services.SettingCreate{
		Key:         req.Key,
		Value:       value,
		Description: description,
		Group:       group,
	}

	setting, err := settingService.CreateSetting(settingData)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteSuccess(w, setting)
}

// UpdateSettingHandler 处理更新设置请求
// 参数: w - HTTP 响应写入器, r - HTTP 请求
func UpdateSettingHandler(w http.ResponseWriter, r *http.Request) {
	settingID, err := GetIntParam(r, "id")
	if err != nil {
		WriteError(w, http.StatusBadRequest, "无效的设置ID")
		return
	}

	var req struct {
		Key         *string `json:"key"`
		Value       *string `json:"value"`
		Description *string `json:"description"`
		Group       *string `json:"group"`
	}

	if err := ParseJSONRequest(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "无效的请求格式")
		return
	}

	var updateData services.SettingUpdate

	if req.Key != nil {
		updateData.Key = req.Key
	}
	if req.Value != nil {
		value := sql.NullString{String: *req.Value, Valid: true}
		updateData.Value = &value
	}
	if req.Description != nil {
		description := sql.NullString{String: *req.Description, Valid: true}
		updateData.Description = &description
	}
	if req.Group != nil {
		group := sql.NullString{String: *req.Group, Valid: true}
		updateData.Group = &group
	}

	setting, err := settingService.UpdateSetting(settingID, updateData)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if setting == nil {
		WriteError(w, http.StatusNotFound, "设置不存在")
		return
	}

	WriteSuccess(w, setting)
}

// BulkSetSettingsHandler 处理批量更新设置请求
// 参数: w - HTTP 响应写入器, r - HTTP 请求
func BulkSetSettingsHandler(w http.ResponseWriter, r *http.Request) {
	var settings map[string]string

	if err := ParseJSONRequest(r, &settings); err != nil {
		WriteError(w, http.StatusBadRequest, "无效的请求格式")
		return
	}

	success, err := settingService.BulkSetSettings(settings)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteSuccess(w, map[string]bool{"success": success})
}

// DeleteSettingHandler 处理删除设置请求
// 参数: w - HTTP 响应写入器, r - HTTP 请求
func DeleteSettingHandler(w http.ResponseWriter, r *http.Request) {
	settingID, err := GetIntParam(r, "id")
	if err != nil {
		WriteError(w, http.StatusBadRequest, "无效的设置ID")
		return
	}

	success, err := settingService.DeleteSetting(settingID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "删除设置失败")
		return
	}

	if !success {
		WriteError(w, http.StatusNotFound, "设置不存在")
		return
	}

	WriteSuccess(w, map[string]bool{"success": true})
}
