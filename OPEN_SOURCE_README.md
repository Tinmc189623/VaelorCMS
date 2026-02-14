# VaelorCMS 开源发布包

本目录由 `scripts/build_opensource.py` 自动生成，包含可安全公开的源代码与文档。

## 已排除的敏感内容

- `config/config.ini` - 数据库密码等配置
- `config/installed.lock` - 安装状态
- `config/install_temp.db` - 安装向导临时数据库
- `venv/` - 虚拟环境
- `storage/uploads/` - 用户上传文件
- `__pycache__/` - Python 缓存
- `staticfiles/` - 构建产物（可运行 collectstatic 重新生成）

## 使用说明

1. 复制 `config/config.ini.sample` 为 `config/config.ini` 并填写配置
2. 运行 `pip install -r requirements.txt`
3. 运行 `python main.py migrate` 初始化数据库
4. 运行 `python main.py` 启动（未安装将进入安装向导）
