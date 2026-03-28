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

// Category 分类模型，存储内容分类信息
type Category struct {
	ID          int64          `json:"id"`
	Name        string         `json:"name"`
	Slug        string         `json:"slug"`
	Description sql.NullString `json:"description"`
	ParentID    sql.NullInt64  `json:"parent_id"`
	CreatedAt   time.Time      `json:"created_at"`
}

// TableName 返回表名
func (c *Category) TableName() string {
	return "categories"
}

// GetCategoryByID 根据ID获取分类
// 参数: id - 分类ID
// 返回: 分类对象指针和可能的错误
func GetCategoryByID(id int64) (*Category, error) {
	var category Category
	found, err := database.FindByID(category.TableName(), id, &category)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return &category, nil
}

// GetCategoryBySlug 根据Slug获取分类
// 参数: slug - 分类URL别名
// 返回: 分类对象指针和可能的错误
func GetCategoryBySlug(slug string) (*Category, error) {
	var categories []*Category
	qb := database.NewQueryBuilder((&Category{}).TableName()).Where("slug = ?", slug)
	query, args := qb.Build()

	rows, err := database.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	err = database.ScanRows(rows, &categories)
	if err != nil {
		return nil, err
	}

	if len(categories) == 0 {
		return nil, nil
	}
	return categories[0], nil
}

// CreateCategory 创建新分类
// 参数: category - 分类对象
// 返回: 插入的ID和可能的错误
func CreateCategory(category *Category) (int64, error) {
	category.CreatedAt = time.Now()

	return database.Create(
		category.TableName(),
		[]string{"name", "slug", "description", "parent_id", "created_at"},
		category.Name,
		category.Slug,
		category.Description,
		category.ParentID,
		category.CreatedAt,
	)
}

// UpdateCategory 更新分类信息
// 参数: category - 分类对象
// 返回: 影响的行数和可能的错误
func UpdateCategory(category *Category) (int64, error) {
	updates := map[string]interface{}{
		"name":        category.Name,
		"slug":        category.Slug,
		"description": category.Description,
		"parent_id":   category.ParentID,
	}

	return database.Update(category.TableName(), category.ID, updates)
}

// DeleteCategory 删除分类
// 参数: id - 分类ID
// 返回: 影响的行数和可能的错误
func DeleteCategory(id int64) (int64, error) {
	return database.Delete((&Category{}).TableName(), id)
}

// GetAllCategories 获取所有分类
// 返回: 分类列表和可能的错误
func GetAllCategories() ([]*Category, error) {
	var categories []*Category
	err := database.FindAll((&Category{}).TableName(), &categories)
	return categories, err
}

// GetChildCategories 获取子分类
// 参数: parentID - 父分类ID
// 返回: 子分类列表和可能的错误
func GetChildCategories(parentID int64) ([]*Category, error) {
	var categories []*Category
	qb := database.NewQueryBuilder((&Category{}).TableName()).Where("parent_id = ?", parentID)
	query, args := qb.Build()

	rows, err := database.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	err = database.ScanRows(rows, &categories)
	return categories, err
}
