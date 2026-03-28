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
	"errors"
	"time"
	"vaelorcms/internal/database"
	"vaelorcms/internal/models"
)

type ContentService struct{}

func NewContentService() *ContentService {
	return &ContentService{}
}

type ContentCreate struct {
	Title       string
	Slug        string
	Content     string
	Excerpt     sql.NullString
	Status      string
	ContentType string
	AuthorID    int64
}

type ContentUpdate struct {
	Title       *string
	Slug        *string
	Content     *string
	Excerpt     *sql.NullString
	Status      *string
	ContentType *string
}

func (s *ContentService) CreateContent(contentData ContentCreate) (*models.Content, error) {
	existingContent, err := models.GetContentBySlug(contentData.Slug)
	if err != nil {
		return nil, err
	}
	if existingContent != nil {
		return nil, errors.New("内容URL别名已存在")
	}

	content := &models.Content{
		Title:       contentData.Title,
		Slug:        contentData.Slug,
		Content:     contentData.Content,
		Excerpt:     contentData.Excerpt,
		Status:      contentData.Status,
		ContentType: contentData.ContentType,
		AuthorID:    contentData.AuthorID,
	}

	if contentData.Status == "published" {
		now := time.Now()
		content.PublishedAt = sql.NullTime{Time: now, Valid: true}
	}

	id, err := models.CreateContent(content)
	if err != nil {
		return nil, err
	}
	content.ID = id

	return content, nil
}

func (s *ContentService) GetContentByID(contentID int64) (*models.Content, error) {
	return models.GetContentByID(contentID)
}

func (s *ContentService) GetContentBySlug(slug string) (*models.Content, error) {
	return models.GetContentBySlug(slug)
}

func (s *ContentService) GetContents(skip, limit int, status *string, contentType *string, search *string) ([]*models.Content, error) {
	qb := database.NewQueryBuilder((&models.Content{}).TableName())

	if status != nil {
		qb.Where("status = ?", *status)
	}

	if contentType != nil {
		qb.Where("content_type = ?", *contentType)
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

	var contents []*models.Content
	err = database.ScanRows(rows, &contents)
	return contents, err
}

func (s *ContentService) CountContents(status *string, contentType *string, search *string) (int64, error) {
	qb := database.NewQueryBuilder((&models.Content{}).TableName())

	if status != nil {
		qb.Where("status = ?", *status)
	}

	if contentType != nil {
		qb.Where("content_type = ?", *contentType)
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

func (s *ContentService) UpdateContent(contentID int64, contentUpdate ContentUpdate) (*models.Content, error) {
	content, err := models.GetContentByID(contentID)
	if err != nil {
		return nil, err
	}
	if content == nil {
		return nil, nil
	}

	if contentUpdate.Title != nil {
		content.Title = *contentUpdate.Title
	}
	if contentUpdate.Slug != nil {
		if *contentUpdate.Slug != content.Slug {
			existingContent, err := models.GetContentBySlug(*contentUpdate.Slug)
			if err != nil {
				return nil, err
			}
			if existingContent != nil {
				return nil, errors.New("内容URL别名已存在")
			}
		}
		content.Slug = *contentUpdate.Slug
	}
	if contentUpdate.Content != nil {
		content.Content = *contentUpdate.Content
	}
	if contentUpdate.Excerpt != nil {
		content.Excerpt = *contentUpdate.Excerpt
	}
	if contentUpdate.Status != nil {
		if *contentUpdate.Status == "published" && !content.PublishedAt.Valid {
			now := time.Now()
			content.PublishedAt = sql.NullTime{Time: now, Valid: true}
		}
		content.Status = *contentUpdate.Status
	}
	if contentUpdate.ContentType != nil {
		content.ContentType = *contentUpdate.ContentType
	}

	_, err = models.UpdateContent(content)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func (s *ContentService) DeleteContent(contentID int64) (bool, error) {
	content, err := models.GetContentByID(contentID)
	if err != nil {
		return false, err
	}
	if content == nil {
		return false, nil
	}

	rows, err := models.DeleteContent(contentID)
	if err != nil {
		return false, err
	}
	return rows > 0, nil
}

func (s *ContentService) GetAllContents() ([]*models.Content, error) {
	return models.GetAllContents()
}

func (s *ContentService) GetContentsByType(contentType string) ([]*models.Content, error) {
	return models.GetContentsByType(contentType)
}

func (s *ContentService) GetContentsByAuthorID(authorID int64) ([]*models.Content, error) {
	return models.GetContentsByAuthorID(authorID)
}
