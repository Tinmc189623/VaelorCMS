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
	"net/http"
	"strings"
	"vaelorcms/internal/services"
	"vaelorcms/pkg/vaelorcore"
)

var authService = services.NewAuthService()

// RegisterAuthRoutes 注册认证相关路由
// 参数: router - Vaelor Core 路由器
func RegisterAuthRoutes(router *vaelorcore.Router) {
	router.POST("/api/auth/login", LoginHandler)
	router.POST("/api/auth/register", RegisterHandler)
	router.POST("/api/auth/change-password", ChangePasswordHandler)
	router.POST("/api/auth/validate", ValidateTokenHandler)
}

// LoginHandler 处理用户登录请求
// 参数: w - HTTP 响应写入器, r - HTTP 请求
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UsernameOrEmail string `json:"username_or_email"`
		Password        string `json:"password"`
	}

	if err := ParseJSONRequest(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "无效的请求格式")
		return
	}

	if req.UsernameOrEmail == "" || req.Password == "" {
		WriteError(w, http.StatusBadRequest, "用户名/邮箱和密码不能为空")
		return
	}

	result, err := authService.Login(services.LoginRequest{
		UsernameOrEmail: req.UsernameOrEmail,
		Password:        req.Password,
	})

	if err != nil {
		WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	WriteSuccess(w, map[string]interface{}{
		"user":  result.User,
		"token": result.Token,
	})
}

// RegisterHandler 处理用户注册请求
// 参数: w - HTTP 响应写入器, r - HTTP 请求
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ParseJSONRequest(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "无效的请求格式")
		return
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		WriteError(w, http.StatusBadRequest, "用户名、邮箱和密码不能为空")
		return
	}

	user, err := authService.Register(services.RegisterRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	WriteSuccess(w, user)
}

// ChangePasswordHandler 处理修改密码请求
// 参数: w - HTTP 响应写入器, r - HTTP 请求
func ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID      int64  `json:"user_id"`
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	if err := ParseJSONRequest(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "无效的请求格式")
		return
	}

	if req.UserID == 0 || req.OldPassword == "" || req.NewPassword == "" {
		WriteError(w, http.StatusBadRequest, "用户ID、旧密码和新密码不能为空")
		return
	}

	success, err := authService.ChangePassword(services.ChangePasswordRequest{
		UserID:      req.UserID,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})

	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	WriteSuccess(w, map[string]bool{"success": success})
}

// ValidateTokenHandler 处理令牌验证请求
// 参数: w - HTTP 响应写入器, r - HTTP 请求
func ValidateTokenHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		WriteError(w, http.StatusUnauthorized, "缺少 Authorization 头部")
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		WriteError(w, http.StatusUnauthorized, "无效的 Authorization 格式")
		return
	}

	userID, username, err := authService.ValidateToken(token)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	WriteSuccess(w, map[string]interface{}{
		"user_id":  userID,
		"username": username,
		"valid":    true,
	})
}
