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

import "vaelorcms/internal/database"

// ContentCategory 内容与分类的多对多关联表
type ContentCategory struct {
	ContentID  int64 `json:"content_id"`
	CategoryID int64 `json:"category_id"`
}

// TableName 返回表名
func (cc *ContentCategory) TableName() string {
	return "content_category"
}

// AddContentCategory 添加内容分类关联
// 参数: contentID - 内容ID, categoryID - 分类ID
// 返回: 可能的错误
func AddContentCategory(contentID, categoryID int64) error {
	_, err := database.Create(
		(&ContentCategory{}).TableName(),
		[]string{"content_id", "category_id"},
		contentID,
		categoryID,
	)
	return err
}

// RemoveContentCategory 移除内容分类关联
// 参数: contentID - 内容ID, categoryID - 分类ID
// 返回: 可能的错误
func RemoveContentCategory(contentID, categoryID int64) error {
	db := database.NewDeleteBuilder((&ContentCategory{}).TableName())
	db.Where("content_id = ?", contentID)
	db.Where("category_id = ?", categoryID)
	query, args := db.Build()

	_, err := database.Exec(query, args...)
	return err
}

// GetContentCategories 获取内容的分类ID列表
// 参数: contentID - 内容ID
// 返回: 分类ID列表和可能的错误
func GetContentCategories(contentID int64) ([]int64, error) {
	var categoryIDs []int64
	qb := database.NewQueryBuilder((&ContentCategory{}).TableName()).Select("category_id").Where("content_id = ?", contentID)
	query, args := qb.Build()

	rows, err := database.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var categoryID int64
		if err := rows.Scan(&categoryID); err != nil {
			return nil, err
		}
		categoryIDs = append(categoryIDs, categoryID)
	}

	return categoryIDs, rows.Err()
}

// GetCategoryContents 获取分类下的内容ID列表
// 参数: categoryID - 分类ID
// 返回: 内容ID列表和可能的错误
func GetCategoryContents(categoryID int64) ([]int64, error) {
	var contentIDs []int64
	qb := database.NewQueryBuilder((&ContentCategory{}).TableName()).Select("content_id").Where("category_id = ?", categoryID)
	query, args := qb.Build()

	rows, err := database.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var contentID int64
		if err := rows.Scan(&contentID); err != nil {
			return nil, err
		}
		contentIDs = append(contentIDs, contentID)
	}

	return contentIDs, rows.Err()
}
