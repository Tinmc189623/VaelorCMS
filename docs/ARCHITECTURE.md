# 架构说明

## 技术栈

Django 4.2 + MySQL + Nginx/Gunicorn

## 应用划分

| 应用 | 职责 |
|------|------|
| config | 项目配置、路由、WSGI |
| users | 自定义 User 模型、认证、登录注册 |
| bbs | 论坛帖子（BbsPost） |
| snippets | 代码分享（CodeSnippet） |
| articles | 文章（Article） |
| site_app | 首页、搜索、关于、帮助、游戏、个人中心、管理面板、API、AdminLog |

## 请求流程

请求 → Nginx/Gunicorn → Django WSGI → 路由匹配 → 视图 → 模板渲染 → 响应

## 数据流

- 配置：config/config.ini → config/settings.py
- 数据库：MySQL（users、bbs_posts、code_snippets、articles、admin_logs、pages、site_settings）
- 静态：assets/ → collectstatic → staticfiles/
