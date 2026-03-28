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

package config

import (
	"crypto/rand"
	"encoding/base64"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// ServerConfig 服务器配置
type ServerConfig struct {
	Port int
}

// Config 是应用程序配置结构
type Config struct {
	// 项目基本配置
	ProjectName     string
	Version         string
	Description     string
	Debug           bool

	// 服务器配置
	Server ServerConfig

	// 安全配置
	SecretKey               string
	Algorithm               string
	AccessTokenExpireMinutes int

	// 数据库配置
	DatabaseURL   string
	DBPoolSize    int
	DBMaxOverflow int
	DBPoolRecycle int
	DBEcho        bool

	// CORS 配置
	CORSOrigins []string

	// API 配置
	APIPrefix string

	// 文件上传配置
	UploadDir     string
	MaxUploadSize int64
}

// globalConfig 是全局配置实例
var globalConfig *Config

// generateSecretKey 生成一个随机的密钥
func generateSecretKey() string {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "default-secret-key-please-change-in-production"
	}
	return base64.URLEncoding.EncodeToString(bytes)
}

// LoadConfig 加载应用程序配置
func LoadConfig() (*Config, error) {
	// 尝试加载 .env 文件
	_ = godotenv.Load()

	v := viper.New()

	// 设置默认值
	setDefaults(v)

	// 绑定环境变量
	bindEnv(v)

	// 自动读取环境变量
	v.AutomaticEnv()

	// 创建配置实例
	config := &Config{
		ProjectName:              v.GetString("PROJECT_NAME"),
		Version:                  v.GetString("VERSION"),
		Description:              v.GetString("DESCRIPTION"),
		Debug:                    v.GetBool("DEBUG"),
		Server: ServerConfig{
			Port: v.GetInt("SERVER_PORT"),
		},
		SecretKey:                v.GetString("SECRET_KEY"),
		Algorithm:                v.GetString("ALGORITHM"),
		AccessTokenExpireMinutes: v.GetInt("ACCESS_TOKEN_EXPIRE_MINUTES"),
		DatabaseURL:              v.GetString("DATABASE_URL"),
		DBPoolSize:               v.GetInt("DB_POOL_SIZE"),
		DBMaxOverflow:            v.GetInt("DB_MAX_OVERFLOW"),
		DBPoolRecycle:            v.GetInt("DB_POOL_RECYCLE"),
		DBEcho:                   v.GetBool("DB_ECHO"),
		APIPrefix:                v.GetString("API_PREFIX"),
		UploadDir:                v.GetString("UPLOAD_DIR"),
		MaxUploadSize:            v.GetInt64("MAX_UPLOAD_SIZE"),
	}

	// 解析 CORS 源
	corsOriginsStr := v.GetString("CORS_ORIGINS")
	isList := strings.HasPrefix(corsOriginsStr, "[") && strings.HasSuffix(corsOriginsStr, "]")
	if isList {
		// 已经是列表形式
		config.CORSOrigins = v.GetStringSlice("CORS_ORIGINS")
	} else {
		// 逗号分隔的字符串
		parts := strings.Split(corsOriginsStr, ",")
		for _, part := range parts {
			trimmed := strings.TrimSpace(part)
			if trimmed != "" {
				config.CORSOrigins = append(config.CORSOrigins, trimmed)
			}
		}
	}

	// 如果配置为空，使用默认值
	if len(config.CORSOrigins) == 0 {
		config.CORSOrigins = []string{"*"}
	}

	// 如果密钥为空，生成一个
	if config.SecretKey == "" {
		config.SecretKey = generateSecretKey()
	}

	globalConfig = config
	return config, nil
}

// setDefaults 设置默认配置值
func setDefaults(v *viper.Viper) {
	v.SetDefault("PROJECT_NAME", "VaelorCMS")
	v.SetDefault("VERSION", "1.0.0")
	v.SetDefault("DESCRIPTION", "现代化的内容管理系统")
	v.SetDefault("DEBUG", true)
	v.SetDefault("SERVER_PORT", 8080)
	v.SetDefault("ALGORITHM", "HS256")
	v.SetDefault("ACCESS_TOKEN_EXPIRE_MINUTES", 30)
	v.SetDefault("DATABASE_URL", "sqlite:///./vaelorcms.db")
	v.SetDefault("DB_POOL_SIZE", 5)
	v.SetDefault("DB_MAX_OVERFLOW", 10)
	v.SetDefault("DB_POOL_RECYCLE", 3600)
	v.SetDefault("DB_ECHO", false)
	v.SetDefault("CORS_ORIGINS", "*")
	v.SetDefault("API_PREFIX", "/api/v1")
	v.SetDefault("UPLOAD_DIR", "./uploads")
	v.SetDefault("MAX_UPLOAD_SIZE", 10*1024*1024) // 10MB
}

// bindEnv 绑定环境变量
func bindEnv(v *viper.Viper) {
	_ = v.BindEnv("PROJECT_NAME")
	_ = v.BindEnv("VERSION")
	_ = v.BindEnv("DESCRIPTION")
	_ = v.BindEnv("DEBUG")
	_ = v.BindEnv("SERVER_PORT")
	_ = v.BindEnv("SECRET_KEY")
	_ = v.BindEnv("ALGORITHM")
	_ = v.BindEnv("ACCESS_TOKEN_EXPIRE_MINUTES")
	_ = v.BindEnv("DATABASE_URL")
	_ = v.BindEnv("DB_POOL_SIZE")
	_ = v.BindEnv("DB_MAX_OVERFLOW")
	_ = v.BindEnv("DB_POOL_RECYCLE")
	_ = v.BindEnv("DB_ECHO")
	_ = v.BindEnv("CORS_ORIGINS")
	_ = v.BindEnv("API_PREFIX")
	_ = v.BindEnv("UPLOAD_DIR")
	_ = v.BindEnv("MAX_UPLOAD_SIZE")
}

// GetConfig 获取全局配置实例
func GetConfig() *Config {
	if globalConfig == nil {
		var err error
		globalConfig, err = LoadConfig()
		if err != nil {
			panic(err)
		}
	}
	return globalConfig
}
