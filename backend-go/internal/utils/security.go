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

package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"vaelorcms/internal/config"
)

// BcryptCost 是 bcrypt 哈希的成本因子
const BcryptCost = 12

// MaxPasswordLength 是密码的最大长度（bcrypt 限制 72 字节）
const MaxPasswordLength = 72

// CustomClaims 是自定义的 JWT 声明结构
type CustomClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// VerifyPassword 验证密码是否正确
// 参数:
//   plainPassword: 明文密码
//   hashedPassword: 哈希后的密码
// 返回:
//   bool: 密码是否匹配
func VerifyPassword(plainPassword, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}

// GetPasswordHash 获取密码的哈希值
// 参数:
//   password: 明文密码
// 返回:
//   string: 哈希后的密码
//   error: 错误信息
func GetPasswordHash(password string) (string, error) {
	// bcrypt 限制密码长度不超过 72 字节，自动截断
	passwordBytes := []byte(password)
	if len(passwordBytes) > MaxPasswordLength {
		passwordBytes = passwordBytes[:MaxPasswordLength]
	}

	hashed, err := bcrypt.GenerateFromPassword(passwordBytes, BcryptCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// CreateAccessToken 创建访问令牌
// 参数:
//   userID: 用户 ID
//   username: 用户名
//   expiresDuration: 过期时间（可选，不提供则使用配置中的默认值）
// 返回:
//   string: JWT 访问令牌
//   error: 错误信息
func CreateAccessToken(userID uint, username string, expiresDuration ...time.Duration) (string, error) {
	cfg := config.GetConfig()

	var expiresIn time.Duration
	if len(expiresDuration) > 0 {
		expiresIn = expiresDuration[0]
	} else {
		expiresIn = time.Duration(cfg.AccessTokenExpireMinutes) * time.Minute
	}

	now := time.Now()
	claims := CustomClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.SecretKey))
}

// VerifyToken 验证并解码 JWT 令牌
// 参数:
//   token: JWT 令牌字符串
// 返回:
//   *CustomClaims: 解码后的声明数据
//   error: 错误信息
func VerifyToken(tokenString string) (*CustomClaims, error) {
	cfg := config.GetConfig()

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(cfg.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
