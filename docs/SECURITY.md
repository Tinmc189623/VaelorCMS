# VaelorCMS 安全说明

## 安全架构

### 1. 防护层

| 层级 | 组件 | 说明 |
|------|------|------|
| 请求 | SecurityHeadersMiddleware | X-Content-Type-Options、X-Frame-Options、Referrer-Policy、Permissions-Policy、CSP、HSTS |
| 请求 | CsrfProxyTrustMiddleware | 反向代理下自动补充 CSRF 可信来源 |
| 请求 | api_throttle_middleware | API 限流，防暴力请求 |
| 认证 | login_throttle | 登录失败锁定，防暴力破解 |
| 输出 | html_sanitizer | HTML 净化，防 XSS（bleach + 自研正则） |
| 输出 | 模板自动转义 | Django 模板默认转义 |

### 2. 安全审计（Thalix）

管理后台 → 安全审计 可查看：

- SECRET_KEY 强度
- DEBUG 模式
- ALLOWED_HOSTS 配置
- CSRF / Session Cookie 安全
- HTTPS 与 Secure Cookie

### 3. 生产部署建议

1. 设置 `DJANGO_SECRET_KEY` 环境变量（≥50 字符随机串）
2. 设置 `DEBUG=0` 或 `DJANGO_DEBUG=0`
3. 配置 `ALLOWED_HOSTS` 为实际域名，避免 `*`
4. 使用 HTTPS，并配置 `CSRF_TRUSTED_ORIGINS`
5. 定期运行安全审计，修复未通过项

### 4. 已知限制

- CSP 使用 `unsafe-inline` 以支持主题定制与阅读进度条，若需更严格策略可自行调整
- 安装向导阶段 `ALLOWED_HOSTS` 为 `*`，安装完成后建议通过 config.ini 或环境变量收紧
