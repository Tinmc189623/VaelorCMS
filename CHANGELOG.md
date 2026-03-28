# 更新日志

> 越往上版本越新，越往下版本越旧

## 1.0.0

### 🚀 里程碑版本：100% Go 语言彻底重写

#### 🔨 技术架构全面升级

- **彻底移除 Python 技术栈**：
  - 完全删除 `backend/` 目录下所有 Python 代码
  - 删除根目录 `requirements.txt` 依赖文件
  - 删除 `frontend/` 目录下所有 Python 相关文件（`__pycache__`、`vaelorcms_admin`、`requirements.txt`、`rxconfig.py`）
  - 项目中不再存在任何 `.py` 文件，实现纯 Go 语言技术栈

- **后端技术栈 100% Go 语言自研**：
  - 使用 Go 1.21+ 作为开发语言
  - 完全自研核心组件，不依赖大型第三方 Web 框架
  - 保持高性能、低内存占用的 Go 语言特性
  - 单二进制文件部署，无需外部运行时依赖

#### 💎 Vaelor Core 内核（完全自研、开源）

- **Web 路由框架（router.go）**：
  - 完全自研的轻量级 HTTP 路由系统
  - 支持所有 HTTP 方法：GET、POST、PUT、DELETE、PATCH、OPTIONS、HEAD
  - 路由参数支持：`/users/:id`、`/articles/:slug`
  - 查询参数解析：自动解析 URL 查询字符串
  - 中间件系统：全局中间件和路由组中间件支持
  - 静态文件服务：内置静态文件托管功能
  - 路由树匹配算法：高效的路由查找性能
  - 实现标准库 `http.Handler` 接口，可无缝集成

- **ORM 层（完全自研）**：
  - 数据库操作抽象层，支持 SQLite
  - 完整的 CRUD 操作：Create、Read、Update、Delete
  - 查询构建器：链式调用构建复杂 SQL 查询
  - 支持 WHERE、JOIN、GROUP BY、HAVING、ORDER BY、LIMIT、OFFSET
  - 关系映射：一对一（HasOne）、一对多（HasMany）、属于（BelongsTo）、多对多（ManyToMany）
  - 关系预加载：支持单个模型和模型切片的关联数据预加载
  - 事务支持：完整的事务管理，自动提交/回滚，保存点功能
  - 连接池管理：数据库连接复用，提升性能
  - 模型扫描：自动将数据库行扫描到 Go 结构体

- **工具库（utils.go）**：
  - 字符串处理工具：截断、包含检查、移除、大小写转换、蛇形/驼峰命名转换
  - 邮箱格式验证、手机号验证、用户名验证
  - 日期时间处理：格式化、解析、相对时间显示、日期计算
  - JSON 处理：序列化、反序列化、深拷贝、有效性验证
  - 随机数生成：随机字符串、随机整数、随机浮点数
  - 切片操作：打乱、随机选择、去重、分块
  - 数值处理：最小/最大值、范围限制、范围判断
  - 指针操作、三元运算符模拟等实用工具

- **错误处理系统**：
  - 统一的错误类型定义
  - 数据库操作错误常量
  - 清晰的错误信息传递
  - 错误堆栈跟踪支持

#### 🔐 用户认证与授权系统

- **JWT Token 认证**：
  - 使用 HS256 算法进行签名
  - Access Token 过期时间可配置（默认 30 分钟）
  - Refresh Token 机制（可扩展）
  - Token 生成和验证完全自研
  - Token 中包含用户 ID、用户名等关键信息

- **密码安全**：
  - 使用 bcrypt 算法进行密码哈希
  - 可配置的密码强度和加密轮数
  - 密码验证时的安全比对（防止时序攻击）

- **认证中间件**：
  - 请求拦截和 Token 验证
  - 从请求头或查询参数提取 Token
  - 将用户信息注入请求上下文
  - 未认证请求自动返回 401 状态码

- **权限控制**：
  - 用户角色系统（管理员、普通用户等）
  - 基于角色的访问控制（RBAC）
  - 资源级别的权限检查

#### 📝 文章管理系统

- **文章 CRUD 操作**：
  - 创建文章：标题、内容、摘要、封面图、状态、发布时间
  - 读取文章：按 ID 或 Slug 获取单篇文章
  - 更新文章：修改文章内容和元数据
  - 删除文章：软删除或硬删除选项

