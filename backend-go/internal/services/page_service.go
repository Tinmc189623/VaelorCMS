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
	"errors"
	"vaelorcms/internal/database"
	"vaelorcms/internal/models"
)

type PageService struct{}

func NewPageService() *PageService {
	return &PageService{}
}

type PageCreate struct {
	Title    string
	Slug     string
	Content  string
	Status   string
	AuthorID int64
}

type PageUpdate struct {
	Title   *string
	Slug    *string
	Content *string
	Status  *string
}

func (s *PageService) CreatePage(pageData PageCreate) (*models.Page, error) {
	existingPage, err := models.GetPageBySlug(pageData.Slug)
	if err != nil {
		return nil, err
	}
	if existingPage != nil {
		return nil, errors.New("页面URL别名已存在")
	}

	page := &models.Page{
		Title:    pageData.Title,
		Slug:     pageData.Slug,
		Content:  pageData.Content,
		Status:   pageData.Status,
		AuthorID: pageData.AuthorID,
	}

	id, err := models.CreatePage(page)
	if err != nil {
		return nil, err
	}
	page.ID = id

	return page, nil
}

func (s *PageService) GetPageByID(pageID int64) (*models.Page, error) {
	return models.GetPageByID(pageID)
}

func (s *PageService) GetPageBySlug(slug string) (*models.Page, error) {
	return models.GetPageBySlug(slug)
}

func (s *PageService) GetPages(skip, limit int, status *string) ([]*models.Page, error) {
	qb := database.NewQueryBuilder((&models.Page{}).TableName())

	if status != nil {
		qb.Where("status = ?", *status)
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

	var pages []*models.Page
	err = database.ScanRows(rows, &pages)
	return pages, err
}

func (s *PageService) CountPages(status *string) (int64, error) {
	qb := database.NewQueryBuilder((&models.Page{}).TableName())

	if status != nil {
		qb.Where("status = ?", *status)
	}

	query, args := qb.BuildCount()
	var count int64
	err := database.QueryRow(query, args...).Scan(&count)
	return count, err
}

func (s *PageService) UpdatePage(pageID int64, pageUpdate PageUpdate) (*models.Page, error) {
	page, err := models.GetPageByID(pageID)
	if err != nil {
		return nil, err
	}
	if page == nil {
		return nil, nil
	}

	if pageUpdate.Title != nil {
		page.Title = *pageUpdate.Title
	}
	if pageUpdate.Slug != nil {
		if *pageUpdate.Slug != page.Slug {
			existingPage, err := models.GetPageBySlug(*pageUpdate.Slug)
			if err != nil {
				return nil, err
			}
			if existingPage != nil {
				return nil, errors.New("页面URL别名已存在")
			}
		}
		page.Slug = *pageUpdate.Slug
	}
	if pageUpdate.Content != nil {
		page.Content = *pageUpdate.Content
	}
	if pageUpdate.Status != nil {
		page.Status = *pageUpdate.Status
	}

	_, err = models.UpdatePage(page)
	if err != nil {
		return nil, err
	}

	return page, nil
}

func (s *PageService) DeletePage(pageID int64) (bool, error) {
	page, err := models.GetPageByID(pageID)
	if err != nil {
		return false, err
	}
	if page == nil {
		return false, nil
	}

	rows, err := models.DeletePage(pageID)
	if err != nil {
		return false, err
	}
	return rows > 0, nil
}

func (s *PageService) GetAllPages() ([]*models.Page, error) {
	return models.GetAllPages()
}

func (s *PageService) GetPagesByAuthorID(authorID int64) ([]*models.Page, error) {
	return models.GetPagesByAuthorID(authorID)
}
