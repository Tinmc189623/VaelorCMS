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

// Tag 标签模型，存储内容标签信息
type Tag struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 返回表名
func (t *Tag) TableName() string {
	return "tags"
}

// GetTagByID 根据ID获取标签
// 参数: id - 标签ID
// 返回: 标签对象指针和可能的错误
func GetTagByID(id int64) (*Tag, error) {
	var tag Tag
	found, err := database.FindByID(tag.TableName(), id, &tag)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return &tag, nil
}

// GetTagBySlug 根据Slug获取标签
// 参数: slug - 标签URL别名
// 返回: 标签对象指针和可能的错误
func GetTagBySlug(slug string) (*Tag, error) {
	var tags []*Tag
	qb := database.NewQueryBuilder((&Tag{}).TableName()).Where("slug = ?", slug)
	query, args := qb.Build()

	rows, err := database.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	err = database.ScanRows(rows, &tags)
	if err != nil {
		return nil, err
	}

	if len(tags) == 0 {
		return nil, nil
	}
	return tags[0], nil
}

// GetTagByName 根据名称获取标签
// 参数: name - 标签名称
// 返回: 标签对象指针和可能的错误
func GetTagByName(name string) (*Tag, error) {
	var tags []*Tag
	qb := database.NewQueryBuilder((&Tag{}).TableName()).Where("name = ?", name)
	query, args := qb.Build()

	rows, err := database.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	err = database.ScanRows(rows, &tags)
	if err != nil {
		return nil, err
	}

	if len(tags) == 0 {
		return nil, nil
	}
	return tags[0], nil
}

// CreateTag 创建新标签
// 参数: tag - 标签对象
// 返回: 插入的ID和可能的错误
func CreateTag(tag *Tag) (int64, error) {
	tag.CreatedAt = time.Now()

	return database.Create(
		tag.TableName(),
		[]string{"name", "slug", "created_at"},
		tag.Name,
		tag.Slug,
		tag.CreatedAt,
	)
}

// UpdateTag 更新标签信息
// 参数: tag - 标签对象
// 返回: 影响的行数和可能的错误
func UpdateTag(tag *Tag) (int64, error) {
	updates := map[string]interface{}{
		"name": tag.Name,
		"slug": tag.Slug,
	}

	return database.Update(tag.TableName(), tag.ID, updates)
}

// DeleteTag 删除标签
// 参数: id - 标签ID
// 返回: 影响的行数和可能的错误
func DeleteTag(id int64) (int64, error) {
	return database.Delete((&Tag{}).TableName(), id)
}

// GetAllTags 获取所有标签
// 返回: 标签列表和可能的错误
func GetAllTags() ([]*Tag, error) {
	var tags []*Tag
	err := database.FindAll((&Tag{}).TableName(), &tags)
	return tags, err
}