- **文章查询与过滤**：
  - 分页查询：支持 skip/limit 分页
  - 状态过滤：已发布、草稿、已归档等
  - 分类过滤：按分类 ID 筛选文章
  - 标签过滤：按标签 ID 或标签名筛选
  - 搜索功能：标题和内容全文搜索
  - 排序：按发布时间、更新时间、浏览量等排序
  - 数量统计：按条件统计文章数量

- **文章元数据**：
  - Slug（URL 别名）：自动生成和唯一性验证
  - 摘要：自动提取或手动输入
  - 封面图：支持本地文件或外部 URL
  - 浏览量统计：自动递增浏览计数
  - 评论数统计：关联评论计数
  - SEO 字段：Meta 标题、Meta 描述、关键词

- **文章关联管理**：
  - 分类关联：一篇文章可属于多个分类
  - 标签关联：一篇文章可拥有多个标签
  - 作者关联：记录文章作者信息
  - 媒体关联：文章中引用的媒体文件

#### 📂 分类管理系统

- **分类 CRUD 操作**：
  - 创建分类：名称、Slug、描述、父分类、排序、图标
  - 读取分类：按 ID 或 Slug 获取单分类，获取分类列表
  - 更新分类：修改分类信息和层级关系
  - 删除分类：删除分类及处理子分类和内容关联

- **分类层级系统**：
  - 无限层级分类支持
  - 父分类关联
  - 分类树构建和遍历
  - 同级分类排序

- **分类内容关联**：
  - 内容与分类多对多关联
  - 批量添加内容到分类
  - 批量从分类移除内容
  - 查询分类下的所有内容

#### 🏷️ 标签管理系统

- **标签 CRUD 操作**：
  - 创建标签：名称、Slug、描述、颜色
  - 读取标签：按 ID 或 Slug 获取单标签，获取标签列表
  - 更新标签：修改标签信息
  - 删除标签：删除标签及处理内容关联

- **标签查询与搜索**：
  - 标签列表分页
  - 按名称搜索标签
  - 按使用频率排序标签
  - 标签云数据生成

- **标签内容关联**：
  - 内容与标签多对多关联
  - 批量添加标签到内容
  - 批量从内容移除标签
  - 查询标签下的所有内容
  - 统计每个标签的使用次数

#### 🖼️ 媒体文件管理系统

- **文件上传功能**：
  - 支持多种文件类型：图片、文档、视频、音频等
  - MIME 类型验证
  - 文件大小限制（可配置，默认 10MB）
  - 唯一文件名生成（防止冲突）
  - 上传目录自动创建
  - 文件元数据记录：原始文件名、文件路径、文件类型、文件大小、上传时间、上传者

- **文件下载功能**：
  - 按 ID 下载文件
  - 正确设置 Content-Type 响应头
  - 正确设置 Content-Disposition 响应头
  - 支持断点续传（可扩展）

- **文件删除功能**：
  - 删除数据库记录
  - 删除物理文件
  - 安全的删除操作（检查文件是否存在）

- **媒体查询与管理**：
  - 媒体列表分页
  - 按文件类型过滤
  - 按上传者过滤
  - 按上传时间排序
  - 媒体数量统计
  - 存储空间使用统计

#### 📄 页面管理系统

- **页面 CRUD 操作**：
  - 创建页面：标题、内容、Slug、状态、模板、排序
  - 读取页面：按 ID 或 Slug 获取单页面，获取页面列表
  - 更新页面：修改页面内容和元数据
  - 删除页面：删除页面

- **页面特性**：
  - Slug（URL 别名）：自动生成和唯一性验证
  - 页面状态：已发布、草稿、已归档
  - 自定义模板：支持不同页面使用不同模板
  - 页面排序：自定义显示顺序
  - 导航栏显示：可配置是否在导航栏显示
  - 内容 HTML 支持：可选择是否按 HTML 渲染内容

- **页面查询**：
  - 页面列表分页
  - 按状态过滤
  - 按 Slug 查询
  - 按排序字段排序

#### ⚙️ 系统设置管理

- **设置 CRUD 操作**：
  - 创建设置：键名、键值、分组、描述、类型
  - 读取设置：按 ID 或键名获取单设置，获取设置列表
  - 更新设置：修改设置值
  - 删除设置：删除设置

