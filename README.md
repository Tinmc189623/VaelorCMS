# VaelorCMS - 现代化内容管理系统

## 项目简介

VaelorCMS 是一个完全自研的现代化内容管理系统，使用 Go 语言开发，包含完整的 Vaelor Core 内核。

## 版本信息

- **版本**: 1.0.0
- **作者**: Tinmc189623
- **团队**: Nexlyh

## 许可证

本项目采用 GNU Affero General Public License v3.0 (AGPL-3.0) 许可证开源。

## 技术栈

- **后端**: Go 100% 自研
- **内核**: Vaelor Core（完全自研，包含路由、ORM 等核心功能）
- **数据库**: SQLite
- **许可证**: GNU AGPL v3

## 项目结构

```
VaelorCMS/
├── backend-go/              # Go 后端完整实现
│   ├── cmd/server/          # 应用程序入口
│   ├── internal/            # 内部包
│   │   ├── api/            # API 路由和处理
│   │   ├── config/         # 配置管理
│   │   ├── database/       # 数据库相关
│   │   ├── models/         # 数据模型
│   │   ├── services/       # 业务逻辑层
│   │   └── utils/          # 工具函数
│   └── pkg/
│       └── vaelorcore/     # Vaelor Core 自研内核
│           ├── router.go    # 自研路由框架
│           ├── orm.go       # 自研 ORM
│           └── ...
└── CHANGELOG.md            # 更新日志
```

## 快速开始

### 编译运行

```bash
cd backend-go
go build -o vaelorcms.exe cmd/server/main.go
./vaelorcms.exe
```

或者直接运行：

```bash
cd backend-go
go run cmd/server/main.go
```

### 访问

- 本机访问: http://localhost:8080
- 局域网访问: http://[你的IP]:8080

## API 端点

- `/` - 首页
- `/health` - 健康检查
- `/api/v1/health` - API 健康检查
- `/api/v1/articles` - 文章列表
- `/api/v1/users` - 用户列表
- `/api/v1/categories` - 分类列表
- `/api/v1/tags` - 标签列表

## 开源声明

Copyright © 2025-2026 Nexlyh. All rights reserved.

本程序是自由软件：你可以重新分发和/或修改它在 GNU Affero 通用公共许可证的条款下发布，版本 3 或 (根据你的选择) 任何更高版本。

本程序是希望它有用，但没有任何保证；甚至没有适销性或特定用途的默示保证。见 GNU Affero 通用公共许可证获取更多细节。

你应该收到 GNU Affero 通用公共许可证的副本与此程序一起。如果没有，请见 &lt;https://www.gnu.org/licenses/&gt;.
