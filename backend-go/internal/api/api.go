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

package api

import (
	"fmt"
	"net/http"

	"vaelorcms/internal/config"
	"vaelorcms/pkg/vaelorcore"
)

// Server API 服务器结构
type Server struct {
	router *vaelorcore.Router
	config *config.Config
}

// NewServer 创建新的 API 服务器
// 参数: cfg - 配置对象
// 返回: 初始化好的服务器实例
func NewServer(cfg *config.Config) *Server {
	return &Server{
		router: vaelorcore.NewRouter(),
		config: cfg,
	}
}

// Setup 初始化服务器，注册路由和中间件
func (s *Server) Setup() {
	s.setupMiddlewares()
	s.setupRoutes()
}

// setupMiddlewares 设置全局中间件
func (s *Server) setupMiddlewares() {
}

// setupRoutes 设置所有 API 路由
func (s *Server) setupRoutes() {
	apiPrefix := s.config.APIPrefix
	
	s.router.GET(apiPrefix+"/health", s.HealthCheck)
	
	RegisterAuthRoutes(s.router)
	RegisterArticleRoutes(s.router)
	RegisterCategoryRoutes(s.router)
	RegisterTagRoutes(s.router)
	RegisterMediaRoutes(s.router)
	RegisterPageRoutes(s.router)
	RegisterSettingRoutes(s.router)
	RegisterContentRoutes(s.router)
}

// HealthCheck 健康检查处理函数
// 参数: w - HTTP 响应写入器, r - HTTP 请求
func (s *Server) HealthCheck(w http.ResponseWriter, r *http.Request) {
	SuccessResponse(w, map[string]interface{}{
		"status":  "ok",
		"version": s.config.Version,
	}, "服务器运行正常")
}

// GetRouter 获取路由实例
// 返回: 路由器实例
func (s *Server) GetRouter() *vaelorcore.Router {
	return s.router
}

// Start 启动 HTTP 服务器
// 参数: addr - 监听地址
// 返回: 错误信息
func (s *Server) Start(addr string) error {
	fmt.Printf("服务器正在 %s 上监听...\n", addr)
	return http.ListenAndServe(addr, s.router)
}

// StartTLS 启动 HTTPS 服务器
// 参数: addr - 监听地址, certFile - 证书文件路径, keyFile - 密钥文件路径
// 返回: 错误信息
func (s *Server) StartTLS(addr, certFile, keyFile string) error {
	fmt.Printf("HTTPS 服务器正在 %s 上监听...\n", addr)
	return http.ListenAndServeTLS(addr, certFile, keyFile, s.router)
}
