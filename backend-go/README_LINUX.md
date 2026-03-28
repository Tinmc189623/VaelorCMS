# VaelorCMS - Linux 部署指南

## 概述

VaelorCMS 是一个完全自研的现代化内容管理系统，使用 Go 语言 100% 重写，支持 Windows、Linux 和 macOS 平台。

## 作者信息

- **作者**: Tinmc189623
- **团队**: Nexlyh
- **许可证**: GNU Affero General Public License v3
- **版本**: 1.0.0

## 系统要求

- Linux 系统（Debian、Ubuntu、CentOS、Arch 等）
- Go 1.21+ （仅编译时需要）
- SQLite 3 数据库

## 快速开始

### 方法 1: 使用预编译的二进制文件

1. 下载适合你系统架构的二进制文件：
   ```bash
   # Linux amd64 (64位)
   vaelorcms-linux-amd64
   
   # Linux arm64 (ARM 64位)
   vaelorcms-linux-arm64
   
   # Linux 386 (32位)
   vaelorcms-linux-386
   
   # Linux arm (ARM 32位)
   vaelorcms-linux-arm
   ```

2. 赋予执行权限：
   ```bash
   chmod +x vaelorcms-linux-amd64
   ```

3. 运行：
   ```bash
   ./vaelorcms-linux-amd64
   ```

### 方法 2: 从源码编译

#### 安装 Go 编译环境

```bash
# Ubuntu/Debian
sudo apt update
sudo apt install -y golang

# CentOS/RHEL
sudo yum install -y golang

# Arch Linux
sudo pacman -S go

# 或者从官方下载
# https://go.dev/dl/
```

#### 编译项目

```bash
# 进入项目目录
cd backend-go

# 使用提供的编译脚本
chmod +x build.sh
./build.sh

# 或者直接编译
go build -o vaelorcms ./cmd/server

# 编译所有平台
./build.sh --all
```

编译后的文件将保存在 `bin/` 目录下。

## 配置

### 环境变量配置

创建 `.env` 文件：

```env
# 项目配置
PROJECT_NAME=VaelorCMS
VERSION=1.0.0
DESCRIPTION=现代化的内容管理系统
DEBUG=true

# 服务器配置
SERVER_PORT=8080

# 安全配置
SECRET_KEY=your-secret-key-here-change-in-production
ALGORITHM=HS256
ACCESS_TOKEN_EXPIRE_MINUTES=30

# 数据库配置
DATABASE_URL=sqlite:///./vaelorcms.db
DB_POOL_SIZE=5
DB_MAX_OVERFLOW=10
DB_POOL_RECYCLE=3600
DB_ECHO=false

# CORS 配置
CORS_ORIGINS=*

# API 配置
API_PREFIX=/api/v1

# 文件上传配置
UPLOAD_DIR=./uploads
MAX_UPLOAD_SIZE=10485760
```

### 使用 systemd 服务（推荐）

创建 `/etc/systemd/system/vaelorcms.service`：

```ini
[Unit]
Description=VaelorCMS - 现代化内容管理系统
After=network.target

[Service]
Type=simple
User=www-data
Group=www-data
WorkingDirectory=/opt/vaelorcms
ExecStart=/opt/vaelorcms/vaelorcms
Restart=always
RestartSec=10

# 环境变量
Environment="PROJECT_NAME=VaelorCMS"
Environment="VERSION=1.0.0"
Environment="DEBUG=false"
Environment="SERVER_PORT=8080"
Environment="SECRET_KEY=your-secret-key-here"
Environment="DATABASE_URL=sqlite:///var/lib/vaelorcms/vaelorcms.db"
Environment="UPLOAD_DIR=/var/lib/vaelorcms/uploads"

# 安全限制
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/var/lib/vaelorcms

[Install]
WantedBy=multi-user.target
```

安装并启动服务：

```bash
# 创建必要的目录
sudo mkdir -p /opt/vaelorcms
sudo mkdir -p /var/lib/vaelorcms
sudo mkdir -p /var/lib/vaelorcms/uploads

# 复制二进制文件
sudo cp vaelorcms /opt/vaelorcms/

# 设置权限
sudo chown -R www-data:www-data /opt/vaelorcms
sudo chown -R www-data:www-data /var/lib/vaelorcms

# 启用并启动服务
sudo systemctl daemon-reload
sudo systemctl enable vaelorcms
sudo systemctl start vaelorcms

# 查看状态
sudo systemctl status vaelorcms

# 查看日志
sudo journalctl -u vaelorcms -f
```

