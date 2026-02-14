# VaelorCMS 安装教程

VaelorCMS 是介于轻量与重量级之间的开源 CMS，功能完整且易于部署。说白了：比纯静态站功能多，比 Drupal 那种大块头好上手。

本文档提供完整安装步骤，支持 **Web 安装向导**（推荐）和 **手动配置** 两种方式。第一次装建议用 Web 向导，跟着浏览器点就行；熟手可以直接改配置文件。

---

## 一、环境要求

| 项目 | 要求 |
|------|------|
| Python | 3.9 及以上 |
| MySQL | 5.7+ 或 MariaDB 10.2+ |
| 操作系统 | Windows / Linux / macOS |

---

## 二、方式一：Web 安装向导（推荐）

适合首次部署，通过浏览器完成配置，不用碰命令行和配置文件。装完就能用。

### 步骤 1：获取代码

将 VaelorCMS 源码上传至服务器或本地目录，例如：

```
/path/to/VaelorCMS/
├── main.py
├── config/
├── install_app/
└── ...
```

### 步骤 2：创建虚拟环境

```bash
# 进入项目目录
cd /path/to/VaelorCMS

# 创建虚拟环境
python -m venv venv

# 激活虚拟环境
# Windows:
venv\Scripts\activate
# Linux / macOS:
source venv/bin/activate
```

### 步骤 3：安装依赖

```bash
pip install -r requirements.txt
```

### 步骤 4：初始化安装数据库

安装向导在首次运行前需要临时数据库（SQLite），执行：

```bash
python main.py migrate
```

### 步骤 5：启动服务

```bash
# 开发模式（默认 0.0.0.0:8000）
python main.py

# 生产模式（Gunicorn）
python main.py --prod
```

启动后会看到类似输出：

```
[VaelorCMS] 未安装，将进入安装向导模式
[VaelorCMS] 启动后请访问 http://127.0.0.1:8000/ 自动跳转至安装页面
```

### 步骤 6：访问安装向导

在浏览器中访问：

- 本地：`http://127.0.0.1:8000/`
- 服务器：`http://你的服务器IP:8000/`

系统会自动跳转到安装页面 `/install/`。

### 步骤 7：按向导完成安装

| 步骤 | 内容 | 说明 |
|------|------|------|
| 1. 许可协议 | 阅读并勾选同意 | 必须同意方可继续 |
| 2. 站点配置 | 站点名称、管理员账号、可选 MySQL | 默认 SQLite，无需修改配置文件 |
| 3. 执行安装 | 自动执行 | 迁移数据库、创建管理员 |

**零配置安装：**

- 默认使用 SQLite，无需创建数据库，安装后**无需重启**
- 可选勾选「使用 MySQL」填写连接信息（需先创建数据库）
- 站点名称等可在安装后于后台「站点设置」中修改

### 步骤 8：重启服务（仅使用 MySQL 时）

若安装时选择了 MySQL，需重启应用以使配置生效。使用 SQLite 则无需重启。

### 步骤 9：收集静态文件（生产环境）

```bash
python main.py collectstatic --noinput
```

将静态文件收集到 `staticfiles/` 目录，供 Nginx 等 Web 服务器使用。

---

## 三、方式二：手动配置

适合已经摸过 config.ini 的人，或者要写脚本自动部署的场景。比 Web 向导快，但得自己保证配置没错。

### 步骤 1～3：同 Web 安装向导

完成环境准备、依赖安装。

### 步骤 4：复制并编辑配置文件

```bash
cp config/config.ini.sample config/config.ini
```

编辑 `config/config.ini`，至少填写：

```ini
[database]
host = 127.0.0.1
port = 3306
dbname = 你的数据库名
username = 数据库用户名
password = 数据库密码
charset = utf8mb4

[site]
name = 你的站点名称
base_url = /
session_name = VAELOR_SESS
cookie_path = /
```

### 步骤 5：创建安装锁文件

```bash
# Linux / macOS
touch config/installed.lock

# Windows (PowerShell)
New-Item config/installed.lock -ItemType File
```

表示系统已安装，跳过 Web 安装向导。

### 步骤 6：数据库迁移与创建管理员

```bash
python main.py migrate
python main.py createsuperuser
```

按提示输入管理员用户名和密码。

### 步骤 7：启动服务

```bash
python main.py        # 开发模式
python main.py --prod # 生产模式
```

---

## 四、常见问题

### 1. 数据库连接失败

先确认这几件事：MySQL 服务有没有启动？`config.ini` 里的主机、端口、用户名、密码对不对？数据库有没有事先建好？用户有没有该库的读写权限？常见坑：端口写错、密码带特殊字符没转义、数据库名不存在。

### 2. 环境检测不通过

- **Python 版本**：必须 3.9 及以上，老版本不行
- **MySQL 驱动**：装一个就行，`pip install mysqlclient` 或 `pip install pymysql`。mysqlclient 要编译，装不上就用 PyMySQL
- **目录权限**：`config/`、`storage/uploads/` 要能写，不然装完写不进去

### 3. 安装后访问 404

多半是没重启。装完一定要停掉服务再开一次，不然还在用旧的 SQLite 配置。再检查一下 `config/installed.lock` 有没有生成。

### 4. 如何重新安装？

删掉 `config/config.ini` 和 `config/installed.lock`，重启服务，就会重新进安装向导。注意：这会清掉你填的配置，数据库里的数据不会自动删，要彻底重来得自己清库。

---

## 五、生产环境部署

生产环境建议配合 Nginx 反向代理，详见 [DEPLOYMENT.md](DEPLOYMENT.md)。

简要流程：

1. 使用 `python main.py --prod` 或 Gunicorn 启动应用
2. 配置 Nginx 反向代理到 `127.0.0.1:8000`
3. 设置 `DEBUG=False`、配置 HTTPS

---

## 六、相关文档

- [DEPLOYMENT.md](DEPLOYMENT.md) - 部署指南
- [CONFIG.md](CONFIG.md) - 配置说明
- [API.md](API.md) - API 规范