- **设置分组系统**：
  - 站点设置分组：站点名称、描述、关键词、时区、语言
  - 安全设置分组：登录限制、密码策略、会话配置、HTTPS
  - 用户设置分组：注册开关、邮箱验证
  - 内容设置分组：论坛审核、游客发帖、文章评论
  - 维护模式分组：维护模式开关、维护信息
  - SEO 设置分组：标题后缀、Meta 描述、robots.txt
  - 主题设置分组：主题选择、自定义 CSS、Aero 特效

- **设置值类型**：
  - 字符串类型
  - 整数类型
  - 布尔类型
  - JSON 类型（复杂配置）

- **批量设置操作**：
  - 批量获取设置（按分组）
  - 批量更新设置
  - 事务支持，确保设置更新的原子性

- **设置缓存**：
  - 设置读取缓存
  - 设置变更时自动失效缓存
  - 提升设置读取性能

#### 📦 通用内容管理系统

- **内容 CRUD 操作**：
  - 创建内容：标题、内容、Slug、类型、状态、摘要、发布时间
  - 读取内容：按 ID 或 Slug 获取单内容，获取内容列表
  - 更新内容：修改内容信息
  - 删除内容：删除内容

- **内容类型系统**：
  - 文章类型
  - 页面类型
  - 论坛帖子类型
  - 代码片段类型
  - 自定义内容类型（可扩展）

- **内容查询与过滤**：
  - 内容列表分页
  - 按类型过滤
  - 按状态过滤
  - 搜索功能（标题和内容）
  - 按发布时间排序
  - 内容数量统计

- **内容关联**：
  - 分类关联（多对多）
  - 标签关联（多对多）
  - 作者关联
  - 媒体关联

#### 🌐 HTTP API 层

- **响应封装系统**：
  - 统一的 API 响应格式
  - 成功响应：数据、消息、状态码
  - 错误响应：错误码、错误消息、错误详情
  - JSON 格式响应
  - 自动设置 Content-Type 响应头

- **请求解析系统**：
  - JSON 请求体解析
  - 查询参数解析（字符串、整数、布尔值）
  - 路径参数解析
  - 表单数据解析（可扩展）
  - 文件上传解析（multipart/form-data）

- **中间件系统**：
  - 日志中间件：记录请求日志、响应时间、状态码
  - 恢复中间件：Panic 恢复，防止服务崩溃
  - CORS 中间件：跨域请求支持，可配置允许的源
  - 认证中间件：JWT Token 验证
  - 内容类型中间件：设置默认 Content-Type
  - 安全头中间件：X-Content-Type-Options、X-Frame-Options 等

- **API 路由实现**：
  - 认证路由：登录、注册、修改密码、验证 Token
  - 文章路由：文章 CRUD、搜索、分类/标签关联
  - 分类路由：分类 CRUD、内容关联管理
  - 标签路由：标签 CRUD、内容关联管理
  - 媒体路由：文件上传、下载、删除、列表
  - 页面路由：页面 CRUD
  - 设置路由：设置 CRUD、批量更新
  - 内容路由：通用内容 CRUD、搜索、过滤
  - 健康检查路由：服务健康状态检查

#### 🗄️ 数据库层

- **数据库连接管理**：
  - SQLite 数据库支持
  - 连接池管理：可配置的最大打开连接数、最大空闲连接数
  - 连接生命周期管理：可配置的连接最大存活时间
  - 连接测试：初始化时测试数据库连接
  - 外键约束：默认启用 SQLite 外键约束

- **数据库操作封装**：
  - Exec 方法：执行不返回行的 SQL
  - Query 方法：执行返回多行的查询
  - QueryRow 方法：执行返回单行的查询
  - BeginTx 方法：开始事务
  - 连接获取和释放

- **数据库迁移**：
  - 迁移脚本系统
  - 迁移版本管理
  - 向上迁移和向下迁移支持
  - 自动创建表结构
  - 初始化默认数据

#### ⚙️ 配置管理系统

- **配置源支持**：
  - 环境变量配置
  - .env 文件配置
  - 默认值配置
  - 优先级：环境变量 > .env 文件 > 默认值

