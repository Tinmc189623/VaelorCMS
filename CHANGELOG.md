# 更新日志

> 越往上版本越新，越往下版本越旧

## Demo-26.02.14.05

### 开源发布

- **open_source 文件夹**：新增 `scripts/build_opensource.py`，运行后生成 `open_source/` 目录
- 包含：源代码、模板、assets、themes、docs、LICENSE、README、requirements.txt、config.ini.sample
- 排除：config.ini、installed.lock、*.db、venv、storage/uploads、__pycache__ 等敏感或生成内容

## Demo-26.02.14.04

### 配置与急救

- **全配置可编辑**：站点设置新增「高级」标签，支持编辑全部 SiteSetting 键值对及新增自定义配置
- **SEO 标题后缀**：SEO 标签页新增「标题后缀」配置项
- **一键恢复**：各配置分类页提供「一键恢复当前分类」按钮，恢复该分类为默认值
- **站点急救**：管理后台新增「站点急救」页面
  - 清除缓存
  - 清除登录锁定
  - 关闭维护模式
  - 恢复全部配置为默认
  - 修复 config.ini（从 sample 恢复）
  - 配置备份导出 / 从 JSON 文件恢复

## Demo-26.02.14.03

### 安全与架构

- **漏洞修复**：HTML 净化增强（svg、template、formaction、style 内 expression）；开放重定向防护（Referer 校验）；自定义 CSS 过滤 behavior:url()
- **安全模块 Thalix**：管理后台新增「安全审计」，检查 SECRET_KEY、DEBUG、ALLOWED_HOSTS、CSRF、Session、HTTPS 等
- **安全头**：Content-Security-Policy、Permissions-Policy 扩展
- **Vaelor Core 内核**：新增 `vaelor_core` 包，提供配置抽象、钩子、安全工具（safe_str、validate_input、validate_slug、validate_email）
- **外部库**：bleach（HTML 净化）、defusedxml（安全 XML 解析）
- **文档**：`docs/THIRD_PARTY.md` 第三方与自研组件声明，`docs/SECURITY.md` 安全说明

## Demo-26.02.14.02

### 新增（CMS 功能增强）

- **文章评论**：支持登录用户与游客评论，可回复，后台可审核
- **文章浏览量**：自动统计阅读次数
- **标签云**：文章列表按标签筛选，标签云展示
- **友情链接**：Link 模型，后台管理，页脚展示

### 优化

- 文章详情页分类/标签可点击筛选
- 文章评论默认开启，可在「管理员设置」- 内容 中关闭

## Demo-26.02.14.01

### 优化

- **安装逻辑**：零配置安装，用户无需修改任何配置文件
  - 流程精简为 3 步：许可协议 → 站点配置 → 执行安装
  - 默认 SQLite，安装后无需重启
  - 可选 MySQL，勾选后填写连接信息
  - 站点名称等可在后台「站点设置」中修改
- **站点设置**：SiteSetting 加入后台管理，支持 Web 面板编辑

## Demo-26.02.13.32

### 新增

- **Redis 缓存**：config.ini [cache] driver=redis 时使用 Redis，支持 REDIS_HOST/PORT/PASSWORD 环境变量

### 优化

- **页面排版**：overflow-x、word-break、pre/code 溢出处理、viewport-fit
- **SEO**：robots meta、twitter meta、JSON-LD WebSite（含 SearchAction）

## Demo-26.02.13.31

### 新增

- **SVG 图标库**：新增 30+ 图标（arrow-left、check、alert、lock、settings、edit、trash、calendar、api 等），替代 emoji 与文本符号
- **模板更新**：返回链接、导航、首页卡片使用图标；`docs/ICONS.md` 图标说明

## Demo-26.02.13.30

### 修复

- **站点 LOGO 与收藏图标**：重新设计 app-logo.svg，新增 favicon.ico，base.html 增加 favicon 与 apple-touch-icon
- **404/500 页面**：使用 `{% static %}` 正确加载 favicon 与样式

## Demo-26.02.13.29

### 优化

- **默认生产配置**：DEBUG 默认关闭（0），本地开发可设 `DJANGO_DEBUG=1`

## Demo-26.02.13.28

