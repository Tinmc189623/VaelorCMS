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

var contentService = services.NewContentService()

// RegisterContentRoutes 注册内容相关路由
// 参数: router - Vaelor Core 路由器
func RegisterContentRoutes(router *vaelorcore.Router) {
	router.GET("/api/content", ListContentHandler)
	router.GET("/api/content/:id", GetContentHandler)
	router.GET("/api/content/slug/:slug", GetContentBySlugHandler)
	router.POST("/api/content", CreateContentHandler)
	router.PUT("/api/content/:id", UpdateContentHandler)
	router.DELETE("/api/content/:id", DeleteContentHandler)
}

// ListContentHandler 处理获取内容列表请求
// 参数: w - HTTP 响应写入器, r - HTTP 请求
func ListContentHandler(w http.ResponseWriter, r *http.Request) {
	skip := GetIntQueryParam(r, "skip", 0)
	limit := GetIntQueryParam(r, "limit", 20)
	statusStr := GetQueryParam(r, "status", "")
	contentTypeStr := GetQueryParam(r, "content_type", "")
	searchStr := GetQueryParam(r, "search", "")

	var status *string
	if statusStr != "" {
		status = &statusStr
	}

	var contentType *string
	if contentTypeStr != "" {
		contentType = &contentTypeStr
	}

	var search *string
	if searchStr != "" {
		search = &searchStr
	}

	contents, err := contentService.GetContents(skip, limit, status, contentType, search)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "获取内容列表失败")
		return
	}

	total, err := contentService.CountContents(status, contentType, search)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "统计内容数量失败")
		return
	}

	WriteSuccess(w, map[string]interface{}{
		"content": contents,
		"total":   total,
		"skip":    skip,
		"limit":   limit,
	})
}

// GetContentHandler 处理根据ID获取内容请求
// 参数: w - HTTP 响应写入器, r - HTTP 请求
func GetContentHandler(w http.ResponseWriter, r *http.Request) {
	contentID, err := GetIntParam(r, "id")
	if err != nil {
		WriteError(w, http.StatusBadRequest, "无效的内容ID")
		return
	}

	content, err := contentService.GetContentByID(contentID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "获取内容失败")
		return
	}

	if content == nil {
		WriteError(w, http.StatusNotFound, "内容不存在")
		return
	}

	WriteSuccess(w, content)
}

// GetContentBySlugHandler 处理根据Slug获取内容请求
// 参数: w - HTTP 响应写入器, r - HTTP 请求
func GetContentBySlugHandler(w http.ResponseWriter, r *http.Request) {
	slug := vaelorcore.GetParam(r, "slug")
	if slug == "" {
		WriteError(w, http.StatusBadRequest, "无效的内容Slug")
		return
	}

	content, err := contentService.GetContentBySlug(slug)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "获取内容失败")
		return
	}

	if content == nil {
		WriteError(w, http.StatusNotFound, "内容不存在")
		return
	}

	WriteSuccess(w, content)
}

// CreateContentHandler 处理创建内容请求
// 参数: w - HTTP 响应写入器, r - HTTP 请求
func CreateContentHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title       string  `json:"title"`
		Slug        string  `json:"slug"`
		Content     string  `json:"content"`
		Excerpt     *string `json:"excerpt"`
		Status      string  `json:"status"`
		ContentType string  `json:"content_type"`
		AuthorID    int64   `json:"author_id"`
	}

	if err := ParseJSONRequest(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "无效的请求格式")
		return
	}

	if req.Title == "" || req.Slug == "" || req.Content == "" || req.Status == "" || req.ContentType == "" || req.AuthorID == 0 {
		WriteError(w, http.StatusBadRequest, "标题、Slug、内容、状态、内容类型和作者ID不能为空")
		return
	}

	var excerpt sql.NullString
	if req.Excerpt != nil {
		excerpt = sql.NullString{String: *req.Excerpt, Valid: true}
	}

	contentData := services.ContentCreate{
		Title:       req.Title,
		Slug:        req.Slug,
		Content:     req.Content,
		Excerpt:     excerpt,
		Status:      req.Status,
		ContentType: req.ContentType,
		AuthorID:    req.AuthorID,
	}

	content, err := contentService.CreateContent(contentData)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteSuccess(w, content)
}

// UpdateContentHandler 处理更新内容请求
// 参数: w - HTTP 响应写入器, r - HTTP 请求
func UpdateContentHandler(w http.ResponseWriter, r *http.Request) {
	contentID, err := GetIntParam(r, "id")
	if err != nil {
		WriteError(w, http.StatusBadRequest, "无效的内容ID")
		return
	}

	var req struct {
		Title       *string `json:"title"`
		Slug        *string `json:"slug"`
		Content     *string `json:"content"`
		Excerpt     *string `json:"excerpt"`
		Status      *string `json:"status"`
		ContentType *string `json:"content_type"`
	}

	if err := ParseJSONRequest(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "无效的请求格式")
		return
	}

	var updateData services.ContentUpdate

	if req.Title != nil {
		updateData.Title = req.Title
	}
	if req.Slug != nil {
		updateData.Slug = req.Slug
	}
	if req.Content != nil {
		updateData.Content = req.Content
	}
	if req.Excerpt != nil {
		excerpt := sql.NullString{String: *req.Excerpt, Valid: true}
		updateData.Excerpt = &excerpt
	}
	if req.Status != nil {
		updateData.Status = req.Status
	}
	if req.ContentType != nil {
		updateData.ContentType = req.ContentType
	}

	content, err := contentService.UpdateContent(contentID, updateData)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if content == nil {
		WriteError(w, http.StatusNotFound, "内容不存在")
		return
	}

	WriteSuccess(w, content)
}

// DeleteContentHandler 处理删除内容请求
// 参数: w - HTTP 响应写入器, r - HTTP 请求
func DeleteContentHandler(w http.ResponseWriter, r *http.Request) {
	contentID, err := GetIntParam(r, "id")
	if err != nil {
		WriteError(w, http.StatusBadRequest, "无效的内容ID")
		return
	}

	success, err := contentService.DeleteContent(contentID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "删除内容失败")
		return
	}

	if !success {
		WriteError(w, http.StatusNotFound, "内容不存在")
		return
	}

	WriteSuccess(w, map[string]bool{"success": true})
}
