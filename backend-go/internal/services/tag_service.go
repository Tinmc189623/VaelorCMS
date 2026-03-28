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
	"vaelorcms/internal/database"
	"vaelorcms/internal/models"
)

// TagService 标签服务类，提供标签相关的业务逻辑操作
type TagService struct{}

// NewTagService 创建新的标签服务实例
// 返回: 标签服务实例指针
func NewTagService() *TagService {
	return &TagService{}
}

// CreateTag 创建新标签
// 参数: name - 标签名称, slug - 标签URL别名
// 返回: 创建的标签对象指针和可能的错误
func (s *TagService) CreateTag(name, slug string) (*models.Tag, error) {
	tag := &models.Tag{
		Name: name,
		Slug: slug,
	}

	id, err := models.CreateTag(tag)
	if err != nil {
		return nil, err
	}

	tag.ID = id
	return tag, nil
}

// GetTagByID 根据标签ID查询标签
// 参数: tagID - 标签ID
// 返回: 标签对象指针和可能的错误
func (s *TagService) GetTagByID(tagID int64) (*models.Tag, error) {
	return models.GetTagByID(tagID)
}

// GetTagBySlug 根据标签URL别名查询标签
// 参数: slug - 标签URL别名
// 返回: 标签对象指针和可能的错误
func (s *TagService) GetTagBySlug(slug string) (*models.Tag, error) {
	return models.GetTagBySlug(slug)
}

// GetTags 获取标签列表，支持搜索
// 参数: skip - 跳过数量, limit - 返回数量限制, search - 搜索关键词（可选）
// 返回: 标签列表和可能的错误
func (s *TagService) GetTags(skip, limit int, search string) ([]*models.Tag, error) {
	var tags []*models.Tag
	qb := database.NewQueryBuilder((&models.Tag{}).TableName())

	if search != "" {
		qb.Where("name LIKE ?", "%"+search+"%")
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

	err = database.ScanRows(rows, &tags)
	return tags, err
}

// CountTags 统计标签数量，支持搜索
// 参数: search - 搜索关键词（可选）
// 返回: 标签数量和可能的错误
func (s *TagService) CountTags(search string) (int64, error) {
	qb := database.NewQueryBuilder((&models.Tag{}).TableName())

	if search != "" {
		qb.Where("name LIKE ?", "%"+search+"%")
	}

	query, args := qb.BuildCount()
	var count int64
	err := database.QueryRow(query, args...).Scan(&count)
	return count, err
}

// UpdateTag 更新标签信息
// 参数: tagID - 标签ID, name - 标签名称, slug - 标签URL别名
// 返回: 更新后的标签对象指针和可能的错误
func (s *TagService) UpdateTag(tagID int64, name, slug string) (*models.Tag, error) {
	tag, err := models.GetTagByID(tagID)
	if err != nil {
		return nil, err
	}
	if tag == nil {
		return nil, nil
	}

	tag.Name = name
	tag.Slug = slug

	_, err = models.UpdateTag(tag)
	if err != nil {
		return nil, err
	}

	return tag, nil
}

// DeleteTag 删除标签
// 参数: tagID - 标签ID
// 返回: 删除是否成功和可能的错误
func (s *TagService) DeleteTag(tagID int64) (bool, error) {
	tag, err := models.GetTagByID(tagID)
	if err != nil {
		return false, err
	}
	if tag == nil {
		return false, nil
	}

	_, err = models.DeleteTag(tagID)
	if err != nil {
		return false, err
	}

	return true, nil
}

// AddContentTag 添加内容标签关联
// 参数: contentID - 内容ID, tagID - 标签ID
// 返回: 可能的错误
func (s *TagService) AddContentTag(contentID, tagID int64) error {
	return models.AddContentTag(contentID, tagID)
}

// RemoveContentTag 移除内容标签关联
// 参数: contentID - 内容ID, tagID - 标签ID
// 返回: 可能的错误
func (s *TagService) RemoveContentTag(contentID, tagID int64) error {
	return models.RemoveContentTag(contentID, tagID)
}

// GetContentTags 获取内容的标签列表
// 参数: contentID - 内容ID
// 返回: 标签列表和可能的错误
func (s *TagService) GetContentTags(contentID int64) ([]*models.Tag, error) {
	tagIDs, err := models.GetContentTags(contentID)
	if err != nil {
		return nil, err
	}

	var tags []*models.Tag
	for _, tagID := range tagIDs {
		tag, err := models.GetTagByID(tagID)
		if err != nil {
			return nil, err
		}
		if tag != nil {
			tags = append(tags, tag)
		}
	}

	return tags, nil
}

// GetTagContents 获取标签下的内容ID列表
// 参数: tagID - 标签ID
// 返回: 内容ID列表和可能的错误
func (s *TagService) GetTagContents(tagID int64) ([]int64, error) {
	return models.GetTagContents(tagID)
}

// SetContentTags 设置内容的标签（先清除所有关联再添加新的）
// 参数: contentID - 内容ID, tagIDs - 标签ID列表
// 返回: 可能的错误
func (s *TagService) SetContentTags(contentID int64, tagIDs []int64) error {
	oldTagIDs, err := models.GetContentTags(contentID)
	if err != nil {
		return err
	}

	for _, tagID := range oldTagIDs {
		err = models.RemoveContentTag(contentID, tagID)
		if err != nil {
			return err
		}
	}

	for _, tagID := range tagIDs {
		err = models.AddContentTag(contentID, tagID)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetOrCreateTag 获取或创建标签
// 参数: name - 标签名称, slug - 标签URL别名
// 返回: 标签对象指针和可能的错误
func (s *TagService) GetOrCreateTag(name, slug string) (*models.Tag, error) {
	tag, err := models.GetTagBySlug(slug)
	if err != nil {
		return nil, err
	}
	if tag != nil {
		return tag, nil
	}

	return s.CreateTag(name, slug)
}