### 修复

- **CSRF 403**：反向代理下自动补充 CSRF_TRUSTED_ORIGINS，云平台部署无需手动配置
- **CSS 无法显示**：集成 WhiteNoise，生产环境由应用直接提供静态文件
- **Secure Cookie**：生产环境（DEBUG=False）自动启用，适配 HTTPS 反向代理

### 优化

- Dockerfile 构建时执行 collectstatic

## Demo-26.02.13.27

### 优化

- **数据库连接复用**：MySQL 配置 `CONN_MAX_AGE=60`，减少建连开销
- **站点设置缓存**：site_settings 上下文处理器缓存 60 秒，保存设置/页面时自动失效
- **缓存失效**：管理员保存设置、创建/编辑/删除 Page 时清除站点设置缓存

## Demo-26.02.13.26

### 优化

- **项目定位**：强调现代 CMS、用户高度自定义与用户自由，移除老旧浏览器兼容代码
- **README / About**：更新项目描述与功能列表

### 新增

- **站点全局 Aero 特效**：body.aero 时启用增强玻璃效果，可配置开关与模糊强度
- **定制化支持**：管理员可设置 Aero 开关、模糊强度、强调色覆盖、自定义 CSS 片段
- **模块化 CSS 架构**：`assets/css/` 拆分为 _variables、_base、_aero、_layout、_components、_admin 模块

### 安全补丁

- **Session/CSRF Cookie**：HttpOnly、SameSite=Lax；HTTPS 下 Secure（`DJANGO_HTTPS=1`）
- **强制 HTTPS**：开启后 HTTP→HTTPS 301 重定向，响应添加 HSTS 头
- **反向代理**：`SECURE_PROXY_SSL_HEADER` 支持 `X-Forwarded-Proto`
- **自定义页面 XSS 防护**：HTML 内容经 `html_sanitizer` 净化（移除 script、iframe、on* 事件等）
- **SECRET_KEY 校验**：生产环境弱密钥时发出警告

### 新增

- **自定义页面**：Page 模型，`/p/<slug>/` 访问，Django 后台管理，可勾选「在导航栏显示」
- **FAQ 页面**：`/faq/` 常见问题，分类展示账号、发帖、个人设置、安全、管理员等
- **FAQ 文档**：`docs/FAQ.md` 完整 FAQ 文档
- **文章分类筛选**：文章列表支持 `?category=` 按分类筛选，分类导航
- **API 文章列表**：`GET /api/v1/articles/?category=` 返回最近 50 篇已发布文章
- **论坛审核**：开启「论坛需审核」后新帖待审，管理员在帖子管理中点击「通过」
- **游客发帖**：开启「游客可发帖」后未登录用户可发帖（显示匿名）
- **自定义页面 HTML**：Page 模型 `content_is_html` 字段，勾选后内容按 HTML 渲染
- **404/500 页面**：自定义错误页，`templates/404.html`、`templates/500.html`

### 优化

- **main.py 启动修复**：`execute_from_command_line` 传入正确脚本路径，解决闪退
- **界面排版**：防止横向溢出、列表项换行、管理表格列数修正、表单溢出处理
- **首页**：增加文章、FAQ、发布文章（登录）、用户中心（登录）入口
- **搜索**：使用 list-wrap 样式，优化空结果提示
- **文章发布**：新增「发布文章」快捷按钮，一键保存为已发布
- **Sitemap**：BbsPost 仅包含已审核帖子，新增 PageSitemap

### 文档

- ROUTES.md 补充自定义页面、FAQ 路由
- API.md 补充 articles 接口
- ADMIN_GUIDE.md 补充论坛审核、游客发帖说明
- README 补充 FAQ 文档链接

## Demo-26.02.13.25

### 新增

- **主题开发规范**：`docs/DEV_THEME_SPEC.md`，规范目录结构、info.ini、CSS 变量、版本号、禁止事项
- **贡献指南**：`docs/CONTRIBUTING.md`，提交规范、代码规范、许可证说明
- **代码规范**：`docs/CODE_STYLE.md`，Python、模板、安全、文档约定
- **API 限流**：`/api/v1/` 按 IP 限流，可配置 `api_rate_limit`（次/分钟），超限返回 429，健康检查不限流
- **安全响应头**：Referrer-Policy、Permissions-Policy、X-Content-Type-Options、X-Frame-Options
- **密码特殊字符**：可选 `require_password_special`，要求密码含 !@#$%^&* 等
- **密码校验统一**：`users/password_validators.py`，注册与修改密码共用同一校验逻辑