- **配置项**：
  - 项目基本配置：项目名称、版本、描述、调试模式
  - 服务器配置：监听端口
  - 安全配置：密钥、算法、Token 过期时间
  - 数据库配置：数据库 URL、连接池大小、最大溢出、连接回收、SQL 日志
  - CORS 配置：允许的源
  - API 配置：API 前缀
  - 文件上传配置：上传目录、最大文件大小

- **配置加载**：
  - 应用启动时自动加载配置
  - 配置验证
  - 密钥自动生成（如果未配置）
  - 全局配置实例访问

#### 📚 数据模型层

- **用户模型（User）**：
  - ID、用户名、邮箱、密码哈希、昵称、头像、简介
  - 是否管理员、是否激活、创建时间、更新时间
  - 密码验证方法
  - 密码哈希生成方法

- **文章模型（Article）**：
  - ID、标题、Slug、内容、摘要、封面图、状态
  - 发布时间、作者 ID、浏览量、评论数
  - 创建时间、更新时间
  - 表名方法

- **分类模型（Category）**：
  - ID、名称、Slug、描述、父分类 ID、排序、图标
  - 创建时间、更新时间
  - 表名方法

- **标签模型（Tag）**：
  - ID、名称、Slug、描述、颜色、使用次数
  - 创建时间、更新时间
  - 表名方法

- **媒体模型（Media）**：
  - ID、文件名、原始文件名、文件路径、文件类型、文件大小
  - 上传者 ID、创建时间
  - 表名方法

- **页面模型（Page）**：
  - ID、标题、Slug、内容、状态、模板、排序
  - 是否在导航栏显示、创建时间、更新时间
  - 表名方法

- **设置模型（Setting）**：
  - ID、键名、键值、分组、描述、类型
  - 创建时间、更新时间
  - 表名方法

- **内容模型（Content）**：
  - ID、标题、Slug、内容、摘要、类型、状态
  - 发布时间、作者 ID、创建时间、更新时间
  - 表名方法

- **内容分类关联模型（ContentCategory）**：
  - 内容 ID、分类 ID
  - 复合主键
  - 表名方法

- **内容标签关联模型（ContentTag）**：
  - 内容 ID、标签 ID
  - 复合主键
  - 表名方法

#### 🏗️ 服务层架构

- **用户服务（UserService）**：
  - 用户创建、查询、更新、删除
  - 用户认证（用户名/邮箱 + 密码）
  - 密码修改
  - 用户列表查询和分页
  - 用户统计

- **认证服务（AuthService）**：
  - 用户登录
  - 用户注册
  - 密码修改
  - Token 验证
  - Token 刷新（可扩展）

- **文章服务（ArticleService）**：
  - 文章创建、查询、更新、删除
  - 文章列表查询、分页、搜索、过滤
  - 文章分类关联管理
  - 文章标签关联管理
  - 文章统计
  - 浏览量递增

- **分类服务（CategoryService）**：
  - 分类创建、查询、更新、删除
  - 分类列表查询、分页
  - 分类树构建
  - 分类内容关联管理
  - 分类统计

- **标签服务（TagService）**：
  - 标签创建、查询、更新、删除
  - 标签列表查询、分页、搜索
  - 标签内容关联管理
  - 标签统计
  - 标签使用次数更新

- **媒体服务（MediaService）**：
  - 文件上传
  - 文件下载
  - 文件删除
  - 媒体列表查询、分页、过滤
  - 媒体统计
  - 上传目录确保
  - 唯一文件名生成

- **页面服务（PageService）**：
  - 页面创建、查询、更新、删除
  - 页面列表查询、分页、过滤
  - 页面统计
  - Slug 唯一性验证

- **设置服务（SettingService）**：
  - 设置创建、查询、更新、删除
  - 按分组获取设置
  - 按键名获取设置值
  - 批量设置更新（事务支持）
  - 设置缓存管理

- **内容服务（ContentService）**：
  - 内容创建、查询、更新、删除
  - 内容列表查询、分页、过滤、搜索
  - 内容统计
  - 发布时间自动处理

#### 🖥️ 多平台支持

- **Windows 平台**：
  - 完整的 Windows 支持
  - PowerShell 编译脚本（build.ps1）
  - 支持当前平台快速编译
  - 支持所有平台编译（-All 参数）
  - 支持清理输出目录（-Clean 参数）
  - 支持自定义输出目录
  - Windows 可执行文件（.exe）

