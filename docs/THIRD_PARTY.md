# VaelorCMS 第三方与自研组件声明

本文档列出 VaelorCMS 使用的所有外部开源库及自研组件，便于合规与审计。

## 一、外部开源库（Third-Party Libraries）

| 库名 | 版本 | 用途 | 许可证 |
|------|------|------|--------|
| Django | 4.2.x | Web 框架、ORM、模板、认证 | BSD-3-Clause |
| Pillow | ≥10.0 | 图片处理、头像裁剪 | HPND |
| Gunicorn | ≥21.0 | WSGI 生产服务器 | MIT |
| mysqlclient | ≥2.2 | MySQL 数据库驱动 | GPL |
| PyMySQL | ≥1.1 | 纯 Python MySQL 驱动 | MIT |
| WhiteNoise | ≥6.0 | 静态文件服务 | MIT |
| Redis | ≥4.5 | 缓存、限流、会话后端 | MIT |
| bleach | ≥6.0 | HTML 净化，防 XSS | Apache-2.0 |
| defusedxml | ≥0.7 | 安全 XML 解析，防 XXE | PSF |

## 二、自研组件（In-House Components）

以下为 VaelorCMS 项目内自研模块，非第三方库，命名已检索避免与现有项目冲突。

### 1. Vaelor Core（vaelor_core）

**说明**：VaelorCMS 内核，提供配置抽象、钩子系统、安全工具等核心能力。

**模块**：
- `vaelor_core.config`：配置读写抽象
- `vaelor_core.hooks`：插件钩子注册与执行
- `vaelor_core.security`：输入验证、安全字符串处理

### 2. Thalix（site_app.thalix_security）

**说明**：安全审计模块，扫描 SECRET_KEY、DEBUG、ALLOWED_HOSTS、CSRF、Session、HTTPS 等配置项，输出检查报告。

**入口**：管理后台 → 安全审计

### 3. 自研 HTML 净化（site_app.html_sanitizer）

**说明**：自定义页面 HTML 净化，防 XSS。优先使用 bleach，无 bleach 时回退到自研正则净化。

---

## 三、命名检索说明

- **Vaelor**：项目主品牌，未发现同名软件商标冲突。
- **Thalix**：安全模块名，检索未发现同名 Python 安全库。
- **Vaelor Core**：内核名称，由 VaelorCMS 去掉 CMS 加 Core 构成。

---

*最后更新：2025*
