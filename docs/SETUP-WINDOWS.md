# VaelorCMS Windows 本地安装指南

## 一、安装数据库（推荐 MariaDB 8.x 或 MySQL 8.x）

### 方式 A：winget 自动安装（需网络）

在 PowerShell（管理员）中执行：

```powershell
# 优先 MariaDB 8.x
winget install MariaDB.Server --source winget --accept-package-agreements --accept-source-agreements

# 若失败，尝试 MySQL 8.x
winget install Oracle.MySQL --source winget --accept-package-agreements --accept-source-agreements
```

或运行项目脚本：

```powershell
cd j:\APP\Web\VaelorCMS
.\scripts\setup-database.ps1
```

### 方式 B：手动下载安装

1. **MariaDB 8.x**（推荐）：https://mariadb.org/download/
2. **MySQL 8.x**：https://dev.mysql.com/downloads/mysql/

安装时记住 root 密码，安装完成后服务会自动启动。

### 创建数据库与用户

以 MariaDB/MySQL 为例，在命令行或 HeidiSQL/phpMyAdmin 中执行：

```sql
CREATE DATABASE vaelor_cms CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'vaelor_user'@'localhost' IDENTIFIED BY '你的密码';
GRANT ALL PRIVILEGES ON vaelor_cms.* TO 'vaelor_user'@'localhost';
FLUSH PRIVILEGES;
```

## 二、安装 Python 依赖

```powershell
cd j:\APP\Web\VaelorCMS
pip install -r requirements.txt
```

## 三、配置与运行

1. 复制配置：`copy config\config.ini.sample config\config.ini`
2. 编辑 `config\config.ini`，填写 [database] 中的密码等
3. 初始化并启动：

```powershell
python main.py migrate --skip-checks
python main.py
```

4. 浏览器访问：http://127.0.0.1:8000/

## 四、无数据库快速体验（SQLite）

若暂不安装 MariaDB/MySQL，可临时使用 SQLite：

1. 删除或重命名 `config\config.ini`
2. 执行 `python main.py migrate --skip-checks`
3. 执行 `python main.py`，通过 Web 安装向导完成配置
4. 安装向导会要求 MySQL，此时需先安装数据库后再执行安装
