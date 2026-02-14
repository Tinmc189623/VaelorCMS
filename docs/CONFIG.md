# 配置说明

## config/config.ini

Django 从 `config/config.ini` 读取业务配置（见 `config/settings.py`）。

### [database]

host、port、dbname、username、password、charset

### [site]

name、base_url、session_name、cookie_path、trusted_origins

- **trusted_origins**：CSRF 可信来源，HTTPS 反向代理部署时必填，逗号分隔

**示例**（config.ini 中）：

```ini
[site]
trusted_origins = https://www.example.com,https://example.com
```

**示例**（环境变量）：

```bash
CSRF_TRUSTED_ORIGINS=https://www.example.com,https://example.com
```

### [security]

login_max_attempts、login_lockout_minutes、force_https

### [limits]

list_limit、bbs_content_max、code_snippet_max

### [pagination]

per_page

### [cache]

driver、host、port、prefix、password。driver=redis 时启用 Redis 缓存。环境变量 REDIS_HOST、REDIS_PORT、REDIS_PASSWORD 可覆盖。

## config/settings.py

Django 主配置。SECRET_KEY、DEBUG、ALLOWED_HOSTS 可通过环境变量覆盖：

- DJANGO_SECRET_KEY：生产环境必设强随机值
- DJANGO_DEBUG（1 为 True，默认 0 即生产模式；本地开发可设 1）
- DJANGO_HTTPS（1/true/yes）：启用 Session/CSRF Cookie Secure，HTTPS 部署时建议设置
- ALLOWED_HOSTS（逗号分隔）
