# 快速上手 VaelorCMS

> 5 分钟搞清楚怎么跑起来、怎么用。不废话，直接上。

---

## 第一次部署？看这里

### 1. 准备环境

- Python 3.9+
- MySQL 5.7+（或 MariaDB 10.2+）
- 先建好一个空数据库，比如：`CREATE DATABASE vaelor_cms CHARACTER SET utf8mb4;`

### 2. 装依赖、跑起来

```bash
cd /path/to/VaelorCMS
python -m venv venv
venv\Scripts\activate   # Windows
# source venv/bin/activate  # Mac/Linux

pip install -r requirements.txt
python main.py migrate
python main.py
```

### 3. 打开浏览器

访问 `http://127.0.0.1:8000/`，会自动跳到安装向导。按提示一步步填：

1. 同意许可协议
2. 环境检测（一般都能过）
3. 填数据库信息（主机、端口、库名、用户名、密码）
4. 填站点名、管理员账号和密码
5. 点「执行安装」

装完后**重启一次** `python main.py`，再用管理员账号登录，就能进管理后台了。

---

## 已经装好了？看这里

### 普通用户

- **注册**：点顶栏「注册」，填用户名密码
- **发帖**：登录 → 论坛 → 发帖
- **发代码**：登录 → 代码 → 提交代码
- **写文章**：登录 → 文章 → 发布文章
- **改资料**：用户中心 → 设置

详细说明见 [用户指南](USER_GUIDE.md)。

### 管理员

- **管理入口**：登录后顶栏多一个「管理」
- **站点设置**：管理 → 站点设置，改站点名、安全策略、注册开关、维护模式等
- **用户/帖子/代码**：管理里都有入口

详细说明见 [管理员指南](ADMIN_GUIDE.md)。

---

## 生产环境怎么跑？

```bash
python main.py --prod
```

或用 Gunicorn：

```bash
gunicorn config.wsgi:application --bind 0.0.0.0:8000 --workers 4
```

前面加 Nginx 做反向代理，配 HTTPS。详见 [部署指南](DEPLOYMENT.md)、[云平台部署](DEPLOY-CLOUD.md)。

---

## 文档导航

| 文档 | 给谁看 | 讲啥 |
|------|--------|------|
| [INSTALL.md](INSTALL.md) | 部署的人 | 完整安装步骤（Web 向导 + 手动） |
| [USER_GUIDE.md](USER_GUIDE.md) | 普通用户 | 注册、发帖、设置、常见问题 |
| [ADMIN_GUIDE.md](ADMIN_GUIDE.md) | 管理员 | 站点设置、安全、维护模式 |
| [CONFIG.md](CONFIG.md) | 改配置的人 | config.ini 说明 |
| [API.md](API.md) | 开发者 | API 接口 |
| [DEPLOYMENT.md](DEPLOYMENT.md) | 运维 | 生产部署 |
| [DEPLOY-CLOUD.md](DEPLOY-CLOUD.md) | 上云的人 | 云平台部署 |
