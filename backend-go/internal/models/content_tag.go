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

// ContentTag 内容与标签的多对多关联表
type ContentTag struct {
	ContentID int64 `json:"content_id"`
	TagID     int64 `json:"tag_id"`
}

// TableName 返回表名
func (ct *ContentTag) TableName() string {
	return "content_tag"
}

// AddContentTag 添加内容标签关联
// 参数: contentID - 内容ID, tagID - 标签ID
// 返回: 可能的错误
func AddContentTag(contentID, tagID int64) error {
	_, err := database.Create(
		(&ContentTag{}).TableName(),
		[]string{"content_id", "tag_id"},
		contentID,
		tagID,
	)
	return err
}

// RemoveContentTag 移除内容标签关联
// 参数: contentID - 内容ID, tagID - 标签ID
// 返回: 可能的错误
func RemoveContentTag(contentID, tagID int64) error {
	db := database.NewDeleteBuilder((&ContentTag{}).TableName())
	db.Where("content_id = ?", contentID)
	db.Where("tag_id = ?", tagID)
	query, args := db.Build()

	_, err := database.Exec(query, args...)
	return err
}

// GetContentTags 获取内容的标签ID列表
// 参数: contentID - 内容ID
// 返回: 标签ID列表和可能的错误
func GetContentTags(contentID int64) ([]int64, error) {
	var tagIDs []int64
	qb := database.NewQueryBuilder((&ContentTag{}).TableName()).Select("tag_id").Where("content_id = ?", contentID)
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

// GetTagContents 获取标签下的内容ID列表
// 参数: tagID - 标签ID
// 返回: 内容ID列表和可能的错误
func GetTagContents(tagID int64) ([]int64, error) {
	var contentIDs []int64
	qb := database.NewQueryBuilder((&ContentTag{}).TableName()).Select("content_id").Where("tag_id = ?", tagID)
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
