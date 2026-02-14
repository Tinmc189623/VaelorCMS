# 路由与文件

## 路由一览

| 功能 | 路径 |
|------|------|
| 首页 | / |
| 论坛列表 | /bbs/ |
| 帖子详情 | /bbs/{id}/ |
| 代码列表 | /code/ |
| 代码详情 | /code/{id}/ |
| 文章列表 | /articles/ |
| 文章详情 | /articles/{id}/ |
| 发布文章 | /articles/new/ |
| 编辑文章 | /articles/{id}/edit/ |
| 登录 | /users/login/ |
| 注册 | /users/register/ |
| 退出 | /users/logout/ |
| 用户中心 | /profile/ |
| 搜索 | /search/?q=关键词 |
| 关于 | /about/ |
| 帮助 | /help/ |
| FAQ | /faq/ |
| 小游戏 | /games/ |
| 管理后台 | /admin-panel/ |
| Django Admin | /admin/ |
| API 统计 | /api/v1/stats/ |
| API 搜索 | /api/v1/search/?q=关键词 |
| API 升级检查 | /api/v1/upgrade/ |
| API 健康检查 | /api/v1/health/ |
| 自定义页面 | /p/{slug}/ |
| 安装向导 | /install/ |

## 应用结构

| 应用 | 职责 |
|------|------|
| site_app | 首页、搜索、关于、帮助、游戏、个人中心、管理面板、API |
| users | 用户模型、认证、登录注册 |
| bbs | 论坛帖子 |
| snippets | 代码分享 |
| articles | 文章 |

## URL 配置

- `config/urls.py`：主路由
- `site_app/urls.py`：站点路由
- `site_app/api_urls.py`：API 路由
- `users/urls.py`：用户路由
- `bbs/urls.py`：论坛路由
- `snippets/urls.py`：代码路由
- `articles/urls.py`：文章路由
