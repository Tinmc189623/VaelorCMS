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
 * 与此程序一起。如果没有, 请见 <https://www.gnu.org/licenses/>.
 */

package services

import (
	"database/sql"
	"time"
	"vaelorcms/internal/database"
	"vaelorcms/internal/models"
)

// ArticleService 文章服务类，提供文章相关的业务逻辑操作
type ArticleService struct{}

// NewArticleService 创建文章服务实例
// 返回: 文章服务指针
func NewArticleService() *ArticleService {
	return &ArticleService{}
}

// ArticleCreate 文章创建请求结构
type ArticleCreate struct {
	Title      string
	Slug       string
	Content    string
	Excerpt    sql.NullString
	Status     string
	AuthorID   int64
	CategoryID sql.NullInt64
	TagIDs     []int64
}

// ArticleUpdate 文章更新请求结构
type ArticleUpdate struct {
	Title      *string
	Slug       *string
	Content    *string
	Excerpt    *sql.NullString
	Status     *string
	CategoryID *sql.NullInt64
	TagIDs     []int64
}

// CreateArticle 创建新文章
// 参数: articleData - 文章创建信息
// 返回: 创建的文章对象指针和可能的错误
func (s *ArticleService) CreateArticle(articleData ArticleCreate) (*models.Article, error) {
	article := &models.Article{
		Title:      articleData.Title,
		Slug:       articleData.Slug,
		Content:    articleData.Content,
		Excerpt:    articleData.Excerpt,
		Status:     articleData.Status,
		AuthorID:   articleData.AuthorID,
		CategoryID: articleData.CategoryID,
	}

	if articleData.Status == "published" {
		now := time.Now()
		article.PublishedAt = sql.NullTime{Time: now, Valid: true}
	}

	id, err := models.CreateArticle(article)
	if err != nil {
		return nil, err
	}
	article.ID = id

	for _, tagID := range articleData.TagIDs {
		if err := models.AddArticleTag(id, tagID); err != nil {
			return nil, err
		}
	}

	return article, nil
}

// GetArticleByID 根据文章ID查询文章
// 参数: articleID - 文章ID
// 返回: 文章对象指针和可能的错误
func (s *ArticleService) GetArticleByID(articleID int64) (*models.Article, error) {
	return models.GetArticleByID(articleID)
}

// GetArticleBySlug 根据文章URL别名查询文章
// 参数: slug - 文章URL别名
// 返回: 文章对象指针和可能的错误
func (s *ArticleService) GetArticleBySlug(slug string) (*models.Article, error) {
	return models.GetArticleBySlug(slug)
}

// GetArticles 获取文章列表，支持过滤和搜索
// 参数: skip - 跳过数量, limit - 返回数量限制, status - 文章状态（可选）, categoryID - 分类ID（可选）, search - 搜索关键词（可选）
// 返回: 文章列表和可能的错误
func (s *ArticleService) GetArticles(skip, limit int, status *string, categoryID *int64, search *string) ([]*models.Article, error) {
	qb := database.NewQueryBuilder((&models.Article{}).TableName())

	if status != nil {
		qb.Where("status = ?", *status)
	}

	if categoryID != nil {
		qb.Where("category_id = ?", *categoryID)
	}

	if search != nil && *search != "" {
		searchPattern := "%" + *search + "%"
		qb.Where("(title LIKE ? OR content LIKE ? OR excerpt LIKE ?)", searchPattern, searchPattern, searchPattern)
	}

	qb.OrderBy("created_at DESC")
	qb.Limit(limit)
	qb.Offset(skip)

	query, args := qb.Build()
	rows, err := database.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []*models.Article
	err = database.ScanRows(rows, &articles)
	return articles, err
}

// CountArticles 统计文章数量，支持过滤和搜索
// 参数: status - 文章状态（可选）, categoryID - 分类ID（可选）, search - 搜索关键词（可选）
// 返回: 文章数量和可能的错误
func (s *ArticleService) CountArticles(status *string, categoryID *int64, search *string) (int64, error) {
	qb := database.NewQueryBuilder((&models.Article{}).TableName())

	if status != nil {
		qb.Where("status = ?", *status)
	}

	if categoryID != nil {
		qb.Where("category_id = ?", *categoryID)
	}

	if search != nil && *search != "" {
		searchPattern := "%" + *search + "%"
		qb.Where("(title LIKE ? OR content LIKE ? OR excerpt LIKE ?)", searchPattern, searchPattern, searchPattern)
	}

	query, args := qb.BuildCount()
	var count int64
	err := database.QueryRow(query, args...).Scan(&count)
	return count, err
}

