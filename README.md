# VaelorCMS

现代开源 CMS，强调**用户高度自定义**与**用户自由**。论坛、代码分享、文章、自定义页面、用户中心、管理后台。基于 **Django 4** + MySQL。

**当前版本：** Demo-26.02.13.26  
**许可证：** [GNU AGPL 3.0](LICENSE)

## 技术栈

- Python 3.9+
- Django 4.2 LTS
- MySQL 5.7+
- Nginx / Gunicorn（生产）

## 功能

- 用户注册与登录
- 论坛（BBS）、代码分享、文章发布
- 文章评论、浏览量、标签云、分类筛选
- 友情链接管理
- 自定义页面、RSS 订阅
- 主题切换、Aero 玻璃特效、自定义 CSS
- 搜索、API、小游戏
- 管理后台（Django Admin + 自定义面板）
- Web 安装向导

## 快速开始

**运行时必须执行 main.py，它是唯一主程序入口。**

```bash
python -m venv venv
venv\Scripts\activate   # Windows
source venv/bin/activate # Mac/Linux

pip install -r requirements.txt
python main.py migrate   # 初始化安装用数据库
python main.py          # 启动（未安装将自动进入安装向导）
```

## 项目结构

```
├── config/          # Django 配置
├── users/           # 用户、认证
├── bbs/             # 论坛
├── snippets/        # 代码分享
├── articles/        # 文章
├── site_app/        # 首页、搜索、管理、API
├── install_app/     # 安装向导
├── templates/       # 模板
├── assets/          # 静态文件（CSS/JS/图片）
├── config/config.ini # 业务配置
├── main.py          # 唯一主程序入口
├── manage.py        # Django 管理命令（由 main.py 委托）
└── requirements.txt
```

## 文档

- [docs/GETTING_STARTED.md](docs/GETTING_STARTED.md) - **快速上手**（5 分钟入门）
- [docs/INSTALL.md](docs/INSTALL.md) - **安装教程**
- [docs/USER_GUIDE.md](docs/USER_GUIDE.md) - **用户指南**（注册、发帖、设置，口语化）
- [docs/ADMIN_GUIDE.md](docs/ADMIN_GUIDE.md) - **管理员指南**（站点设置、安全、维护）
- [docs/SETUP-WINDOWS.md](docs/SETUP-WINDOWS.md) - Windows 本地安装（含 MariaDB/MySQL 8.x）
- [docs/API.md](docs/API.md) - API 规范
- [docs/ROUTES.md](docs/ROUTES.md) - 路由
- [docs/DEPLOYMENT.md](docs/DEPLOYMENT.md) - 部署
- [docs/DEPLOY-CLOUD.md](docs/DEPLOY-CLOUD.md) - 云平台部署
- [docs/CONFIG.md](docs/CONFIG.md) - 配置
- [docs/DEV_README.md](docs/DEV_README.md) - **开发者文档**（主题、插件、API）
- [docs/DEV_THEME_SPEC.md](docs/DEV_THEME_SPEC.md) - 主题开发规范
- [docs/DEV_CSS_MODULES.md](docs/DEV_CSS_MODULES.md) - CSS 模块化架构
- [docs/CONTRIBUTING.md](docs/CONTRIBUTING.md) - 贡献指南
- [docs/CODE_STYLE.md](docs/CODE_STYLE.md) - 代码规范
- [docs/SECURITY.md](docs/SECURITY.md) - 安全策略
- [docs/FAQ.md](docs/FAQ.md) - 常见问题

## 配置

- `config/config.ini`：数据库、站点、安全（从 config.ini.sample 复制）
- `config/settings.py`：Django 主配置，自动读取 config.ini

## 开源发布包

构建不含敏感信息的开源发布包：

```bash
python scripts/build_opensource.py
```

生成 `open_source/` 目录，包含：源代码、模板、静态资源、文档、LICENSE、config.ini.sample 等。已排除：config.ini、installed.lock、*.db、venv、用户上传等。
