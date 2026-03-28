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
	"time"
	"vaelorcms/internal/database"
)

// Page 页面模型，存储静态页面信息
type Page struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Slug      string    `json:"slug"`
	Content   string    `json:"content"`
	Status    string    `json:"status"`
	AuthorID  int64     `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 返回表名
func (p *Page) TableName() string {
	return "pages"
}

// GetPageByID 根据ID获取页面
// 参数: id - 页面ID
// 返回: 页面对象指针和可能的错误
func GetPageByID(id int64) (*Page, error) {
	var page Page
	found, err := database.FindByID(page.TableName(), id, &page)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return &page, nil
}

// GetPageBySlug 根据Slug获取页面
// 参数: slug - 页面URL别名
// 返回: 页面对象指针和可能的错误
func GetPageBySlug(slug string) (*Page, error) {
	var pages []*Page
	qb := database.NewQueryBuilder((&Page{}).TableName()).Where("slug = ?", slug)
	query, args := qb.Build()

	rows, err := database.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	err = database.ScanRows(rows, &pages)
	if err != nil {
		return nil, err
	}

	if len(pages) == 0 {
		return nil, nil
	}
	return pages[0], nil
}

// CreatePage 创建新页面
// 参数: page - 页面对象
// 返回: 插入的ID和可能的错误
func CreatePage(page *Page) (int64, error) {
	now := time.Now()
	page.CreatedAt = now
	page.UpdatedAt = now

	return database.Create(
		page.TableName(),
		[]string{"title", "slug", "content", "status", "author_id", "created_at", "updated_at"},
		page.Title,
		page.Slug,
		page.Content,
		page.Status,
		page.AuthorID,
		page.CreatedAt,
		page.UpdatedAt,
	)
}

// UpdatePage 更新页面信息
// 参数: page - 页面对象
// 返回: 影响的行数和可能的错误
func UpdatePage(page *Page) (int64, error) {
	page.UpdatedAt = time.Now()

	updates := map[string]interface{}{
		"title":      page.Title,
		"slug":       page.Slug,
		"content":    page.Content,
		"status":     page.Status,
		"author_id":  page.AuthorID,
		"updated_at": page.UpdatedAt,
	}

	return database.Update(page.TableName(), page.ID, updates)
}

// DeletePage 删除页面
// 参数: id - 页面ID
// 返回: 影响的行数和可能的错误
func DeletePage(id int64) (int64, error) {
	return database.Delete((&Page{}).TableName(), id)
}

// GetAllPages 获取所有页面
// 返回: 页面列表和可能的错误
func GetAllPages() ([]*Page, error) {
	var pages []*Page
	err := database.FindAll((&Page{}).TableName(), &pages)
	return pages, err
}

// GetPagesByAuthorID 根据作者ID获取页面列表
// 参数: authorID - 作者ID
// 返回: 页面列表和可能的错误
func GetPagesByAuthorID(authorID int64) ([]*Page, error) {
	var pages []*Page
	qb := database.NewQueryBuilder((&Page{}).TableName()).Where("author_id = ?", authorID)
	query, args := qb.Build()

	rows, err := database.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	err = database.ScanRows(rows, &pages)
	return pages, err
}
