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
	"errors"
	"vaelorcms/internal/database"
	"vaelorcms/internal/models"
)

type SettingService struct{}

func NewSettingService() *SettingService {
	return &SettingService{}
}

type SettingCreate struct {
	Key         string
	Value       sql.NullString
	Description sql.NullString
	Group       sql.NullString
}

type SettingUpdate struct {
	Key         *string
	Value       *sql.NullString
	Description *sql.NullString
	Group       *sql.NullString
}

func (s *SettingService) CreateSetting(settingData SettingCreate) (*models.Setting, error) {
	existingSetting, err := models.GetSettingByKey(settingData.Key)
	if err != nil {
		return nil, err
	}
	if existingSetting != nil {
		return nil, errors.New("设置键名已存在")
	}

	setting := &models.Setting{
		Key:         settingData.Key,
		Value:       settingData.Value,
		Description: settingData.Description,
		Group:       settingData.Group,
	}

	id, err := models.CreateSetting(setting)
	if err != nil {
		return nil, err
	}
	setting.ID = id

	return setting, nil
}

func (s *SettingService) GetSettingByID(settingID int64) (*models.Setting, error) {
	return models.GetSettingByID(settingID)
}

func (s *SettingService) GetSettingByKey(key string) (*models.Setting, error) {
	return models.GetSettingByKey(key)
}

func (s *SettingService) GetSettings(group *string) ([]*models.Setting, error) {
	if group == nil {
		return models.GetAllSettings()
	}
	return models.GetSettingsByGroup(*group)
}

func (s *SettingService) UpdateSetting(settingID int64, settingUpdate SettingUpdate) (*models.Setting, error) {
	setting, err := models.GetSettingByID(settingID)
	if err != nil {
		return nil, err
	}
	if setting == nil {
		return nil, nil
	}

	if settingUpdate.Key != nil {
		if *settingUpdate.Key != setting.Key {
			existingSetting, err := models.GetSettingByKey(*settingUpdate.Key)
			if err != nil {
				return nil, err
			}
			if existingSetting != nil {
				return nil, errors.New("设置键名已存在")
			}
		}
		setting.Key = *settingUpdate.Key
	}
	if settingUpdate.Value != nil {
		setting.Value = *settingUpdate.Value
	}
	if settingUpdate.Description != nil {
		setting.Description = *settingUpdate.Description
	}
	if settingUpdate.Group != nil {
		setting.Group = *settingUpdate.Group
	}

	_, err = models.UpdateSetting(setting)
	if err != nil {
		return nil, err
	}

	return setting, nil
}

func (s *SettingService) SetSettingValue(key string, value string) (bool, error) {
	setting, err := models.GetSettingByKey(key)
	if err != nil {
		return false, err
	}
	if setting == nil {
		return false, errors.New("设置不存在")
	}

	setting.Value = sql.NullString{String: value, Valid: true}

	_, err = models.UpdateSetting(setting)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *SettingService) GetSettingValue(key string) (string, error) {
	setting, err := models.GetSettingByKey(key)
	if err != nil {
		return "", err
	}
	if setting == nil {
		return "", errors.New("设置不存在")
	}
	if !setting.Value.Valid {
		return "", nil
	}
	return setting.Value.String, nil
}

func (s *SettingService) DeleteSetting(settingID int64) (bool, error) {
	setting, err := models.GetSettingByID(settingID)
	if err != nil {
		return false, err
	}
	if setting == nil {
		return false, nil
	}

	rows, err := models.DeleteSetting(settingID)
	if err != nil {
		return false, err
	}
	return rows > 0, nil
}

func (s *SettingService) GetAllSettings() ([]*models.Setting, error) {
	return models.GetAllSettings()
}

func (s *SettingService) GetSettingsByGroup(group string) ([]*models.Setting, error) {
	return models.GetSettingsByGroup(group)
}

func (s *SettingService) BulkSetSettings(settings map[string]string) (bool, error) {
	tx, err := database.BeginTx()
	if err != nil {
		return false, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	for key, value := range settings {
		setting, err := models.GetSettingByKey(key)
		if err != nil {
			return false, err
		}
		if setting == nil {
			continue
		}

		setting.Value = sql.NullString{String: value, Valid: true}
		_, err = models.UpdateSetting(setting)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}
