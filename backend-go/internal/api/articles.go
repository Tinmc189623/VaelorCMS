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

var articleService = services.NewArticleService()

// RegisterArticleRoutes 注册文章相关路由
// 参数: router - Vaelor Core 路由器
func RegisterArticleRoutes(router *vaelorcore.Router) {
	router.GET("/api/articles", ListArticlesHandler)
	router.GET("/api/articles/:id", GetArticleHandler)
	router.GET("/api/articles/slug/:slug", GetArticleBySlugHandler)
	router.POST("/api/articles", CreateArticleHandler)
	router.PUT("/api/articles/:id", UpdateArticleHandler)
	router.DELETE("/api/articles/:id", DeleteArticleHandler)
	router.GET("/api/articles/:id/tags", GetArticleTagsHandler)
	router.PUT("/api/articles/:id/tags", SetArticleTagsHandler)
}

// ListArticlesHandler 处理获取文章列表请求
// 参数: w - HTTP 响应写入器, r - HTTP 请求
func ListArticlesHandler(w http.ResponseWriter, r *http.Request) {
	skip := GetIntQueryParam(r, "skip", 0)
	limit := GetIntQueryParam(r, "limit", 20)
	statusStr := GetQueryParam(r, "status", "")
	search := GetQueryParam(r, "search", "")

	var status *string
	if statusStr != "" {
		status = &statusStr
	}

	var categoryID *int64
	categoryIDStr := GetQueryParam(r, "category_id", "")
	if categoryIDStr != "" {
		var id int64
		_, err := GetIntParam(r, "category_id")
		if err == nil {
			id, _ = GetIntParam(r, "category_id")
			categoryID = &id
		}
	}

	articles, err := articleService.GetArticles(skip, limit, status, categoryID, &search)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "获取文章列表失败")
		return
	}

	total, err := articleService.CountArticles(status, categoryID, &search)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "统计文章数量失败")
		return
	}

	WriteSuccess(w, map[string]interface{}{
		"articles": articles,
		"total":    total,
		"skip":     skip,
		"limit":    limit,
	})
}

// GetArticleHandler 处理根据ID获取文章请求
// 参数: w - HTTP 响应写入器, r - HTTP 请求
func GetArticleHandler(w http.ResponseWriter, r *http.Request) {
	articleID, err := GetIntParam(r, "id")
	if err != nil {
		WriteError(w, http.StatusBadRequest, "无效的文章ID")
		return
	}

	article, err := articleService.GetArticleByID(articleID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "获取文章失败")
		return
	}

	if article == nil {
		WriteError(w, http.StatusNotFound, "文章不存在")
		return
	}

	WriteSuccess(w, article)
}

// GetArticleBySlugHandler 处理根据Slug获取文章请求
// 参数: w - HTTP 响应写入器, r - HTTP 请求
func GetArticleBySlugHandler(w http.ResponseWriter, r *http.Request) {
	slug := vaelorcore.GetParam(r, "slug")
	if slug == "" {
		WriteError(w, http.StatusBadRequest, "无效的文章Slug")
		return
	}

	article, err := articleService.GetArticleBySlug(slug)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "获取文章失败")
		return
	}

	if article == nil {
		WriteError(w, http.StatusNotFound, "文章不存在")
		return
	}

	WriteSuccess(w, article)
}

