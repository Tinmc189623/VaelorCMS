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
	"vaelorcms/internal/database"
	"vaelorcms/internal/models"
)

// CategoryService 分类服务类，提供分类相关的业务逻辑操作
type CategoryService struct{}

// NewCategoryService 创建分类服务实例
// 返回: CategoryService 指针
func NewCategoryService() *CategoryService {
	return &CategoryService{}
}

// CreateCategory 创建新分类
// 参数: name - 分类名称, slug - 分类URL别名, description - 分类描述, parentID - 父分类ID（可选）
// 返回: 创建的分类对象指针和可能的错误
func (s *CategoryService) CreateCategory(name string, slug string, description *string, parentID *int64) (*models.Category, error) {
	category := &models.Category{
		Name: name,
		Slug: slug,
	}

	if description != nil {
		category.Description = sql.NullString{String: *description, Valid: true}
	}

	if parentID != nil {
		category.ParentID = sql.NullInt64{Int64: *parentID, Valid: true}
	}

	id, err := models.CreateCategory(category)
	if err != nil {
		return nil, err
	}

	category.ID = id
	return category, nil
}

// GetCategoryByID 根据分类ID查询分类
// 参数: categoryID - 分类ID
// 返回: 分类对象指针（不存在则为nil）和可能的错误
func (s *CategoryService) GetCategoryByID(categoryID int64) (*models.Category, error) {
	return models.GetCategoryByID(categoryID)
}

// GetCategoryBySlug 根据分类URL别名查询分类
// 参数: slug - 分类URL别名
// 返回: 分类对象指针（不存在则为nil）和可能的错误
func (s *CategoryService) GetCategoryBySlug(slug string) (*models.Category, error) {
	return models.GetCategoryBySlug(slug)
}

// GetCategories 获取分类列表
// 参数: skip - 跳过数量, limit - 返回数量限制, parentID - 父分类ID（可选）
// 返回: 分类列表和可能的错误
func (s *CategoryService) GetCategories(skip int, limit int, parentID *int64) ([]*models.Category, error) {
	var categories []*models.Category

	qb := database.NewQueryBuilder((&models.Category{}).TableName())

	if parentID != nil {
		qb.Where("parent_id = ?", *parentID)
	} else {
		qb.Where("parent_id IS NULL")
	}

	qb.OrderBy("created_at DESC")
	qb.Offset(skip)
	qb.Limit(limit)

	query, args := qb.Build()
	rows, err := database.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	err = database.ScanRows(rows, &categories)
	return categories, err
}

// CountCategories 统计分类数量
// 参数: parentID - 父分类ID（可选）
// 返回: 分类数量和可能的错误
func (s *CategoryService) CountCategories(parentID *int64) (int64, error) {
	qb := database.NewQueryBuilder((&models.Category{}).TableName())

	if parentID != nil {
		qb.Where("parent_id = ?", *parentID)
	} else {
		qb.Where("parent_id IS NULL")
	}

	query, args := qb.BuildCount()
	var count int64
	err := database.QueryRow(query, args...).Scan(&count)
	return count, err
}

// UpdateCategory 更新分类信息
// 参数: categoryID - 分类ID, name - 分类名称, slug - 分类URL别名, description - 分类描述, parentID - 父分类ID（可选）
// 返回: 更新后的分类对象指针（不存在则为nil）和可能的错误
func (s *CategoryService) UpdateCategory(categoryID int64, name *string, slug *string, description *string, parentID *int64) (*models.Category, error) {
	category, err := models.GetCategoryByID(categoryID)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, nil
	}

	if name != nil {
		category.Name = *name
	}

	if slug != nil {
		category.Slug = *slug
	}

	if description != nil {
		category.Description = sql.NullString{String: *description, Valid: true}
	} else {
		category.Description = sql.NullString{Valid: false}
	}

	if parentID != nil {
		category.ParentID = sql.NullInt64{Int64: *parentID, Valid: true}
	} else {
		category.ParentID = sql.NullInt64{Valid: false}
	}

	_, err = models.UpdateCategory(category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

// DeleteCategory 删除分类
// 参数: categoryID - 分类ID
// 返回: 删除是否成功和可能的错误
func (s *CategoryService) DeleteCategory(categoryID int64) (bool, error) {
	rowsAffected, err := models.DeleteCategory(categoryID)
	if err != nil {
		return false, err
	}
	return rowsAffected > 0, nil
}

// AddContentToCategory 添加内容到分类
// 参数: contentID - 内容ID, categoryID - 分类ID
// 返回: 可能的错误
func (s *CategoryService) AddContentToCategory(contentID int64, categoryID int64) error {
	return models.AddContentCategory(contentID, categoryID)
}

// RemoveContentFromCategory 从分类中移除内容
// 参数: contentID - 内容ID, categoryID - 分类ID
// 返回: 可能的错误
func (s *CategoryService) RemoveContentFromCategory(contentID int64, categoryID int64) error {
	return models.RemoveContentCategory(contentID, categoryID)
}

// GetContentCategories 获取内容的分类列表
// 参数: contentID - 内容ID
// 返回: 分类列表和可能的错误
func (s *CategoryService) GetContentCategories(contentID int64) ([]*models.Category, error) {
	categoryIDs, err := models.GetContentCategories(contentID)
	if err != nil {
		return nil, err
	}

	var categories []*models.Category
	for _, id := range categoryIDs {
		category, err := models.GetCategoryByID(id)
		if err != nil {
			return nil, err
		}
		if category != nil {
			categories = append(categories, category)
		}
	}

	return categories, nil
}

// GetCategoryContents 获取分类下的内容ID列表
// 参数: categoryID - 分类ID
// 返回: 内容ID列表和可能的错误
func (s *CategoryService) GetCategoryContents(categoryID int64) ([]int64, error) {
	return models.GetCategoryContents(categoryID)
}

// SetContentCategories 设置内容的分类（替换所有现有分类）
// 参数: contentID - 内容ID, categoryIDs - 分类ID列表
// 返回: 可能的错误
func (s *CategoryService) SetContentCategories(contentID int64, categoryIDs []int64) error {
	oldCategoryIDs, err := models.GetContentCategories(contentID)
	if err != nil {
		return err
	}

	oldMap := make(map[int64]bool)
	for _, id := range oldCategoryIDs {
		oldMap[id] = true
	}

	newMap := make(map[int64]bool)
	for _, id := range categoryIDs {
		newMap[id] = true
	}

	for _, id := range oldCategoryIDs {
		if !newMap[id] {
			if err := models.RemoveContentCategory(contentID, id); err != nil {
				return err
			}
		}
	}

	for _, id := range categoryIDs {
		if !oldMap[id] {
			if err := models.AddContentCategory(contentID, id); err != nil {
				return err
			}
		}
	}

	return nil
}
