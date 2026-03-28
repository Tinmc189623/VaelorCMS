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

// User 用户模型，存储系统用户信息
type User struct {
	ID             int64     `json:"id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"-"`
	IsActive       bool      `json:"is_active"`
	IsAdmin        bool      `json:"is_admin"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// TableName 返回表名
func (u *User) TableName() string {
	return "users"
}

// GetUserByID 根据ID获取用户
// 参数: id - 用户ID
// 返回: 用户对象指针和可能的错误
func GetUserByID(id int64) (*User, error) {
	var user User
	found, err := database.FindByID(user.TableName(), id, &user)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return &user, nil
}

// GetUserByUsername 根据用户名获取用户
// 参数: username - 用户名
// 返回: 用户对象指针和可能的错误
func GetUserByUsername(username string) (*User, error) {
	var users []*User
	qb := database.NewQueryBuilder((&User{}).TableName()).Where("username = ?", username)
	query, args := qb.Build()

	rows, err := database.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	err = database.ScanRows(rows, &users)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, nil
	}
	return users[0], nil
}

// GetUserByEmail 根据邮箱获取用户
// 参数: email - 邮箱地址
// 返回: 用户对象指针和可能的错误
func GetUserByEmail(email string) (*User, error) {
	var users []*User
	qb := database.NewQueryBuilder((&User{}).TableName()).Where("email = ?", email)
	query, args := qb.Build()

	rows, err := database.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	err = database.ScanRows(rows, &users)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, nil
	}
	return users[0], nil
}

// CreateUser 创建新用户
// 参数: user - 用户对象
// 返回: 插入的ID和可能的错误
func CreateUser(user *User) (int64, error) {
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	return database.Create(
		user.TableName(),
		[]string{"username", "email", "hashed_password", "is_active", "is_admin", "created_at", "updated_at"},
		user.Username,
		user.Email,
		user.HashedPassword,
		user.IsActive,
		user.IsAdmin,
		user.CreatedAt,
		user.UpdatedAt,
	)
}

// UpdateUser 更新用户信息
// 参数: user - 用户对象
// 返回: 影响的行数和可能的错误
func UpdateUser(user *User) (int64, error) {
	user.UpdatedAt = time.Now()

	updates := map[string]interface{}{
		"username":        user.Username,
		"email":           user.Email,
		"hashed_password": user.HashedPassword,
		"is_active":       user.IsActive,
		"is_admin":        user.IsAdmin,
		"updated_at":      user.UpdatedAt,
	}

	return database.Update(user.TableName(), user.ID, updates)
}

// DeleteUser 删除用户
// 参数: id - 用户ID
// 返回: 影响的行数和可能的错误
func DeleteUser(id int64) (int64, error) {
	return database.Delete((&User{}).TableName(), id)
}

// GetAllUsers 获取所有用户
// 返回: 用户列表和可能的错误
func GetAllUsers() ([]*User, error) {
	var users []*User
	err := database.FindAll((&User{}).TableName(), &users)
	return users, err
}