- **Linux 平台**：
  - 完整的 Linux 支持
  - Bash 编译脚本（build.sh）
  - 支持当前平台快速编译
  - 支持所有平台编译（--all 参数）
  - 支持清理输出目录（--clean 参数）
  - 支持自定义输出目录
  - 详细的 Linux 部署指南（README_LINUX.md）

- **macOS 平台**：
  - 完整的 macOS 支持
  - Bash 编译脚本（build.sh）
  - Intel（amd64）和 Apple Silicon（arm64）支持
  - 与 Linux 相同的使用体验

- **跨平台兼容性**：
  - 使用 `filepath` 包处理所有路径操作
  - 正确处理路径分隔符（Windows \ vs Linux/macOS /）
  - 正确处理环境变量
  - 正确处理行尾符
  - 所有平台使用相同的代码库

#### 📖 部署指南与文档

- **Linux 部署指南（README_LINUX.md）**：
  - 快速开始：预编译二进制文件使用、源码编译
  - 系统要求：Linux 发行版、Go 版本、数据库
  - 环境变量配置：完整的配置项说明和示例
  - systemd 服务配置：生产环境推荐的服务管理方式
  - Nginx 反向代理配置：完整的 Nginx 配置示例
  - Docker 部署：Dockerfile 和 Docker Compose 配置
  - 数据库备份和恢复：完整的备份恢复流程
  - 安全建议：生产环境安全最佳实践
  - 故障排除：常见问题和解决方案

- **编译脚本**：
  - build.ps1（Windows PowerShell）
  - build.sh（Linux/macOS Bash）
  - 支持 8 个目标平台：
    - Windows amd64
    - Windows 386
    - Linux amd64
    - Linux 386
    - Linux arm64
    - Linux arm
    - macOS amd64
    - macOS arm64

#### 📁 项目结构优化

```
VaelorCMS/
├── backend-go/                    # Go 后端项目（唯一后端）
│   ├── cmd/
│   │   └── server/
│   │       └── main.go           # 应用程序入口
│   ├── internal/                  # 内部包（不对外导出）
│   │   ├── api/                   # API 层
│   │   │   ├── api.go            # API 服务器主文件
│   │   │   ├── auth.go           # 认证路由
│   │   │   ├── articles.go       # 文章路由
│   │   │   ├── categories.go     # 分类路由
│   │   │   ├── tags.go           # 标签路由
│   │   │   ├── media.go          # 媒体路由
│   │   │   ├── pages.go          # 页面路由
│   │   │   ├── settings.go       # 设置路由
│   │   │   ├── content.go        # 内容路由
│   │   │   ├── middleware.go     # 中间件
│   │   │   ├── request.go        # 请求解析
│   │   │   └── response.go       # 响应封装
│   │   ├── config/                # 配置管理
│   │   │   └── config.go
│   │   ├── database/              # 数据库层
│   │   │   ├── database.go       # 数据库连接管理
│   │   │   ├── migration.go      # 数据库迁移
│   │   │   └── orm.go            # ORM 基础
│   │   ├── models/                # 数据模型
│   │   │   ├── models.go         # 模型基础
│   │   │   ├── user.go           # 用户模型
│   │   │   ├── article.go        # 文章模型
│   │   │   ├── category.go       # 分类模型
│   │   │   ├── tag.go            # 标签模型
│   │   │   ├── media.go          # 媒体模型
│   │   │   ├── page.go           # 页面模型
│   │   │   ├── setting.go        # 设置模型
│   │   │   ├── content.go        # 内容模型
│   │   │   ├── content_category.go  # 内容分类关联
│   │   │   └── content_tag.go    # 内容标签关联
│   │   ├── services/              # 业务逻辑层
│   │   │   ├── services.go       # 服务基础
│   │   │   ├── user_service.go   # 用户服务
│   │   │   ├── auth_service.go   # 认证服务
│   │   │   ├── article_service.go  # 文章服务
│   │   │   ├── category_service.go # 分类服务
│   │   │   ├── tag_service.go    # 标签服务
│   │   │   ├── media_service.go  # 媒体服务
│   │   │   ├── page_service.go   # 页面服务
│   │   │   ├── setting_service.go # 设置服务
│   │   │   └── content_service.go # 内容服务
│   │   └── utils/                 # 工具函数
│   │       ├── utils.go          # 工具基础
│   │       ├── security.go       # 安全工具（密码、JWT）
│   │       ├── slugify.go        # Slug 生成
│   │       └── utils_test.go     # 工具测试
│   ├── pkg/                       # 公共包（可对外导出）
│   │   ├── vaelorcore/           # Vaelor Core 内核
│   │   │   ├── vaelorcore.go     # 包文档
│   │   │   ├── router.go         # Web 路由框架
│   │   │   ├── model.go          # 模型基础
│   │   │   ├── query.go          # 查询构建器
│   │   │   ├── orm.go            # ORM 核心
│   │   │   ├── relations.go      # 关系映射
│   │   │   ├── transaction.go    # 事务支持
│   │   │   ├── errors.go         # 错误定义
│   │   │   └── example.go        # 使用示例
│   │   └── pkg.go
│   ├── vaelorcms.exe             # Windows 可执行文件
│   ├── build.ps1                  # Windows 编译脚本
│   ├── build.sh                   # Linux/macOS 编译脚本
│   ├── README_LINUX.md            # Linux 部署指南
│   ├── go.mod                     # Go 模块定义
│   └── go.sum                     # 依赖版本锁定
├── frontend/                      # 前端静态文件
│   └── .web/                      # 前端构建输出
├── .trae/                         # IDE 配置和规格文档
│   └── specs/
│       └── vaelorcms-100-percent-go-migration/
│           ├── spec.md            # 产品需求文档
│           ├── tasks.md           # 实现计划
│           └── checklist.md       # 验证检查清单
├── CHANGELOG.md                   # 更新日志
├── LICENSE                        # 许可证（GNU AGPL v3）
└── .gitignore                     # Git 忽略文件
```

