# VaelorCMS 开发者文档

> 写给要二次开发、做主题、写插件、调 API 的开发者。这里汇总了所有你需要的东西。

---

## 文档导航

| 文档 | 讲啥 |
|------|------|
| [主题开发](DEV_THEME.md) | 做主题、覆盖样式、发布主题 |
| [主题开发规范](DEV_THEME_SPEC.md) | 主题目录结构、info.ini、CSS 变量、版本号规范 |
| [CSS 模块化](DEV_CSS_MODULES.md) | 模块结构、Aero 特效、定制方式 |
| [插件开发](DEV_PLUGIN.md) | 钩子、扩展点、写插件 |
| [API 文档](API.md) | REST API、健康检查、搜索、统计、限流 |
| [代码规范](CODE_STYLE.md) | Python、模板、安全约定 |
| [贡献指南](CONTRIBUTING.md) | 提交规范、PR 流程 |

---

## 技术栈

- **后端**：Python 3.9+、Django 4.2 LTS
- **数据库**：MySQL 5.7+ / MariaDB 10.2+（安装向导阶段用 SQLite）
- **前端**：原生 HTML/CSS/JS，无框架

---

## 项目结构（开发者视角）

```
VaelorCMS/
├── config/           # Django 配置
├── vaelor/           # 核心扩展框架（主题、插件）
├── themes/           # 主题目录
├── users/            # 用户、认证
├── bbs/              # 论坛
├── snippets/         # 代码分享
├── articles/         # 文章
├── site_app/         # 首页、搜索、管理、API、SEO、sitemap
├── install_app/      # 安装向导
├── templates/        # 模板
├── assets/           # 默认静态文件
├── main.py           # 唯一入口
└── docs/             # 文档
```

---

## 扩展点速查

| 扩展方式 | 入口 | 文档 |
|----------|------|------|
| 主题 | `themes/<id>/style.css` | [DEV_THEME.md](DEV_THEME.md) |
| 插件钩子 | `vaelor.plugin.add_hook` | [DEV_PLUGIN.md](DEV_PLUGIN.md) |
| API | `/api/v1/` | [API.md](API.md) |
| 站点设置 | 管理后台 → 站点设置 | [ADMIN_GUIDE.md](ADMIN_GUIDE.md) |

---

## 常用 API

```
GET /api/v1/stats/           # 统计
GET /api/v1/search/?q=关键词  # 搜索
GET /api/v1/health/          # 健康检查
GET /api/v1/upgrade/         # 升级检查
GET /sitemap.xml             # 站点地图
GET /robots.txt              # 爬虫规则
```

---

## 对开发者友好

- **主题**：纯 CSS，覆盖变量即可，无需改模板
- **插件**：钩子机制，不侵入核心
- **API**：RESTful JSON，无鉴权（公开接口）
- **配置**：数据库 + config.ini，可脚本化
- **文档**：Markdown，可本地预览

---

## 对站长友好

- **SEO**：meta、canonical、sitemap、robots 可配
- **主题**：后台切换，无需改代码
- **维护模式**：一键开关
- **注册开关**：可关开放注册
