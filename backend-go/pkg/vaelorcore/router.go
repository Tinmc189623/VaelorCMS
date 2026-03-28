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

package vaelorcore

import (
	"context"
	"net/http"
	"path/filepath"
	"strings"
)

// ContextKey 定义请求上下文中存储路由参数的键类型
// 用于在请求上下文中传递和获取路由参数
type ContextKey string

const (
	// ParamsKey 路由参数在上下文中的键
	ParamsKey ContextKey = "params"
)

// HandlerFunc 定义路由处理函数类型
// 所有路由处理函数都应该是此类型
type HandlerFunc func(http.ResponseWriter, *http.Request)

// MiddlewareFunc 定义中间件函数类型
// 中间件接收一个 HandlerFunc 并返回一个新的 HandlerFunc
type MiddlewareFunc func(HandlerFunc) HandlerFunc

// node 表示路由树中的一个节点
// 用于构建和匹配路由路径
type node struct {
	path      string
	children  []*node
	isParam   bool
	paramName string
	methods   map[string]HandlerFunc
}

// Router 轻量级 Web 路由框架的核心结构
// 提供路由注册、匹配、中间件和静态文件服务功能
type Router struct {
	root         *node
	middlewares  []MiddlewareFunc
	notFound     HandlerFunc
	staticRoutes map[string]string
}

// NewRouter 创建并返回一个新的 Router 实例
// 返回: 初始化好的 Router 对象
func NewRouter() *Router {
	return &Router{
		root: &node{
			children: make([]*node, 0),
			methods:  make(map[string]HandlerFunc),
		},
		middlewares:  make([]MiddlewareFunc, 0),
		notFound:     defaultNotFound,
		staticRoutes: make(map[string]string),
	}
}

// Use 添加全局中间件
// 参数: middleware - 要添加的中间件函数
func (r *Router) Use(middleware MiddlewareFunc) {
	r.middlewares = append(r.middlewares, middleware)
}

// GET 注册 GET 方法的路由
// 参数: path - 路由路径, handler - 处理函数
func (r *Router) GET(path string, handler HandlerFunc) {
	r.register("GET", path, handler)
}

// POST 注册 POST 方法的路由
// 参数: path - 路由路径, handler - 处理函数
func (r *Router) POST(path string, handler HandlerFunc) {
	r.register("POST", path, handler)
}

// PUT 注册 PUT 方法的路由
// 参数: path - 路由路径, handler - 处理函数
func (r *Router) PUT(path string, handler HandlerFunc) {
	r.register("PUT", path, handler)
}

// DELETE 注册 DELETE 方法的路由
// 参数: path - 路由路径, handler - 处理函数
func (r *Router) DELETE(path string, handler HandlerFunc) {
	r.register("DELETE", path, handler)
}

// PATCH 注册 PATCH 方法的路由
// 参数: path - 路由路径, handler - 处理函数
func (r *Router) PATCH(path string, handler HandlerFunc) {
	r.register("PATCH", path, handler)
}

// register 向路由树中注册一个路由
// 参数: method - HTTP 方法, path - 路由路径, handler - 处理函数
func (r *Router) register(method, path string, handler HandlerFunc) {
	parts := strings.Split(path, "/")
	current := r.root

	for _, part := range parts {
		if part == "" {
			continue
		}

		found := false
		var child *node
		for _, c := range current.children {
			if c.path == part || (c.isParam && strings.HasPrefix(part, ":")) {
				child = c
				found = true
				break
			}
		}

		if !found {
			child = &node{
				path:      part,
				children:  make([]*node, 0),
				methods:   make(map[string]HandlerFunc),
				isParam:   strings.HasPrefix(part, ":"),
				paramName: strings.TrimPrefix(part, ":"),
			}
			current.children = append(current.children, child)
		}
		current = child
	}

	current.methods[method] = handler
}

// Static 注册静态文件服务路由
// 参数: prefix - URL 前缀, dir - 静态文件所在目录
func (r *Router) Static(prefix, dir string) {
	r.staticRoutes[prefix] = dir
}

// ServeHTTP 实现 http.Handler 接口
// 参数: w - HTTP 响应写入器, req - HTTP 请求
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for prefix, dir := range r.staticRoutes {
		if strings.HasPrefix(req.URL.Path, prefix) {
			r.serveStatic(w, req, prefix, dir)
			return
		}
	}

	handler, params := r.match(req.Method, req.URL.Path)
	if handler == nil {
		r.notFound(w, req)
		return
	}

	ctx := context.WithValue(req.Context(), ParamsKey, params)
	req = req.WithContext(ctx)

	wrapped := handler
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		wrapped = r.middlewares[i](wrapped)
	}

	wrapped(w, req)
}

// match 在路由树中匹配请求路径
// 参数: method - HTTP 方法, path - 请求路径
// 返回: 匹配的处理函数, 路由参数 map
func (r *Router) match(method, path string) (HandlerFunc, map[string]string) {
	parts := strings.Split(path, "/")
	current := r.root
	params := make(map[string]string)

	for _, part := range parts {
		if part == "" {
			continue
		}

		found := false
		var child *node
		for _, c := range current.children {
			if c.path == part {
				child = c
				found = true
				break
			} else if c.isParam {
				child = c
				found = true
				params[c.paramName] = part
				break
			}
		}

		if !found {
			return nil, nil
		}
		current = child
	}

	handler, ok := current.methods[method]
	if !ok {
		return nil, nil
	}

	return handler, params
}

// GetParam 从请求上下文中获取路由参数
// 参数: req - HTTP 请求, key - 参数名
// 返回: 参数值，不存在则返回空字符串
func GetParam(req *http.Request, key string) string {
	params, ok := req.Context().Value(ParamsKey).(map[string]string)
	if !ok {
		return ""
	}
	return params[key]
}

// GetParams 从请求上下文中获取所有路由参数
// 参数: req - HTTP 请求
// 返回: 所有路由参数的 map
func GetParams(req *http.Request) map[string]string {
	params, ok := req.Context().Value(ParamsKey).(map[string]string)
	if !ok {
		return make(map[string]string)
	}
	return params
}

// serveStatic 提供静态文件服务
// 参数: w - HTTP 响应写入器, req - HTTP 请求, prefix - URL 前缀, dir - 静态文件目录
func (r *Router) serveStatic(w http.ResponseWriter, req *http.Request, prefix, dir string) {
	relPath := req.URL.Path[len(prefix):]
	if relPath == "" {
		relPath = "/"
	}
	filePath := filepath.Join(dir, relPath)
	http.ServeFile(w, req, filePath)
}

// defaultNotFound 默认的 404 处理函数
// 参数: w - HTTP 响应写入器, req - HTTP 请求
func defaultNotFound(w http.ResponseWriter, req *http.Request) {
	http.NotFound(w, req)
}

// NotFound 设置自定义的 404 处理函数
// 参数: handler - 自定义的 404 处理函数
func (r *Router) NotFound(handler HandlerFunc) {
	r.notFound = handler
}