#### 🔒 安全特性

- **密码安全**：
  - bcrypt 密码哈希
  - 可配置的密码强度
  - 安全的密码比对

- **JWT 认证**：
  - HS256 签名算法
  - 可配置的 Token 过期时间
  - 密钥自动生成
  - Token 验证中间件

- **CORS 安全**：
  - 可配置的允许源
  - 安全的 CORS 头设置

- **安全头**：
  - X-Content-Type-Options
  - X-Frame-Options
  - Content-Security-Policy（可扩展）
  - Permissions-Policy（可扩展）

- **输入验证**：
  - 邮箱格式验证
  - 手机号验证
  - 用户名验证
  - Slug 格式验证
  - SQL 注入防护（使用参数化查询）
  - XSS 防护（可扩展）

#### 📊 性能优化

- **数据库连接池**：
  - 连接复用，减少连接开销
  - 可配置的连接池大小
  - 连接生命周期管理

- **查询优化**：
  - 使用参数化查询
  - 索引支持（可扩展）
  - 查询结果缓存（可扩展）

- **设置缓存**：
  - 设置读取缓存
  - 设置变更自动失效
  - 提升设置读取性能

- **静态文件**：
  - 高效的静态文件服务
  - 支持缓存头（可扩展）

- **Go 语言特性**：
  - 高性能的 HTTP 服务
  - 低内存占用
  - 优秀的并发处理能力
  - 快速的启动时间

#### 👨‍💻 代码质量

- **完整的注释**：
  - 所有文件都有头部注释
  - 所有公共函数都有函数级注释
  - 所有结构体都有注释
  - 所有类型都有注释

- **清晰的代码结构**：
  - 分层架构（API、Service、Model、Database）
  - 模块化设计
  - 单一职责原则
  - 高内聚低耦合

- **版本管理**：
  - Go modules 依赖管理
  - 依赖版本锁定（go.sum）
  - 语义化版本（1.0.0）

- **许可证合规**：
  - 所有文件都有正确的许可证声明
  - 使用 GNU Affero General Public License v3
  - 第三方依赖许可证兼容

#### 🏷️ 版本与作者信息

- **版本**：1.0.0（正式版）
- **作者**：Tinmc189623
- **团队**：Nexlyh
- **许可证**：GNU Affero General Public License v3
- **发布日期**：2026-03-28

---

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