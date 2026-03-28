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

// Content 内容模型，存储各类内容信息
type Content struct {
	ID          int64          `json:"id"`
	Title       string         `json:"title"`
	Slug        string         `json:"slug"`
	Content     string         `json:"content"`
	Excerpt     sql.NullString `json:"excerpt"`
	Status      string         `json:"status"`
	ContentType string         `json:"content_type"`
	AuthorID    int64          `json:"author_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	PublishedAt sql.NullTime   `json:"published_at"`
}

// TableName 返回表名
func (c *Content) TableName() string {
	return "contents"
}

// GetContentByID 根据ID获取内容
// 参数: id - 内容ID
// 返回: 内容对象指针和可能的错误
func GetContentByID(id int64) (*Content, error) {
	var content Content
	found, err := database.FindByID(content.TableName(), id, &content)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return &content, nil
}

// GetContentBySlug 根据Slug获取内容
// 参数: slug - 内容URL别名
// 返回: 内容对象指针和可能的错误
func GetContentBySlug(slug string) (*Content, error) {
	var contents []*Content
	qb := database.NewQueryBuilder((&Content{}).TableName()).Where("slug = ?", slug)
	query, args := qb.Build()

	rows, err := database.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	err = database.ScanRows(rows, &contents)
	if err != nil {
		return nil, err
	}

	if len(contents) == 0 {
		return nil, nil
	}
	return contents[0], nil
}

// CreateContent 创建新内容
// 参数: content - 内容对象
// 返回: 插入的ID和可能的错误
func CreateContent(content *Content) (int64, error) {
	now := time.Now()
	content.CreatedAt = now
	content.UpdatedAt = now

	return database.Create(
		content.TableName(),
		[]string{"title", "slug", "content", "excerpt", "status", "content_type", "author_id", "created_at", "updated_at", "published_at"},
		content.Title,
		content.Slug,
		content.Content,
		content.Excerpt,
		content.Status,
		content.ContentType,
		content.AuthorID,
		content.CreatedAt,
		content.UpdatedAt,
		content.PublishedAt,
	)
}

// UpdateContent 更新内容信息
// 参数: content - 内容对象
// 返回: 影响的行数和可能的错误
func UpdateContent(content *Content) (int64, error) {
	content.UpdatedAt = time.Now()

	updates := map[string]interface{}{
		"title":        content.Title,
		"slug":         content.Slug,
		"content":      content.Content,
		"excerpt":      content.Excerpt,
		"status":       content.Status,
		"content_type": content.ContentType,
		"author_id":    content.AuthorID,
		"updated_at":   content.UpdatedAt,
		"published_at": content.PublishedAt,
	}

	return database.Update(content.TableName(), content.ID, updates)
}

// DeleteContent 删除内容
// 参数: id - 内容ID
// 返回: 影响的行数和可能的错误
func DeleteContent(id int64) (int64, error) {
	return database.Delete((&Content{}).TableName(), id)
}

// GetAllContents 获取所有内容
// 返回: 内容列表和可能的错误
func GetAllContents() ([]*Content, error) {
	var contents []*Content
	err := database.FindAll((&Content{}).TableName(), &contents)
	return contents, err
}

// GetContentsByType 根据类型获取内容
// 参数: contentType - 内容类型
// 返回: 内容列表和可能的错误
func GetContentsByType(contentType string) ([]*Content, error) {
	var contents []*Content
	qb := database.NewQueryBuilder((&Content{}).TableName()).Where("content_type = ?", contentType)
	query, args := qb.Build()

	rows, err := database.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	err = database.ScanRows(rows, &contents)
	return contents, err
}

// GetContentsByAuthorID 根据作者ID获取内容列表
// 参数: authorID - 作者ID
// 返回: 内容列表和可能的错误
func GetContentsByAuthorID(authorID int64) ([]*Content, error) {
	var contents []*Content
	qb := database.NewQueryBuilder((&Content{}).TableName()).Where("author_id = ?", authorID)
	query, args := qb.Build()

	rows, err := database.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	err = database.ScanRows(rows, &contents)
	return contents, err
}
