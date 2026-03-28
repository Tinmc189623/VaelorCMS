/*
 * VaelorCMS - 现代化内容管理系统
 *
 * Copyright © 2025-2026 Nexlyh. All rights reserved.
 *
 * 作者: Tinmc189623
 * 团队: Nexlyh
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

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"vaelorcms/internal/api"
	"vaelorcms/internal/config"
	"vaelorcms/internal/database"
)

// AppServer 应用服务器结构体，封装所有运行时组件
type AppServer struct {
	config *config.Config
	server *http.Server
	api    *api.Server
}

// NewAppServer 创建新的应用服务器实例
// 参数: cfg - 配置对象
// 返回: 初始化好的应用服务器实例
func NewAppServer(cfg *config.Config) *AppServer {
	apiServer := api.NewServer(cfg)
	apiServer.Setup()

	return &amp;AppServer{
		config: cfg,
		api:    apiServer,
		server: &amp;http.Server{
			Addr:         fmt.Sprintf("0.0.0.0:%d", cfg.Server.Port),
			Handler:      apiServer.GetRouter(),
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}
}

// Start 启动应用服务器
func (a *AppServer) Start() error {
	fmt.Printf("%s v%s\n", a.config.ProjectName, a.config.Version)
	fmt.Println("==========================")
	fmt.Printf("描述: %s\n", a.config.Description)
	fmt.Printf("调试模式: %v\n", a.config.Debug)
	fmt.Printf("监听地址: %s\n", a.server.Addr)
	fmt.Println("服务器初始化成功！")
	fmt.Println("按 Ctrl+C 停止服务器")

	go func() {
		if err := a.server.ListenAndServe(); err != nil &amp;&amp; err != http.ErrServerClosed {
			log.Fatalf("启动服务器失败: %v", err)
		}
	}()

	return nil
}

// Shutdown 优雅关闭应用服务器
// 参数: ctx - 上下文，用于设置超时
// 返回: 错误信息
func (a *AppServer) Shutdown(ctx context.Context) error {
	fmt.Println("\n正在停止 HTTP 服务...")
	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("HTTP 服务关闭失败: %w", err)
	}

	fmt.Println("正在关闭数据库连接...")
	if err := database.CloseDB(); err != nil {
		return fmt.Errorf("数据库关闭失败: %w", err)
	}

	fmt.Println("服务器已成功关闭")
	return nil
}

// initApp 初始化应用程序组件
// 返回: 配置对象和可能的错误
func initApp() (*config.Config, error) {
	fmt.Println("正在加载配置...")
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("加载配置失败: %w", err)
	}

	fmt.Println("正在初始化数据库...")
	if err := database.InitDB(cfg); err != nil {
		return nil, fmt.Errorf("初始化数据库失败: %w", err)
	}

	return cfg, nil
}

// main 是程序的入口点
func main() {
	cfg, err := initApp()
	if err != nil {
		log.Fatalf("应用初始化失败: %v", err)
	}

	app := NewAppServer(cfg)
	if err := app.Start(); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	&lt;-sigChan

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.Shutdown(ctx); err != nil {
		log.Fatalf("关闭服务器失败: %v", err)
	}
}
