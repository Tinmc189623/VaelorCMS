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

var pageService = services.NewPageService()

// RegisterPageRoutes 注册页面相关路由
// 参数: router - Vaelor Core 路由器
func RegisterPageRoutes(router *vaelorcore.Router) {
	router.GET("/api/pages", ListPagesHandler)
	router.GET("/api/pages/:id", GetPageHandler)
	router.GET("/api/pages/slug/:slug", GetPageBySlugHandler)
	router.POST("/api/pages", CreatePageHandler)
	router.PUT("/api/pages/:id", UpdatePageHandler)
	router.DELETE("/api/pages/:id", DeletePageHandler)
}

// ListPagesHandler 处理获取页面列表请求
// 参数: w - HTTP 响应写入器, r - HTTP 请求
func ListPagesHandler(w http.ResponseWriter, r *http.Request) {
	skip := GetIntQueryParam(r, "skip", 0)
	limit := GetIntQueryParam(r, "limit", 20)
	statusStr := GetQueryParam(r, "status", "")

	var status *string
	if statusStr != "" {
		status = &statusStr
	}

	pages, err := pageService.GetPages(skip, limit, status)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "获取页面列表失败")
		return
	}

	total, err := pageService.CountPages(status)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "统计页面数量失败")
		return
	}

	WriteSuccess(w, map[string]interface{}{
		"pages": pages,
		"total": total,
		"skip":  skip,
		"limit": limit,
	})
}

// GetPageHandler 处理根据ID获取页面请求
// 参数: w - HTTP 响应写入器, r - HTTP 请求
func GetPageHandler(w http.ResponseWriter, r *http.Request) {
	pageID, err := GetIntParam(r, "id")
	if err != nil {
		WriteError(w, http.StatusBadRequest, "无效的页面ID")
		return
	}

	page, err := pageService.GetPageByID(pageID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "获取页面失败")
		return
	}

	if page == nil {
		WriteError(w, http.StatusNotFound, "页面不存在")
		return
	}

	WriteSuccess(w, page)
}

// GetPageBySlugHandler 处理根据Slug获取页面请求
// 参数: w - HTTP 响应写入器, r - HTTP 请求
func GetPageBySlugHandler(w http.ResponseWriter, r *http.Request) {
	slug := vaelorcore.GetParam(r, "slug")
	if slug == "" {
		WriteError(w, http.StatusBadRequest, "无效的页面Slug")
		return
	}

	page, err := pageService.GetPageBySlug(slug)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "获取页面失败")
		return
	}

	if page == nil {
		WriteError(w, http.StatusNotFound, "页面不存在")
		return
	}

	WriteSuccess(w, page)
}

// CreatePageHandler 处理创建页面请求
// 参数: w - HTTP 响应写入器, r - HTTP 请求
func CreatePageHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title    string `json:"title"`
		Slug     string `json:"slug"`
		Content  string `json:"content"`
		Status   string `json:"status"`
		AuthorID int64  `json:"author_id"`
	}

	if err := ParseJSONRequest(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "无效的请求格式")
		return
	}

	if req.Title == "" || req.Slug == "" || req.Content == "" || req.Status == "" || req.AuthorID == 0 {
		WriteError(w, http.StatusBadRequest, "标题、Slug、内容、状态和作者ID不能为空")
		return
	}

	pageData := services.PageCreate{
		Title:    req.Title,
		Slug:     req.Slug,
		Content:  req.Content,
		Status:   req.Status,
		AuthorID: req.AuthorID,
	}

	page, err := pageService.CreatePage(pageData)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteSuccess(w, page)
}

// UpdatePageHandler 处理更新页面请求
// 参数: w - HTTP 响应写入器, r - HTTP 请求
func UpdatePageHandler(w http.ResponseWriter, r *http.Request) {
	pageID, err := GetIntParam(r, "id")
	if err != nil {
		WriteError(w, http.StatusBadRequest, "无效的页面ID")
		return
	}

	var req struct {
		Title   *string `json:"title"`
		Slug    *string `json:"slug"`
		Content *string `json:"content"`
		Status  *string `json:"status"`
	}

	if err := ParseJSONRequest(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "无效的请求格式")
		return
	}

	var updateData services.PageUpdate

	if req.Title != nil {
		updateData.Title = req.Title
	}
	if req.Slug != nil {
		updateData.Slug = req.Slug
	}
	if req.Content != nil {
		updateData.Content = req.Content
	}
	if req.Status != nil {
		updateData.Status = req.Status
	}

	page, err := pageService.UpdatePage(pageID, updateData)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if page == nil {
		WriteError(w, http.StatusNotFound, "页面不存在")
		return
	}

	WriteSuccess(w, page)
}

// DeletePageHandler 处理删除页面请求
// 参数: w - HTTP 响应写入器, r - HTTP 请求
func DeletePageHandler(w http.ResponseWriter, r *http.Request) {
	pageID, err := GetIntParam(r, "id")
	if err != nil {
		WriteError(w, http.StatusBadRequest, "无效的页面ID")
		return
	}

	success, err := pageService.DeletePage(pageID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "删除页面失败")
		return
	}

	if !success {
		WriteError(w, http.StatusNotFound, "页面不存在")
		return
	}

	WriteSuccess(w, map[string]bool{"success": true})
}