// UpdateArticle 更新文章信息
// 参数: articleID - 文章ID, articleUpdate - 文章更新信息
// 返回: 更新后的文章对象指针和可能的错误
func (s *ArticleService) UpdateArticle(articleID int64, articleUpdate ArticleUpdate) (*models.Article, error) {
	article, err := models.GetArticleByID(articleID)
	if err != nil {
		return nil, err
	}
	if article == nil {
		return nil, nil
	}

	if articleUpdate.Title != nil {
		article.Title = *articleUpdate.Title
	}
	if articleUpdate.Slug != nil {
		article.Slug = *articleUpdate.Slug
	}
	if articleUpdate.Content != nil {
		article.Content = *articleUpdate.Content
	}
	if articleUpdate.Excerpt != nil {
		article.Excerpt = *articleUpdate.Excerpt
	}
	if articleUpdate.Status != nil {
		if *articleUpdate.Status == "published" && !article.PublishedAt.Valid {
			now := time.Now()
			article.PublishedAt = sql.NullTime{Time: now, Valid: true}
		}
		article.Status = *articleUpdate.Status
	}
	if articleUpdate.CategoryID != nil {
		article.CategoryID = *articleUpdate.CategoryID
	}

	_, err = models.UpdateArticle(article)
	if err != nil {
		return nil, err
	}

	if articleUpdate.TagIDs != nil {
		currentTags, err := models.GetArticleTags(articleID)
		if err != nil {
			return nil, err
		}

		tagMap := make(map[int64]bool)
		for _, tagID := range articleUpdate.TagIDs {
			tagMap[tagID] = true
		}

		for _, tagID := range currentTags {
			if !tagMap[tagID] {
				if err := models.RemoveArticleTag(articleID, tagID); err != nil {
					return nil, err
				}
			}
		}

		for _, tagID := range articleUpdate.TagIDs {
			found := false
			for _, currentTagID := range currentTags {
				if tagID == currentTagID {
					found = true
					break
				}
			}
			if !found {
				if err := models.AddArticleTag(articleID, tagID); err != nil {
					return nil, err
				}
			}
		}
	}

	return article, nil
}

// DeleteArticle 删除文章
// 参数: articleID - 文章ID
// 返回: 删除是否成功和可能的错误
func (s *ArticleService) DeleteArticle(articleID int64) (bool, error) {
	article, err := models.GetArticleByID(articleID)
	if err != nil {
		return false, err
	}
	if article == nil {
		return false, nil
	}

	tags, err := models.GetArticleTags(articleID)
	if err != nil {
		return false, err
	}

	for _, tagID := range tags {
		if err := models.RemoveArticleTag(articleID, tagID); err != nil {
			return false, err
		}
	}

	rows, err := models.DeleteArticle(articleID)
	if err != nil {
		return false, err
	}
	return rows > 0, nil
}

// GetArticleTags 获取文章的标签ID列表
// 参数: articleID - 文章ID
// 返回: 标签ID列表和可能的错误
func (s *ArticleService) GetArticleTags(articleID int64) ([]int64, error) {
	return models.GetArticleTags(articleID)
}

// SetArticleTags 设置文章的标签
// 参数: articleID - 文章ID, tagIDs - 标签ID列表
// 返回: 可能的错误
func (s *ArticleService) SetArticleTags(articleID int64, tagIDs []int64) error {
	currentTags, err := models.GetArticleTags(articleID)
	if err != nil {
		return err
	}

	tagMap := make(map[int64]bool)
	for _, tagID := range tagIDs {
		tagMap[tagID] = true
	}

	for _, tagID := range currentTags {
		if !tagMap[tagID] {
			if err := models.RemoveArticleTag(articleID, tagID); err != nil {
				return err
			}
		}
	}

	for _, tagID := range tagIDs {
		found := false
		for _, currentTagID := range currentTags {
			if tagID == currentTagID {
				found = true
				break
			}
		}
		if !found {
			if err := models.AddArticleTag(articleID, tagID); err != nil {
				return err
			}
		}
	}

	return nil
}
