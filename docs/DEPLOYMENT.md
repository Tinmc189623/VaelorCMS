# 部署指南

## 环境

- Python 3.9+
- MySQL 5.7+ / MariaDB 10.2+
- Nginx 1.18+（推荐）或 Apache + mod_wsgi

## 方式一：Web 安装向导（推荐）

1. 上传代码至站点根目录  
2. 创建虚拟环境并安装依赖：
   ```bash
   python -m venv venv
   source venv/bin/activate  # 或 Windows: venv\Scripts\activate
   pip install -r requirements.txt
   ```
3. 初始化安装用数据库（使用 SQLite 临时库）：
   ```bash
   python main.py migrate
   ```
4. 启动服务并访问站点，将自动进入安装向导：
   ```bash
   python main.py        # 开发模式
   python main.py --prod # 生产模式（Gunicorn）
   ```
5. 按向导完成：许可协议 → 环境检测 → 数据库配置 → 站点配置 → 执行安装
6. 安装完成后**重启应用服务器**，使 MySQL 配置生效
7. 收集静态文件：`python main.py collectstatic`

## 方式二：手动配置

1. 上传代码并安装依赖（同上）
2. 复制 `config/config.ini.sample` 为 `config/config.ini`，填写 [database]、[site] 等  
3. 创建 `config/installed.lock` 空文件（表示已安装，跳过安装向导）
4. 数据库迁移：
   ```bash
   python main.py migrate
   python main.py createsuperuser
   ```
5. 收集静态文件：`python main.py collectstatic`

## 已有数据库

若已有旧版数据，先执行：

```sql
ALTER TABLE users ADD COLUMN last_login datetime NULL;
```

然后执行 `python main.py migrate`，必要时使用 `--fake-initial`。

## 生产运行

**推荐使用 main.py 作为唯一入口：**

```bash
python main.py --prod
```

或直接使用 Gunicorn：

```bash
gunicorn config.wsgi:application --bind 0.0.0.0:8000 --workers 4
```

## Nginx 配置示例

```nginx
server {
    listen 80;
    server_name your-domain.com;
    location /static/ {
        alias /path/to/project/staticfiles/;
    }
    location / {
        proxy_pass http://127.0.0.1:8000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## Docker / 容器部署

若 `migrate` 在系统检查阶段报错，可跳过检查：

```bash
python main.py migrate --noinput --skip-checks
```

建议使用 Python 3.9+。

## 安全

- DEBUG = False
- ALLOWED_HOSTS 配置正确
- HTTPS、数据库强密码、定期备份
