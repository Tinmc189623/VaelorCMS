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

// Media 媒体文件模型，存储上传的媒体文件信息
type Media struct {
	ID               int64     `json:"id"`
	Filename         string    `json:"filename"`
	OriginalFilename string    `json:"original_filename"`
	FilePath         string    `json:"file_path"`
	FileType         string    `json:"file_type"`
	FileSize         int64     `json:"file_size"`
	UploadedByID     int64     `json:"uploaded_by_id"`
	CreatedAt        time.Time `json:"created_at"`
}

// TableName 返回表名
func (m *Media) TableName() string {
	return "media"
}

// GetMediaByID 根据ID获取媒体文件
// 参数: id - 媒体文件ID
// 返回: 媒体文件对象指针和可能的错误
func GetMediaByID(id int64) (*Media, error) {
	var media Media
	found, err := database.FindByID(media.TableName(), id, &media)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return &media, nil
}

// CreateMedia 创建新媒体文件记录
// 参数: media - 媒体文件对象
// 返回: 插入的ID和可能的错误
func CreateMedia(media *Media) (int64, error) {
	media.CreatedAt = time.Now()

	return database.Create(
		media.TableName(),
		[]string{"filename", "original_filename", "file_path", "file_type", "file_size", "uploaded_by_id", "created_at"},
		media.Filename,
		media.OriginalFilename,
		media.FilePath,
		media.FileType,
		media.FileSize,
		media.UploadedByID,
		media.CreatedAt,
	)
}

// DeleteMedia 删除媒体文件记录
// 参数: id - 媒体文件ID
// 返回: 影响的行数和可能的错误
func DeleteMedia(id int64) (int64, error) {
	return database.Delete((&Media{}).TableName(), id)
}

// GetAllMedia 获取所有媒体文件
// 返回: 媒体文件列表和可能的错误
func GetAllMedia() ([]*Media, error) {
	var mediaList []*Media
	err := database.FindAll((&Media{}).TableName(), &mediaList)
	return mediaList, err
}

// GetMediaByUploaderID 根据上传者ID获取媒体文件列表
// 参数: uploadedByID - 上传者ID
// 返回: 媒体文件列表和可能的错误
func GetMediaByUploaderID(uploadedByID int64) ([]*Media, error) {
	var mediaList []*Media
	qb := database.NewQueryBuilder((&Media{}).TableName()).Where("uploaded_by_id = ?", uploadedByID)
	query, args := qb.Build()

	rows, err := database.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	err = database.ScanRows(rows, &mediaList)
	return mediaList, err
}
