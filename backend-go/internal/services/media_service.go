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
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"vaelorcms/internal/config"
	"vaelorcms/internal/database"
	"vaelorcms/internal/models"
)

// MediaService 媒体服务类，提供媒体文件相关的业务逻辑操作
type MediaService struct {
	config *config.Config
}

// NewMediaService 创建新媒体服务实例
// 返回: 媒体服务实例指针
func NewMediaService() *MediaService {
	return &MediaService{
		config: config.GetConfig(),
	}
}

// EnsureUploadDir 确保上传目录存在
// 返回: 上传目录路径和可能的错误
func (s *MediaService) EnsureUploadDir() (string, error) {
	uploadDir := s.config.UploadDir
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err = os.MkdirAll(uploadDir, 0755)
		if err != nil {
			return "", err
		}
	}
	return uploadDir, nil
}

// GenerateUniqueFilename 生成唯一的文件名
// 参数: originalFilename - 原始文件名
// 返回: 唯一的文件名
func (s *MediaService) GenerateUniqueFilename(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(bytes) + ext
}

// UploadMedia 上传媒体文件
// 参数:
//     file - 文件内容读取器
//     originalFilename - 原始文件名
//     fileType - 文件MIME类型
//     uploadedByID - 上传者ID
// 返回: 媒体文件对象指针和可能的错误
func (s *MediaService) UploadMedia(file io.Reader, originalFilename string, fileType string, uploadedByID int64) (*models.Media, error) {
	uploadDir, err := s.EnsureUploadDir()
	if err != nil {
		return nil, err
	}

	if fileType == "" {
		ext := filepath.Ext(originalFilename)
		fileType = mime.TypeByExtension(ext)
		if fileType == "" {
			fileType = "application/octet-stream"
		}
	}

	filename := s.GenerateUniqueFilename(originalFilename)
	filePath := filepath.Join(uploadDir, filename)

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	fileSize, err := io.Copy(dst, file)
	if err != nil {
		return nil, err
	}

	media := &models.Media{
		Filename:         filename,
		OriginalFilename: originalFilename,
		FilePath:         filePath,
		FileType:         fileType,
		FileSize:         fileSize,
		UploadedByID:     uploadedByID,
	}

	id, err := models.CreateMedia(media)
	if err != nil {
		os.Remove(filePath)
		return nil, err
	}

	media.ID = id
	return media, nil
}

// GetMediaByID 根据媒体ID查询媒体
// 参数: mediaID - 媒体ID
// 返回: 媒体对象指针和可能的错误
func (s *MediaService) GetMediaByID(mediaID int64) (*models.Media, error) {
	return models.GetMediaByID(mediaID)
}

// DownloadMedia 下载媒体文件
// 参数: mediaID - 媒体ID
// 返回: 文件内容读取器、文件名、文件类型和可能的错误
func (s *MediaService) DownloadMedia(mediaID int64) (io.ReadCloser, string, string, error) {
	media, err := s.GetMediaByID(mediaID)
	if err != nil {
		return nil, "", "", err
	}
	if media == nil {
		return nil, "", "", fmt.Errorf("media not found")
	}

	file, err := os.Open(media.FilePath)
	if err != nil {
		return nil, "", "", err
	}

	return file, media.OriginalFilename, media.FileType, nil
}

// GetMediaList 获取媒体列表，支持过滤
// 参数:
//     skip - 跳过数量
//     limit - 返回数量限制
//     fileType - 文件类型（可选）
//     uploadedByID - 上传者ID（可选）
// 返回: 媒体文件列表和可能的错误
func (s *MediaService) GetMediaList(skip int, limit int, fileType string, uploadedByID *int64) ([]*models.Media, error) {
	qb := database.NewQueryBuilder((&models.Media{}).TableName())

	if fileType != "" {
		qb = qb.Where("file_type LIKE ?", "%"+fileType+"%")
	}

	if uploadedByID != nil {
		qb = qb.Where("uploaded_by_id = ?", *uploadedByID)
	}

	qb = qb.OrderBy("created_at DESC").Offset(skip).Limit(limit)

	query, args := qb.Build()

	rows, err := database.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mediaList []*models.Media
	err = database.ScanRows(rows, &mediaList)
	return mediaList, err
}

// CountMedia 统计媒体数量，支持过滤
// 参数:
//     fileType - 文件类型（可选）
//     uploadedByID - 上传者ID（可选）
// 返回: 媒体数量和可能的错误
func (s *MediaService) CountMedia(fileType string, uploadedByID *int64) (int64, error) {
	qb := database.NewQueryBuilder((&models.Media{}).TableName()).Select("COUNT(*) as count")

	if fileType != "" {
		qb = qb.Where("file_type LIKE ?", "%"+fileType+"%")
	}

	if uploadedByID != nil {
		qb = qb.Where("uploaded_by_id = ?", *uploadedByID)
	}

	query, args := qb.Build()

	row := database.QueryRow(query, args...)

	var count int64
	err := row.Scan(&count)
	return count, err
}

// DeleteMedia 删除媒体记录和文件
// 参数: mediaID - 媒体ID
// 返回: 删除是否成功和可能的错误
func (s *MediaService) DeleteMedia(mediaID int64) (bool, error) {
	media, err := s.GetMediaByID(mediaID)
	if err != nil {
		return false, err
	}
	if media == nil {
		return false, nil
	}

	if _, err := os.Stat(media.FilePath); err == nil {
		err = os.Remove(media.FilePath)
		if err != nil {
		}
	}

	rowsAffected, err := models.DeleteMedia(mediaID)
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}
