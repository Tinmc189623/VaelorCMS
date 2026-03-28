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
 * 与此程序一起。如果没有, 请见 &lt;https://www.gnu.org/licenses/&gt;.
 */

package api

import (
	"io"
	"net/http"
	"vaelorcms/internal/services"
	"vaelorcms/pkg/vaelorcore"
)

var mediaService = services.NewMediaService()

// RegisterMediaRoutes 注册媒体相关路由
// 参数: router - Vaelor Core 路由器
func RegisterMediaRoutes(router *vaelorcore.Router) {
	router.GET("/api/media", ListMediaHandler)
	router.GET("/api/media/:id", GetMediaHandler)
	router.POST("/api/media", UploadMediaHandler)
	router.GET("/api/media/:id/download", DownloadMediaHandler)
	router.DELETE("/api/media/:id", DeleteMediaHandler)
}

// ListMediaHandler 处理获取媒体列表请求
// 参数: w - HTTP 响应写入器, r - HTTP 请求
func ListMediaHandler(w http.ResponseWriter, r *http.Request) {
	skip := GetIntQueryParam(r, "skip", 0)
	limit := GetIntQueryParam(r, "limit", 20)
	fileType := GetQueryParam(r, "file_type", "")
	
	var uploadedByID *int64
	uploadedByIDStr := GetQueryParam(r, "uploaded_by_id", "")
	if uploadedByIDStr != "" {
		var id int64
		tempID, err := GetIntParam(r, "uploaded_by_id")
		if err == nil {
			id = tempID
			uploadedByID = &id
		}
	}

	mediaList, err := mediaService.GetMediaList(skip, limit, fileType, uploadedByID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "获取媒体列表失败")
		return
	}

	total, err := mediaService.CountMedia(fileType, uploadedByID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "统计媒体数量失败")
		return
	}

	WriteSuccess(w, map[string]interface{}{
		"media": mediaList,
		"total": total,
		"skip":  skip,
		"limit": limit,
	})
}

// GetMediaHandler 处理根据ID获取媒体请求
// 参数: w - HTTP 响应写入器, r - HTTP 请求
func GetMediaHandler(w http.ResponseWriter, r *http.Request) {
	mediaID, err := GetIntParam(r, "id")
	if err != nil {
		WriteError(w, http.StatusBadRequest, "无效的媒体ID")
		return
	}

	media, err := mediaService.GetMediaByID(mediaID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "获取媒体失败")
		return
	}

	if media == nil {
		WriteError(w, http.StatusNotFound, "媒体不存在")
		return
	}

	WriteSuccess(w, media)
}

// UploadMediaHandler 处理上传媒体文件请求
// 参数: w - HTTP 响应写入器, r - HTTP 请求
func UploadMediaHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "无法解析表单数据")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		WriteError(w, http.StatusBadRequest, "请选择要上传的文件")
		return
	}
	defer file.Close()

	uploadedByIDStr := r.FormValue("uploaded_by_id")
	if uploadedByIDStr == "" {
		WriteError(w, http.StatusBadRequest, "上传者ID不能为空")
		return
	}

	var uploadedByID int64
	_, err = GetIntParam(r, "uploaded_by_id")
	if err != nil {
		WriteError(w, http.StatusBadRequest, "无效的上传者ID")
		return
	}
	uploadedByID, _ = GetIntParam(r, "uploaded_by_id")

	media, err := mediaService.UploadMedia(file, header.Filename, header.Header.Get("Content-Type"), uploadedByID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "文件上传失败")
		return
	}

	WriteSuccess(w, media)
}

// DownloadMediaHandler 处理下载媒体文件请求
// 参数: w - HTTP 响应写入器, r - HTTP 请求
func DownloadMediaHandler(w http.ResponseWriter, r *http.Request) {
	mediaID, err := GetIntParam(r, "id")
	if err != nil {
		WriteError(w, http.StatusBadRequest, "无效的媒体ID")
		return
	}

	file, filename, fileType, err := mediaService.DownloadMedia(mediaID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "获取文件失败")
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", fileType)
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	io.Copy(w, file)
}

// DeleteMediaHandler 处理删除媒体请求
// 参数: w - HTTP 响应写入器, r - HTTP 请求
func DeleteMediaHandler(w http.ResponseWriter, r *http.Request) {
	mediaID, err := GetIntParam(r, "id")
	if err != nil {
		WriteError(w, http.StatusBadRequest, "无效的媒体ID")
		return
	}

	success, err := mediaService.DeleteMedia(mediaID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "删除媒体失败")
		return
	}

	if !success {
		WriteError(w, http.StatusNotFound, "媒体不存在")
		return
	}

	WriteSuccess(w, map[string]bool{"success": true})
}