### 安全

- 安全响应头中间件，增强 XSS、点击劫持等防护
- 密码策略增强：支持强制特殊字符
- 管理员设置新增「API 限流」「密码须含特殊字符」

### 许可证

- 采用 **GNU AGPL 3.0** 许可证，LICENSE 文件含完整版权声明与许可条款

### 文档

- SECURITY.md 补充 API 限流、响应头、密码策略说明
- API.md 补充限流说明
- ADMIN_GUIDE.md 补充新安全设置项
- README 增加文档链接

## Demo-26.02.13.24

### 新增

- **文章 RSS 订阅**：`/articles/feed/` 提供 RSS 2.0 订阅，输出最近 50 篇已发布文章，站点头部增加 RSS 链接

### 安全

- **登录失败锁定**：按 IP 记录失败次数，超过 `login_max_attempts` 后锁定 `login_lockout_minutes` 分钟，支持 X-Forwarded-For 代理场景

## Demo-26.02.13.23

### 新增

- **主题系统**：`themes/<id>/` 目录，支持多主题切换，内置深色主题
- **插件系统**：基于钩子（head_extra、nav_extra、footer_extra 等），可注入 HTML
- **SEO 优化**：meta、canonical、Open Graph、sitemap.xml、动态 robots.txt、文章 Schema.org
- **开发者文档**：DEV_THEME.md、DEV_PLUGIN.md、DEV_README.md，口语化易读
- **API 文档**：补充 sitemap、robots，对开发者友好
- **阅读体验**：文章页阅读时间估算、语义化 HTML、时间标签

### 优化

- 管理员设置新增「主题」「SEO」标签
- 文章详情页 SEO：标题、描述、关键词、结构化数据
- 模板支持 `request.page_title`、`request.page_description` 等

## Demo-26.02.13.22

- 搜索支持文章模块
- CSRF 可信来源配置（环境变量、config.ini），解决 HTTPS 部署 403
- 安装向导：端口校验、config 写入安全
- 版本号统一、管理后台与 API 补充文章统计
- 文档：云平台部署通用化，示例仅保留在文档

## Demo-26.02.13.21

### 新增

- 用户设置：资料（昵称、简介、头像、隐私、通知、偏好）、修改密码、账户与安全（登出其他设备）
- 管理员设置：站点（名称、描述、关键词、时区、语言）、安全（登录限制、密码策略、会话、HTTPS）、用户（注册开关、邮箱验证）、内容（论坛审核、游客发帖、文章评论）、维护模式
- SiteSetting 模型与 settings_service，数据库存储可编辑设置
- UserProfile 模型，扩展用户资料与偏好
- 维护模式中间件，管理员可正常访问

### 安全

- 注册与密码修改遵循管理员配置的密码策略（最小长度、强密码）
- 关闭注册时禁止新用户注册
- 登出其他设备仅清除当前用户会话

## Demo-26.02.13.20

- 技术栈迁移：PHP → Django 4，保留 config.ini、assets、sql
- 应用：users、bbs、snippets、articles、site_app、install_app
- 路由：/、/bbs/、/code/、/articles/、/users/、/search/、/admin-panel/、/api/v1/、/install/
- Web 安装向导：许可协议 → 环境检测 → 数据库配置 → 站点配置 → 执行安装
- 健康检查 `/api/v1/health/`、升级接口 `/api/v1/upgrade/`
- 统一品牌 VaelorCMS，配置项 VAELOR_SESS、vaelor_

## Demo-26.02.13.19（初始版本）

- 基于 PHP + MySQL 搭建的 CMS
- 用户、论坛、代码分享、文章、搜索、管理后台
- 路径统一使用 .php，子域名路由（www、bbs、m）
- 配置：config/config.ini，文档：ROUTES.md、API.md、DEPLOYMENT.md
