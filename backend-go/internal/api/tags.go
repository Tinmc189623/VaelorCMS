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
	"vaelorcms/internal/services"
	"vaelorcms/pkg/vaelorcore"
)

var tagService = services.NewTagService()

// RegisterTagRoutes 注册标签相关路由
// 参数: router - Vaelor Core 路由器
func RegisterTagRoutes(router *vaelorcore.Router) {
	router.GET("/api/tags", ListTagsHandler)
	router.GET("/api/tags/:id", GetTagHandler)
	router.GET("/api/tags/slug/:slug", GetTagBySlugHandler)
	router.POST("/api/tags", CreateTagHandler)
	router.PUT("/api/tags/:id", UpdateTagHandler)
	router.DELETE("/api/tags/:id", DeleteTagHandler)
}

// ListTagsHandler 处理获取标签列表请求
// 参数: w - HTTP 响应写入器, r - HTTP 请求
func ListTagsHandler(w http.ResponseWriter, r *http.Request) {
	skip := GetIntQueryParam(r, "skip", 0)
	limit := GetIntQueryParam(r, "limit", 20)
	search := GetQueryParam(r, "search", "")

	tags, err := tagService.GetTags(skip, limit, search)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "获取标签列表失败")
		return
	}

	total, err := tagService.CountTags(search)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "统计标签数量失败")
		return
	}

	WriteSuccess(w, map[string]interface{}{
		"tags":  tags,
		"total": total,
		"skip":  skip,
		"limit": limit,
	})
}

// GetTagHandler 处理根据ID获取标签请求
// 参数: w - HTTP 响应写入器, r - HTTP 请求
func GetTagHandler(w http.ResponseWriter, r *http.Request) {
	tagID, err := GetIntParam(r, "id")
	if err != nil {
		WriteError(w, http.StatusBadRequest, "无效的标签ID")
		return
	}

	tag, err := tagService.GetTagByID(tagID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "获取标签失败")
		return
	}

	if tag == nil {
		WriteError(w, http.StatusNotFound, "标签不存在")
		return
	}

	WriteSuccess(w, tag)
}

// GetTagBySlugHandler 处理根据Slug获取标签请求
// 参数: w - HTTP 响应写入器, r - HTTP 请求
func GetTagBySlugHandler(w http.ResponseWriter, r *http.Request) {
	slug := vaelorcore.GetParam(r, "slug")
	if slug == "" {
		WriteError(w, http.StatusBadRequest, "无效的标签Slug")
		return
	}

	tag, err := tagService.GetTagBySlug(slug)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "获取标签失败")
		return
	}

	if tag == nil {
		WriteError(w, http.StatusNotFound, "标签不存在")
		return
	}

	WriteSuccess(w, tag)
}

// CreateTagHandler 处理创建标签请求
// 参数: w - HTTP 响应写入器, r - HTTP 请求
func CreateTagHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
	}

	if err := ParseJSONRequest(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "无效的请求格式")
		return
	}

	if req.Name == "" || req.Slug == "" {
		WriteError(w, http.StatusBadRequest, "标签名称和Slug不能为空")
		return
	}

	tag, err := tagService.CreateTag(req.Name, req.Slug)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteSuccess(w, tag)
}

// UpdateTagHandler 处理更新标签请求
// 参数: w - HTTP 响应写入器, r - HTTP 请求
func UpdateTagHandler(w http.ResponseWriter, r *http.Request) {
	tagID, err := GetIntParam(r, "id")
	if err != nil {
		WriteError(w, http.StatusBadRequest, "无效的标签ID")
		return
	}

	var req struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
	}

	if err := ParseJSONRequest(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "无效的请求格式")
		return
	}

	if req.Name == "" || req.Slug == "" {
		WriteError(w, http.StatusBadRequest, "标签名称和Slug不能为空")
		return
	}

	tag, err := tagService.UpdateTag(tagID, req.Name, req.Slug)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if tag == nil {
		WriteError(w, http.StatusNotFound, "标签不存在")
		return
	}

	WriteSuccess(w, tag)
}

// DeleteTagHandler 处理删除标签请求
// 参数: w - HTTP 响应写入器, r - HTTP 请求
func DeleteTagHandler(w http.ResponseWriter, r *http.Request) {
	tagID, err := GetIntParam(r, "id")
	if err != nil {
		WriteError(w, http.StatusBadRequest, "无效的标签ID")
		return
	}

	success, err := tagService.DeleteTag(tagID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "删除标签失败")
		return
	}

	if !success {
		WriteError(w, http.StatusNotFound, "标签不存在")
		return
	}

	WriteSuccess(w, map[string]bool{"success": true})
}