// CreateArticleHandler 处理创建文章请求
// 参数: w - HTTP 响应写入器, r - HTTP 请求
func CreateArticleHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title      string  `json:"title"`
		Slug       string  `json:"slug"`
		Content    string  `json:"content"`
		Excerpt    *string `json:"excerpt"`
		Status     string  `json:"status"`
		AuthorID   int64   `json:"author_id"`
		CategoryID *int64  `json:"category_id"`
		TagIDs     []int64 `json:"tag_ids"`
	}

	if err := ParseJSONRequest(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "无效的请求格式")
		return
	}

	if req.Title == "" || req.Slug == "" || req.Content == "" || req.Status == "" || req.AuthorID == 0 {
		WriteError(w, http.StatusBadRequest, "标题、Slug、内容、状态和作者ID不能为空")
		return
	}

	var excerpt sql.NullString
	if req.Excerpt != nil {
		excerpt = sql.NullString{String: *req.Excerpt, Valid: true}
	}

	var categoryID sql.NullInt64
	if req.CategoryID != nil {
		categoryID = sql.NullInt64{Int64: *req.CategoryID, Valid: true}
	}

	articleData := services.ArticleCreate{
		Title:      req.Title,
		Slug:       req.Slug,
		Content:    req.Content,
		Excerpt:    excerpt,
		Status:     req.Status,
		AuthorID:   req.AuthorID,
		CategoryID: categoryID,
		TagIDs:     req.TagIDs,
	}

	article, err := articleService.CreateArticle(articleData)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteSuccess(w, article)
}

// UpdateArticleHandler 处理更新文章请求
// 参数: w - HTTP 响应写入器, r - HTTP 请求
func UpdateArticleHandler(w http.ResponseWriter, r *http.Request) {
	articleID, err := GetIntParam(r, "id")
	if err != nil {
		WriteError(w, http.StatusBadRequest, "无效的文章ID")
		return
	}

	var req struct {
		Title      *string `json:"title"`
		Slug       *string `json:"slug"`
		Content    *string `json:"content"`
		Excerpt    *string `json:"excerpt"`
		Status     *string `json:"status"`
		CategoryID *int64  `json:"category_id"`
		TagIDs     []int64 `json:"tag_ids"`
	}

	if err := ParseJSONRequest(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "无效的请求格式")
		return
	}

	var updateData services.ArticleUpdate

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
	if req.CategoryID != nil {
		catID := sql.NullInt64{Int64: *req.CategoryID, Valid: true}
		updateData.CategoryID = &catID
	}
	if req.TagIDs != nil {
		updateData.TagIDs = req.TagIDs
	}

	article, err := articleService.UpdateArticle(articleID, updateData)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if article == nil {
		WriteError(w, http.StatusNotFound, "文章不存在")
		return
	}

	WriteSuccess(w, article)
}

// DeleteArticleHandler 处理删除文章请求
// 参数: w - HTTP 响应写入器, r - HTTP 请求
func DeleteArticleHandler(w http.ResponseWriter, r *http.Request) {
	articleID, err := GetIntParam(r, "id")
	if err != nil {
		WriteError(w, http.StatusBadRequest, "无效的文章ID")
		return
	}

	success, err := articleService.DeleteArticle(articleID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "删除文章失败")
		return
	}

	if !success {
		WriteError(w, http.StatusNotFound, "文章不存在")
		return
	}

	WriteSuccess(w, map[string]bool{"success": true})
}

// GetArticleTagsHandler 处理获取文章标签请求
// 参数: w - HTTP 响应写入器, r - HTTP 请求
func GetArticleTagsHandler(w http.ResponseWriter, r *http.Request) {
	articleID, err := GetIntParam(r, "id")
	if err != nil {
		WriteError(w, http.StatusBadRequest, "无效的文章ID")
		return
	}

	tagIDs, err := articleService.GetArticleTags(articleID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "获取文章标签失败")
		return
	}

	WriteSuccess(w, map[string][]int64{"tag_ids": tagIDs})
}

// SetArticleTagsHandler 处理设置文章标签请求
// 参数: w - HTTP 响应写入器, r - HTTP 请求
func SetArticleTagsHandler(w http.ResponseWriter, r *http.Request) {
	articleID, err := GetIntParam(r, "id")
	if err != nil {
		WriteError(w, http.StatusBadRequest, "无效的文章ID")
		return
	}

	var req struct {
		TagIDs []int64 `json:"tag_ids"`
	}

	if err := ParseJSONRequest(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "无效的请求格式")
		return
	}

	err = articleService.SetArticleTags(articleID, req.TagIDs)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "设置文章标签失败")
		return
	}

	WriteSuccess(w, map[string]bool{"success": true})
}
