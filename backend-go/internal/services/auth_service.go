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

type AuthService struct {
	userService *UserService
}

func NewAuthService() *AuthService {
	return &AuthService{
		userService: NewUserService(),
	}
}

type LoginRequest struct {
	UsernameOrEmail string
	Password        string
}

type LoginResponse struct {
	User  *models.User
	Token string
}

func (s *AuthService) Login(request LoginRequest) (*LoginResponse, error) {
	user, err := s.userService.AuthenticateUser(request.UsernameOrEmail, request.Password)
	if err != nil {
		return nil, err
	}

	token, err := utils.CreateAccessToken(uint(user.ID), user.Username)
	if err != nil {
		return nil, errors.New("生成令牌失败")
	}

	return &LoginResponse{
		User:  user,
		Token: token,
	}, nil
}

type RegisterRequest struct {
	Username string
	Email    string
	Password string
}

func (s *AuthService) Register(request RegisterRequest) (*models.User, error) {
	return s.userService.CreateUser(request.Username, request.Email, request.Password)
}

type ChangePasswordRequest struct {
	UserID      int64
	OldPassword string
	NewPassword string
}

func (s *AuthService) ChangePassword(request ChangePasswordRequest) (bool, error) {
	user, err := s.userService.GetUserByID(request.UserID)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, errors.New("用户不存在")
	}

	if !utils.VerifyPassword(request.OldPassword, user.HashedPassword) {
		return false, errors.New("旧密码错误")
	}

	updates := map[string]interface{}{
		"password": request.NewPassword,
	}

	_, err = s.userService.UpdateUser(request.UserID, updates)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *AuthService) ValidateToken(token string) (int64, string, error) {
	claims, err := utils.VerifyToken(token)
	if err != nil {
		return 0, "", errors.New("无效的令牌")
	}
	return int64(claims.UserID), claims.Username, nil
}
