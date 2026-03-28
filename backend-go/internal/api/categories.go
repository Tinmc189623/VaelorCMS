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

var categoryService = services.NewCategoryService()

// RegisterCategoryRoutes 注册分类相关路由
// 参数: router - Vaelor Core 路由器
func RegisterCategoryRoutes(router *vaelorcore.Router) {
	router.GET("/api/categories", ListCategoriesHandler)
	router.GET("/api/categories/:id", GetCategoryHandler)
	router.GET("/api/categories/slug/:slug", GetCategoryBySlugHandler)
	router.POST("/api/categories", CreateCategoryHandler)
	router.PUT("/api/categories/:id", UpdateCategoryHandler)
	router.DELETE("/api/categories/:id", DeleteCategoryHandler)
}

// ListCategoriesHandler 处理获取分类列表请求
// 参数: w - HTTP 响应写入器, r - HTTP 请求
func ListCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	skip := GetIntQueryParam(r, "skip", 0)
	limit := GetIntQueryParam(r, "limit", 20)
	
	var parentID *int64
	parentIDStr := GetQueryParam(r, "parent_id", "")
	if parentIDStr != "" {
		var id int64
		tempID, err := GetIntParam(r, "parent_id")
		if err == nil {
			id = tempID
			parentID = &amp;id
		}
	}

	categories, err := categoryService.GetCategories(skip, limit, parentID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "获取分类列表失败")
		return
	}

	total, err := categoryService.CountCategories(parentID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "统计分类数量失败")
		return
	}

	WriteSuccess(w, map[string]interface{}{
		"categories": categories,
		"total":      total,
		"skip":       skip,
		"limit":      limit,
	})
}

// GetCategoryHandler 处理根据ID获取分类请求
// 参数: w - HTTP 响应写入器, r - HTTP 请求
func GetCategoryHandler(w http.ResponseWriter, r *http.Request) {
	categoryID, err := GetIntParam(r, "id")
	if err != nil {
		WriteError(w, http.StatusBadRequest, "无效的分类ID")
		return
	}

	category, err := categoryService.GetCategoryByID(categoryID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "获取分类失败")
		return
	}

	if category == nil {
		WriteError(w, http.StatusNotFound, "分类不存在")
		return
	}

	WriteSuccess(w, category)
}

// GetCategoryBySlugHandler 处理根据Slug获取分类请求
// 参数: w - HTTP 响应写入器, r - HTTP 请求
func GetCategoryBySlugHandler(w http.ResponseWriter, r *http.Request) {
	slug := vaelorcore.GetParam(r, "slug")
	if slug == "" {
		WriteError(w, http.StatusBadRequest, "无效的分类Slug")
		return
	}

	category, err := categoryService.GetCategoryBySlug(slug)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "获取分类失败")
		return
	}

	if category == nil {
		WriteError(w, http.StatusNotFound, "分类不存在")
		return
	}

	WriteSuccess(w, category)
}

// CreateCategoryHandler 处理创建分类请求
// 参数: w - HTTP 响应写入器, r - HTTP 请求
func CreateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string  `json:"name"`
		Slug        string  `json:"slug"`
		Description *string `json:"description"`
		ParentID    *int64  `json:"parent_id"`
	}

	if err := ParseJSONRequest(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "无效的请求格式")
		return
	}

	if req.Name == "" || req.Slug == "" {
		WriteError(w, http.StatusBadRequest, "分类名称和Slug不能为空")
		return
	}

	category, err := categoryService.CreateCategory(req.Name, req.Slug, req.Description, req.ParentID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteSuccess(w, category)
}

// UpdateCategoryHandler 处理更新分类请求
// 参数: w - HTTP 响应写入器, r - HTTP 请求
func UpdateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	categoryID, err := GetIntParam(r, "id")
	if err != nil {
		WriteError(w, http.StatusBadRequest, "无效的分类ID")
		return
	}

	var req struct {
		Name        *string `json:"name"`
		Slug        *string `json:"slug"`
		Description *string `json:"description"`
		ParentID    *int64  `json:"parent_id"`
	}

	if err := ParseJSONRequest(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "无效的请求格式")
		return
	}

	category, err := categoryService.UpdateCategory(categoryID, req.Name, req.Slug, req.Description, req.ParentID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if category == nil {
		WriteError(w, http.StatusNotFound, "分类不存在")
		return
	}

	WriteSuccess(w, category)
}

// DeleteCategoryHandler 处理删除分类请求
// 参数: w - HTTP 响应写入器, r - HTTP 请求
func DeleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	categoryID, err := GetIntParam(r, "id")
	if err != nil {
		WriteError(w, http.StatusBadRequest, "无效的分类ID")
		return
	}

	success, err := categoryService.DeleteCategory(categoryID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "删除分类失败")
		return
	}

	if !success {
		WriteError(w, http.StatusNotFound, "分类不存在")
		return
	}

	WriteSuccess(w, map[string]bool{"success": true})
}
