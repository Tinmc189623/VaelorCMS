# VaelorCMS 云平台部署指南

适用于各类 PaaS 云平台、容器（Docker/K8s）、1Panel 等。

## 容器部署必读：数据库连接

**错误示例**：`Can't connect to server on '127.0.0.1' (115)`

在容器内，`127.0.0.1` 指向**容器自身**，无法连接宿主机或独立 MySQL 服务。

### 自动回退（推荐）

当检测到容器环境且未设置 `DB_HOST` 时，应用会**自动使用 SQLite** 进入安装向导模式，无需额外配置即可启动。

### 使用 MySQL

在应用环境变量中设置 `DB_HOST` 等，指向 MySQL 实际地址：

| 场景 | DB_HOST 应填 |
|------|--------------|
| MySQL 在另一容器（如 docker-compose） | 服务名，如 `mysql`、`db` |
| MySQL 在宿主机（Docker Desktop） | `host.docker.internal` |
| 1Panel / 平台提供的 MySQL | 平台给出的主机名或内网 IP |
| 外网 MySQL | 域名或公网 IP |

**配置方式**：设置 `DB_HOST`、`DB_PORT`、`DB_NAME`、`DB_USER`、`DB_PASSWORD`。

## 部署步骤

- **构建命令**：`pip install -r requirements.txt`（通常自动）
- **启动命令**：`python main.py --prod`
- **预启动**（若平台支持）：`python main.py migrate --noinput --skip-checks`
- **Python 版本**：3.9+
- **数据库**：不配置 MySQL 时自动使用 SQLite，可通过 Web 安装向导配置

## 环境变量（可选）

| 变量 | 说明 | 示例 |
|------|------|------|
| DB_HOST | 数据库主机 | mysql.example.com |
| DB_PORT | 端口 | 3306 |
| DB_NAME | 数据库名 | vaelor_cms |
| DB_USER | 用户名 | vaelor_user |
| DB_PASSWORD | 密码 | *** |
| CSRF_TRUSTED_ORIGINS | CSRF 可信来源（逗号分隔，反向代理下可自动补充） | https://app.example.com |
| CSRF_TRUSTED_ORIGINS_AUTO | 反向代理下自动补充（1=启用） | 1 |
| DJANGO_SECRET_KEY | 密钥 | 随机字符串 |
| DJANGO_DEBUG | 调试（默认 0 生产，本地开发设 1） | 0 |
| PORT | 监听端口 | 8000（部分平台自动注入） |
| REDIS_HOST | Redis 主机（cache driver=redis 时） | 127.0.0.1 |
| REDIS_PORT | Redis 端口 | 6379 |
| REDIS_PASSWORD | Redis 密码 | 可选 |

## 启动命令

- **默认**：`python main.py --prod`
- **或**：`gunicorn config.wsgi:application --bind 0.0.0.0:$PORT`

## 数据库

- 无 config.ini 且未设置 DB_* 时，使用 SQLite（安装向导模式）
- 配置 MySQL 后，通过 Web 安装向导或环境变量完成设置

## 端口

应用监听 `0.0.0.0:8000`，部分平台需设置 `PORT` 环境变量并调整启动命令。
