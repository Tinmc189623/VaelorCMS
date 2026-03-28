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

// Setting 设置模型，存储系统配置信息
type Setting struct {
	ID          int64          `json:"id"`
	Key         string         `json:"key"`
	Value       sql.NullString `json:"value"`
	Description sql.NullString `json:"description"`
	Group       sql.NullString `json:"group"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// TableName 返回表名
func (s *Setting) TableName() string {
	return "settings"
}

// GetSettingByID 根据ID获取设置
// 参数: id - 设置ID
// 返回: 设置对象指针和可能的错误
func GetSettingByID(id int64) (*Setting, error) {
	var setting Setting
	found, err := database.FindByID(setting.TableName(), id, &setting)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return &setting, nil
}

// GetSettingByKey 根据键名获取设置
// 参数: key - 设置键名
// 返回: 设置对象指针和可能的错误
func GetSettingByKey(key string) (*Setting, error) {
	var settings []*Setting
	qb := database.NewQueryBuilder((&Setting{}).TableName()).Where("key = ?", key)
	query, args := qb.Build()

	rows, err := database.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	err = database.ScanRows(rows, &settings)
	if err != nil {
		return nil, err
	}

	if len(settings) == 0 {
		return nil, nil
	}
	return settings[0], nil
}

// CreateSetting 创建新设置
// 参数: setting - 设置对象
// 返回: 插入的ID和可能的错误
func CreateSetting(setting *Setting) (int64, error) {
	now := time.Now()
	setting.CreatedAt = now
	setting.UpdatedAt = now

	return database.Create(
		setting.TableName(),
		[]string{"key", "value", "description", "group", "created_at", "updated_at"},
		setting.Key,
		setting.Value,
		setting.Description,
		setting.Group,
		setting.CreatedAt,
		setting.UpdatedAt,
	)
}

// UpdateSetting 更新设置信息
// 参数: setting - 设置对象
// 返回: 影响的行数和可能的错误
func UpdateSetting(setting *Setting) (int64, error) {
	setting.UpdatedAt = time.Now()

	updates := map[string]interface{}{
		"key":         setting.Key,
		"value":       setting.Value,
		"description": setting.Description,
		"group":       setting.Group,
		"updated_at":  setting.UpdatedAt,
	}

	return database.Update(setting.TableName(), setting.ID, updates)
}

// DeleteSetting 删除设置
// 参数: id - 设置ID
// 返回: 影响的行数和可能的错误
func DeleteSetting(id int64) (int64, error) {
	return database.Delete((&Setting{}).TableName(), id)
}

// GetAllSettings 获取所有设置
// 返回: 设置列表和可能的错误
func GetAllSettings() ([]*Setting, error) {
	var settings []*Setting
	err := database.FindAll((&Setting{}).TableName(), &settings)
	return settings, err
}

// GetSettingsByGroup 根据分组获取设置
// 参数: group - 设置分组
// 返回: 设置列表和可能的错误
func GetSettingsByGroup(group string) ([]*Setting, error) {
	var settings []*Setting
	qb := database.NewQueryBuilder((&Setting{}).TableName()).Where("group = ?", group)
	query, args := qb.Build()

	rows, err := database.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	err = database.ScanRows(rows, &settings)
	return settings, err
}
