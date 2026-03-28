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
	"errors"
	"vaelorcms/internal/models"
	"vaelorcms/internal/utils"
)

// UserService 用户服务类，提供用户相关的业务逻辑操作
type UserService struct{}

// NewUserService 创建用户服务实例
// 返回: 用户服务实例指针
func NewUserService() *UserService {
	return &UserService{}
}

// CreateUser 创建新用户
// 参数: username - 用户名, email - 邮箱, password - 密码
// 返回: 创建的用户对象指针和可能的错误
func (s *UserService) CreateUser(username, email, password string) (*models.User, error) {
	// 检查用户名是否已存在
	existingUser, err := models.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	existingUser, err = models.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("邮箱已存在")
	}

	// 生成密码哈希
	hashedPassword, err := utils.GetPasswordHash(password)
	if err != nil {
		return nil, err
	}

	// 创建用户对象
	user := &models.User{
		Username:       username,
		Email:          email,
		HashedPassword: hashedPassword,
		IsActive:       true,
		IsAdmin:        false,
	}

	// 保存到数据库
	id, err := models.CreateUser(user)
	if err != nil {
		return nil, err
	}

	// 设置用户ID
	user.ID = id

	return user, nil
}

// GetUserByID 根据用户ID查询用户
// 参数: userID - 用户ID
// 返回: 用户对象指针和可能的错误
func (s *UserService) GetUserByID(userID int64) (*models.User, error) {
	return models.GetUserByID(userID)
}

// GetUserByUsername 根据用户名查询用户
// 参数: username - 用户名
// 返回: 用户对象指针和可能的错误
func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	return models.GetUserByUsername(username)
}

// GetUserByEmail 根据邮箱查询用户
// 参数: email - 邮箱地址
// 返回: 用户对象指针和可能的错误
func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	return models.GetUserByEmail(email)
}

// GetUserByUsernameOrEmail 根据用户名或邮箱查询用户
// 参数: usernameOrEmail - 用户名或邮箱
// 返回: 用户对象指针和可能的错误
func (s *UserService) GetUserByUsernameOrEmail(usernameOrEmail string) (*models.User, error) {
	// 首先尝试按用户名查找
	user, err := s.GetUserByUsername(usernameOrEmail)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return user, nil
	}

	// 用户名未找到，尝试按邮箱查找
	return s.GetUserByEmail(usernameOrEmail)
}

// AuthenticateUser 验证用户密码
// 参数: usernameOrEmail - 用户名或邮箱, password - 密码
// 返回: 验证成功返回用户对象指针，失败返回 nil 和错误
func (s *UserService) AuthenticateUser(usernameOrEmail, password string) (*models.User, error) {
	// 查找用户
	user, err := s.GetUserByUsernameOrEmail(usernameOrEmail)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("用户不存在")
	}

	// 检查用户是否激活
	if !user.IsActive {
		return nil, errors.New("用户已被禁用")
	}

	// 验证密码
	if !utils.VerifyPassword(password, user.HashedPassword) {
		return nil, errors.New("密码错误")
	}

	return user, nil
}

// UpdateUser 更新用户信息
// 参数: userID - 用户ID, updates - 要更新的字段映射
// 返回: 更新后的用户对象指针和可能的错误
func (s *UserService) UpdateUser(userID int64, updates map[string]interface{}) (*models.User, error) {
	// 查找用户
	user, err := s.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("用户不存在")
	}

	// 处理密码更新
	if password, ok := updates["password"]; ok {
		hashedPassword, err := utils.GetPasswordHash(password.(string))
		if err != nil {
			return nil, err
		}
		updates["hashed_password"] = hashedPassword
		delete(updates, "password")
	}

	// 更新用户字段
	for key, value := range updates {
		switch key {
		case "username":
			user.Username = value.(string)
		case "email":
			user.Email = value.(string)
		case "hashed_password":
			user.HashedPassword = value.(string)
		case "is_active":
			user.IsActive = value.(bool)
		case "is_admin":
			user.IsAdmin = value.(bool)
		}
	}

	// 保存更新
	_, err = models.UpdateUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser 删除用户
// 参数: userID - 用户ID
// 返回: 删除是否成功和可能的错误
func (s *UserService) DeleteUser(userID int64) (bool, error) {
	// 查找用户
	user, err := s.GetUserByID(userID)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, errors.New("用户不存在")
	}

	// 删除用户
	rowsAffected, err := models.DeleteUser(userID)
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

// GetAllUsers 获取所有用户
// 返回: 用户列表和可能的错误
func (s *UserService) GetAllUsers() ([]*models.User, error) {
	return models.GetAllUsers()
}
