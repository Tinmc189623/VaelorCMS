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

package models

import (
	"database/sql"
	"time"
	"vaelorcms/internal/database"
)

// Article 文章模型，存储文章内容信息
type Article struct {
	ID          int64          `json:"id"`
	Title       string         `json:"title"`
	Slug        string         `json:"slug"`
	Content     string         `json:"content"`
	Excerpt     sql.NullString `json:"excerpt"`
	Status      string         `json:"status"`
	AuthorID    int64          `json:"author_id"`
	CategoryID  sql.NullInt64  `json:"category_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	PublishedAt sql.NullTime   `json:"published_at"`
}

// TableName 返回表名
func (a *Article) TableName() string {
	return "articles"
}

// ArticleTag 文章与标签的多对多关联表
type ArticleTag struct {
	ArticleID int64 `json:"article_id"`
	TagID     int64 `json:"tag_id"`
}

// TableName 返回表名
func (at *ArticleTag) TableName() string {
	return "article_tag"
}

// GetArticleByID 根据ID获取文章
// 参数: id - 文章ID
// 返回: 文章对象指针和可能的错误
func GetArticleByID(id int64) (*Article, error) {
	var article Article
	found, err := database.FindByID(article.TableName(), id, &article)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return &article, nil
}

// GetArticleBySlug 根据Slug获取文章
// 参数: slug - 文章URL别名
// 返回: 文章对象指针和可能的错误
func GetArticleBySlug(slug string) (*Article, error) {
	var articles []*Article
	qb := database.NewQueryBuilder((&Article{}).TableName()).Where("slug = ?", slug)
	query, args := qb.Build()

	rows, err := database.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	err = database.ScanRows(rows, &articles)
	if err != nil {
		return nil, err
	}

	if len(articles) == 0 {
		return nil, nil
	}
	return articles[0], nil
}

// CreateArticle 创建新文章
// 参数: article - 文章对象
// 返回: 插入的ID和可能的错误
func CreateArticle(article *Article) (int64, error) {
	now := time.Now()
	article.CreatedAt = now
	article.UpdatedAt = now

	return database.Create(
		article.TableName(),
		[]string{"title", "slug", "content", "excerpt", "status", "author_id", "category_id", "created_at", "updated_at", "published_at"},
		article.Title,
		article.Slug,
		article.Content,
		article.Excerpt,
		article.Status,
		article.AuthorID,
		article.CategoryID,
		article.CreatedAt,
		article.UpdatedAt,
		article.PublishedAt,
	)
}

// UpdateArticle 更新文章信息
// 参数: article - 文章对象
// 返回: 影响的行数和可能的错误
func UpdateArticle(article *Article) (int64, error) {
	article.UpdatedAt = time.Now()

	updates := map[string]interface{}{
		"title":        article.Title,
		"slug":         article.Slug,
		"content":      article.Content,
		"excerpt":      article.Excerpt,
		"status":       article.Status,
		"author_id":    article.AuthorID,
		"category_id":  article.CategoryID,
		"updated_at":   article.UpdatedAt,
		"published_at": article.PublishedAt,
	}

	return database.Update(article.TableName(), article.ID, updates)
}

// DeleteArticle 删除文章
// 参数: id - 文章ID
// 返回: 影响的行数和可能的错误
func DeleteArticle(id int64) (int64, error) {
	return database.Delete((&Article{}).TableName(), id)
}

// GetAllArticles 获取所有文章
// 返回: 文章列表和可能的错误
func GetAllArticles() ([]*Article, error) {
	var articles []*Article
	err := database.FindAll((&Article{}).TableName(), &articles)
	return articles, err
}

// GetArticlesByAuthorID 根据作者ID获取文章列表
// 参数: authorID - 作者ID
// 返回: 文章列表和可能的错误
func GetArticlesByAuthorID(authorID int64) ([]*Article, error) {
	var articles []*Article
	qb := database.NewQueryBuilder((&Article{}).TableName()).Where("author_id = ?", authorID)
	query, args := qb.Build()

	rows, err := database.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	err = database.ScanRows(rows, &articles)
	return articles, err
}

// AddArticleTag 添加文章标签关联
// 参数: articleID - 文章ID, tagID - 标签ID
// 返回: 可能的错误
func AddArticleTag(articleID, tagID int64) error {
	_, err := database.Create(
		(&ArticleTag{}).TableName(),
		[]string{"article_id", "tag_id"},
		articleID,
		tagID,
	)
	return err
}

// RemoveArticleTag 移除文章标签关联
// 参数: articleID - 文章ID, tagID - 标签ID
// 返回: 可能的错误
func RemoveArticleTag(articleID, tagID int64) error {
	db := database.NewDeleteBuilder((&ArticleTag{}).TableName())
	db.Where("article_id = ?", articleID)
	db.Where("tag_id = ?", tagID)
	query, args := db.Build()

	_, err := database.Exec(query, args...)
	return err
}

// GetArticleTags 获取文章的标签ID列表
// 参数: articleID - 文章ID
// 返回: 标签ID列表和可能的错误
func GetArticleTags(articleID int64) ([]int64, error) {
	var tagIDs []int64
	qb := database.NewQueryBuilder((&ArticleTag{}).TableName()).Select("tag_id").Where("article_id = ?", articleID)
	query, args := qb.Build()

	rows, err := database.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tagID int64
		if err := rows.Scan(&tagID); err != nil {
			return nil, err
		}
		tagIDs = append(tagIDs, tagID)
	}

	return tagIDs, rows.Err()
}
