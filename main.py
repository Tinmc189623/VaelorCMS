#!/usr/bin/env python
# -*- coding: utf-8 -*-
"""
VaelorCMS 主程序入口 - 唯一启动入口

启动后自动检测系统是否已完成安装：
- 若已安装：正常进入 CMS 主业务流程
- 若未安装：自动跳转到安装程序（访问任意页面将重定向至 /install/）

用法：
  python main.py              # 启动服务（开发模式）
  python main.py --prod       # 启动服务（生产模式，Gunicorn）
  python main.py migrate      # 管理命令委托给 manage.py
"""
import os
import sys
from pathlib import Path

# 项目根目录（相对 main.py 所在目录）
BASE_DIR = Path(__file__).resolve().parent
INSTALLED_LOCK_PATH = BASE_DIR / 'config' / 'installed.lock'

try:
    from config.version import __version__ as CMS_VERSION
except ImportError:
    CMS_VERSION = 'Demo-26.02.13.26'


def is_installed() -> bool:
    """检测系统是否已完成安装"""
    return INSTALLED_LOCK_PATH.exists()


def run_server(use_gunicorn: bool = False) -> None:
    """启动 Web 服务"""
    os.environ.setdefault('DJANGO_SETTINGS_MODULE', 'config.settings')
    project_root = BASE_DIR
    os.chdir(project_root)
    if str(project_root) not in sys.path:
        sys.path.insert(0, str(project_root))

    if use_gunicorn:
        try:
            from gunicorn.app.base import BaseApplication

            class StandaloneApplication(BaseApplication):
                def __init__(self, app, options=None):
                    self.options = options or {}
                    self.application = app
                    super().__init__()

                def load_config(self):
                    for key, value in self.options.items():
                        self.cfg.set(key, value)

                def load(self):
                    return self.application

            from config.wsgi import application
            port = os.environ.get('PORT', '8000')
            StandaloneApplication(
                application,
                {'bind': f'0.0.0.0:{port}', 'workers': 4}
            ).run()
        except ImportError:
            _run_runserver()
    else:
        _run_runserver()


def _run_runserver() -> None:
    """使用 Django 内置 runserver 启动"""
    from django.core.management import execute_from_command_line
    port = os.environ.get('PORT', '8000')
    manage_py = BASE_DIR / 'manage.py'
    execute_from_command_line([str(manage_py), 'runserver', f'0.0.0.0:{port}'])


def main() -> None:
    # 管理命令（migrate、createsuperuser、collectstatic 等）委托给 manage.py
    if len(sys.argv) > 1 and sys.argv[1] not in ('run', '--prod'):
        os.chdir(BASE_DIR)
        # 过滤掉仅用于启动服务的参数，避免传给 Django 命令导致报错
        skip_args = {'--prod', '--pord'}  # --pord 为 --prod 常见拼写错误
        delegate_args = [a for a in sys.argv[1:] if a not in skip_args]
        os.execv(sys.executable, [sys.executable, 'manage.py'] + delegate_args)
        return

    # 安装状态检测
    installed = is_installed()
    print(f'[VaelorCMS] 版本 {CMS_VERSION}')
    if installed:
        print('[VaelorCMS] 已安装，正在启动主业务流程...')
    else:
        print('[VaelorCMS] 未安装，将进入安装向导模式')
        print('[VaelorCMS] 启动后请访问 http://127.0.0.1:8000/ 自动跳转至安装页面')

    # 启动服务（中间件会根据安装状态自动重定向）
    use_gunicorn = len(sys.argv) > 1 and sys.argv[1] == '--prod'
    run_server(use_gunicorn=use_gunicorn)


if __name__ == '__main__':
    main()