## 使用 Nginx 反向代理

创建 `/etc/nginx/sites-available/vaelorcms`：

```nginx
server {
    listen 80;
    server_name your-domain.com;

    # 客户端上传文件大小限制
    client_max_body_size 50M;

    # 日志
    access_log /var/log/nginx/vaelorcms_access.log;
    error_log /var/log/nginx/vaelorcms_error.log;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # WebSocket 支持
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }

    # 静态文件缓存
    location /static/ {
        proxy_pass http://127.0.0.1:8080;
        expires 30d;
        add_header Cache-Control "public, immutable";
    }
}
```

启用站点：

```bash
# 创建符号链接
sudo ln -s /etc/nginx/sites-available/vaelorcms /etc/nginx/sites-enabled/

# 测试配置
sudo nginx -t

# 重启 Nginx
sudo systemctl restart nginx
```

## 使用 Docker

创建 `Dockerfile`：

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./
RUN go mod download

# 复制源码
COPY . .

# 编译
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o vaelorcms ./cmd/server

# 最终镜像
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# 从 builder 复制二进制文件
COPY --from=builder /app/vaelorcms .

# 创建必要的目录
RUN mkdir -p /app/data /app/uploads

# 环境变量
ENV PROJECT_NAME=VaelorCMS
ENV VERSION=1.0.0
ENV DEBUG=false
ENV SERVER_PORT=8080
ENV DATABASE_URL=sqlite:///app/data/vaelorcms.db
ENV UPLOAD_DIR=/app/uploads

# 暴露端口
EXPOSE 8080

# 运行
CMD ["./vaelorcms"]
```

构建和运行：

```bash
# 构建镜像
docker build -t vaelorcms:1.0.0 .

# 运行容器
docker run -d \
  --name vaelorcms \
  -p 8080:8080 \
  -v vaelorcms_data:/app/data \
  -v vaelorcms_uploads:/app/uploads \
  vaelorcms:1.0.0
```

## 数据库备份和恢复

### 备份数据库

```bash
# 创建备份目录
mkdir -p backups

# 备份 SQLite 数据库
cp vaelorcms.db backups/vaelorcms-$(date +%Y%m%d_%H%M%S).db

# 使用 SQLite 命令备份
sqlite3 vaelorcms.db ".backup backups/vaelorcms-$(date +%Y%m%d_%H%M%S).db"
```

### 恢复数据库

```bash
# 停止服务
sudo systemctl stop vaelorcms

# 恢复备份
cp backups/vaelorcms-YYYYMMDD_HHMMSS.db vaelorcms.db

# 启动服务
sudo systemctl start vaelorcms
```

## 安全建议

1. **修改默认密钥**: 务必在生产环境中修改 `SECRET_KEY`
2. **使用 HTTPS**: 配置 SSL/TLS 证书
3. **防火墙**: 只开放必要的端口
4. **定期备份**: 设置自动备份任务
5. **更新系统**: 保持系统和依赖包更新
6. **文件权限**: 设置正确的文件和目录权限

## 故障排除

### 查看日志

```bash
# systemd 日志
sudo journalctl -u vaelorcms -f

# 应用程序日志（如果有）
tail -f /var/log/vaelorcms.log
```

### 常见问题

**Q: 端口被占用？**
```bash
# 查找占用端口的进程
sudo lsof -i :8080
sudo netstat -tulpn | grep 8080
```

**Q: 权限错误？**
```bash
# 检查文件权限
ls -la /opt/vaelorcms
ls -la /var/lib/vaelorcms

# 修复权限
sudo chown -R www-data:www-data /opt/vaelorcms
sudo chown -R www-data:www-data /var/lib/vaelorcms
```

**Q: 数据库文件锁定？**
```bash
# 确保没有其他进程在使用数据库
sudo lsof vaelorcms.db
```

## 许可证

Copyright © 2025-2026 Nexlyh. All rights reserved.

本程序使用 GNU Affero General Public License v3 许可证发布。
