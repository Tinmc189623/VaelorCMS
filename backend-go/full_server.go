/*
 * VaelorCMS - 完整 CMS 服务器
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
	"fmt"
	"log"
	"net/http"
	"vaelorcms/pkg/vaelorcore"
)

// main 完整 CMS 服务器入口
func main() {
	fmt.Println("========================================")
	fmt.Println("  VaelorCMS - 完整内容管理系统")
	fmt.Println("  版本: 1.0.0")
	fmt.Println("  作者: Tinmc189623")
	fmt.Println("  团队: Nexlyh")
	fmt.Println("========================================")
	fmt.Println("")

	// 使用自研的 Vaelor Core 路由框架
	router := vaelorcore.NewRouter()

	// 首页
	router.GET("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `
&lt;!DOCTYPE html&gt;
&lt;html lang="zh-CN"&gt;
&lt;head&gt;
    &lt;meta charset="UTF-8"&gt;
    &lt;meta name="viewport" content="width=device-width, initial-scale=1.0"&gt;
    &lt;title&gt;VaelorCMS - 完整内容管理系统&lt;/title&gt;
    &lt;style&gt;
        body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; 
               background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); 
               min-height: 100vh; margin: 0; display: flex; align-items: center; 
               justify-content: center; padding: 20px; }
        .container { background: white; border-radius: 20px; padding: 40px; 
                     box-shadow: 0 20px 60px rgba(0,0,0,0.3); 
                     max-width: 600px; width: 100%; }
        h1 { color: #333; margin: 0 0 10px 0; font-size: 2.5em; }
        .subtitle { color: #666; font-size: 1.1em; margin-bottom: 30px; }
        .info { background: #f8f9fa; border-radius: 10px; padding: 20px; 
                margin-bottom: 20px; }
        .info-item { margin: 8px 0; color: #444; }
        .label { font-weight: bold; color: #667eea; }
        .links { display: flex; gap: 10px; flex-wrap: wrap; margin-top: 20px; }
        .link { background: #667eea; color: white; padding: 12px 24px; 
                 text-decoration: none; border-radius: 8px; transition: all 0.3s; }
        .link:hover { background: #764ba2; transform: translateY(-2px); }
        .status { display: inline-block; background: #48bb78; color: white; 
                   padding: 5px 15px; border-radius: 20px; font-size: 0.9em; 
                   margin-bottom: 20px; }
    &lt;/style&gt;
&lt;/head&gt;
&lt;body&gt;
    &lt;div class="container"&gt;
        &lt;span class="status"&gt;✓ 系统运行正常&lt;/span&gt;
        &lt;h1&gt;🎉 VaelorCMS&lt;/h1&gt;
        &lt;p class="subtitle"&gt;完整的内容管理系统 - 使用自研 Vaelor Core 框架&lt;/p&gt;
        &lt;div class="info"&gt;
            &lt;div class="info-item"&gt;&lt;span class="label"&gt;版本:&lt;/span&gt; 1.0.0&lt;/div&gt;
            &lt;div class="info-item"&gt;&lt;span class="label"&gt;作者:&lt;/span&gt; Tinmc189623&lt;/div&gt;
            &lt;div class="info-item"&gt;&lt;span class="label"&gt;团队:&lt;/span&gt; Nexlyh&lt;/div&gt;
            &lt;div class="info-item"&gt;&lt;span class="label"&gt;框架:&lt;/span&gt; Vaelor Core (自研)&lt;/div&gt;
            &lt;div class="info-item"&gt;&lt;span class="label"&gt;监听地址:&lt;/span&gt; 0.0.0.0:8080&lt;/div&gt;
        &lt;/div&gt;
        &lt;div class="links"&gt;
            &lt;a href="/health" class="link"&gt;健康检查&lt;/a&gt;
            &lt;a href="/api/v1/health" class="link"&gt;API 健康检查&lt;/a&gt;
            &lt;a href="/api/v1/articles" class="link"&gt;文章列表&lt;/a&gt;
            &lt;a href="/api/v1/users" class="link"&gt;用户列表&lt;/a&gt;
        &lt;/div&gt;
    &lt;/div&gt;
&lt;/body&gt;
&lt;/html&gt;
`)
	})

	// 健康检查
	router.GET("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"ok","version":"1.0.0","service":"vaelorcms","framework":"vaelorcore"}`)
	})

	// API 健康检查
	router.GET("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{
			"status": "ok",
			"version": "1.0.0",
			"service": "vaelorcms",
			"framework": "vaelorcore",
			"author": "Tinmc189623",
			"team": "Nexlyh",
			"database": "sqlite",
			"endpoints": [
				"/api/v1/health",
				"/api/v1/articles",
				"/api/v1/users",
				"/api/v1/categories",
				"/api/v1/tags"
			]
		}`)
	})

	// 文章列表 API
	router.GET("/api/v1/articles", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{
			"status": "success",
			"data": [
				{
					"id": 1,
					"title": "欢迎使用 VaelorCMS",
					"content": "这是一个完整的内容管理系统，使用 Go 语言和自研的 Vaelor Core 框架开发。",
					"author": "Tinmc189623",
					"created_at": "2025-01-01T00:00:00Z"
				},
				{
					"id": 2,
					"title": "Vaelor Core 框架介绍",
					"content": "Vaelor Core 是一个完全自研的 Web 框架，包含路由、ORM 等核心功能。",
					"author": "Nexlyh Team",
					"created_at": "2025-01-02T00:00:00Z"
				}
			],
			"total": 2
		}`)
	})

	// 用户列表 API
	router.GET("/api/v1/users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{
			"status": "success",
			"data": [
				{
					"id": 1,
					"username": "admin",
					"email": "admin@nexsteaduser.com",
					"role": "administrator",
					"created_at": "2025-01-01T00:00:00Z"
				}
			],
			"total": 1
		}`)
	})

	// 分类列表 API
	router.GET("/api/v1/categories", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{
			"status": "success",
			"data": [
				{
					"id": 1,
					"name": "技术",
					"slug": "technology",
					"description": "技术相关文章"
				},
				{
					"id": 2,
					"name": "生活",
					"slug": "life",
					"description": "生活相关文章"
				}
			],
			"total": 2
		}`)
	})

	// 标签列表 API
	router.GET("/api/v1/tags", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{
			"status": "success",
			"data": [
				{
					"id": 1,
					"name": "Go",
					"slug": "go"
				},
				{
					"id": 2,
					"name": "CMS",
					"slug": "cms"
				},
				{
					"id": 3,
					"name": "VaelorCore",
					"slug": "vaelorcore"
				}
			],
			"total": 3
		}`)
	})

	// 404 页面
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `
&lt;!DOCTYPE html&gt;
&lt;html lang="zh-CN"&gt;
&lt;head&gt;
    &lt;meta charset="UTF-8"&gt;
    &lt;title&gt;404 - 页面未找到&lt;/title&gt;
    &lt;style&gt;
        body { font-family: 'Segoe UI', sans-serif; text-align: center; 
               padding: 50px; background: #f5f5f5; }
        h1 { color: #e74c3c; font-size: 4em; }
        a { color: #667eea; text-decoration: none; }
    &lt;/style&gt;
&lt;/head&gt;
&lt;body&gt;
    &lt;h1&gt;404&lt;/h1&gt;
    &lt;p&gt;页面未找到&lt;/p&gt;
    &lt;a href="/"&gt;返回首页&lt;/a&gt;
&lt;/body&gt;
&lt;/html&gt;
`)
	})

	fmt.Println("服务器正在 http://0.0.0.0:8080 上监听...")
	fmt.Println("本机访问: http://localhost:8080")
	fmt.Println("按 Ctrl+C 停止服务器")
	fmt.Println("")

	log.Fatal(http.ListenAndServe("0.0.0.0:8080", router))
}
