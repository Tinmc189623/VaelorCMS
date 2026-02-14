"""
VaelorCMS 安装向导 - 零配置安装，仅需在 Web 填写信息
流程：许可协议 → 一站式配置（站点 + 管理员，可选 MySQL）→ 执行安装
"""
import os
import sys
import subprocess
from pathlib import Path

from django.shortcuts import render, redirect
from django.views.decorators.http import require_http_methods
from django.conf import settings
from django.contrib import messages

BASE_DIR = Path(settings.BASE_DIR)


def _sanitize_ini_value(s: str) -> str:
    """移除可能破坏 INI 格式的字符"""
    if not s:
        return s
    return ''.join(c for c in s if c not in '\n\r')


def _is_installed():
    lock = getattr(settings, 'INSTALLED_LOCK_PATH', None)
    return lock and lock.exists()


def _run_environment_checks():
    """环境检测"""
    checks = []
    py_ver = sys.version_info
    py_ok = py_ver >= (3, 9)
    checks.append(('Python 3.9+', f'{py_ver.major}.{py_ver.minor}.{py_ver.micro}', py_ok))

    try:
        import django
        dj_ver = django.VERSION
        checks.append(('VaelorCMS 框架', f'{dj_ver[0]}.{dj_ver[1]}', dj_ver >= (4, 0)))
    except ImportError:
        checks.append(('VaelorCMS 框架', '未安装', False))

    try:
        import MySQLdb
        checks.append(('MySQL 驱动', '已安装', True))
    except ImportError:
        try:
            import pymysql
            pymysql.install_as_MySQLdb()
            checks.append(('MySQL 驱动', '已安装', True))
        except ImportError:
            checks.append(('MySQL 驱动', '未安装', False))

    config_dir = BASE_DIR / 'config'
    config_writable = config_dir.exists() and os.access(config_dir, os.W_OK)
    checks.append(('config/ 可写', '是' if config_writable else '否', config_writable))

    storage_dir = BASE_DIR / 'storage' / 'uploads'
    storage_ok = True
    try:
        storage_dir.mkdir(parents=True, exist_ok=True)
        storage_ok = os.access(storage_dir, os.W_OK)
    except OSError:
        storage_ok = False
    checks.append(('storage/uploads/ 可写', '是' if storage_ok else '否', storage_ok))

    return checks


# ========== 步骤 1：许可协议 ==========
@require_http_methods(['GET', 'POST'])
def step_license(request):
    if _is_installed():
        return redirect('home')
    step = 1
    if request.method == 'POST':
        if request.POST.get('agree') == '1':
            return redirect('install_step2')
        messages.error(request, '请阅读并同意许可协议')
    return render(request, 'install/step1_license.html', {'step': step})


# ========== 步骤 2：一站式配置（环境 + 站点 + 可选 MySQL）==========
@require_http_methods(['GET', 'POST'])
def step_setup(request):
    if _is_installed():
        return redirect('home')
    step = 2
    checks = _run_environment_checks()
    all_ok = all(c[2] for c in checks)

    if request.method == 'POST':
        if not all_ok:
            messages.error(request, '请先解决环境检测中的问题')
            return render(request, 'install/step2_setup.html', {
                'step': step, 'checks': checks, 'all_ok': all_ok,
                'site_name': request.POST.get('site_name', 'VaelorCMS'),
                'admin_user': request.POST.get('admin_user', ''),
                'use_mysql': request.POST.get('use_mysql') == '1',
            })

        site_name = request.POST.get('site_name', 'VaelorCMS').strip() or 'VaelorCMS'
        admin_user = request.POST.get('admin_user', '').strip()
        admin_pass = request.POST.get('admin_pass', '')
        admin_pass2 = request.POST.get('admin_pass2', '')
        use_mysql = request.POST.get('use_mysql') == '1'

        if not admin_user:
            messages.error(request, '请填写管理员用户名')
        elif len(admin_pass) < 6:
            messages.error(request, '密码至少 6 位')
        elif admin_pass != admin_pass2:
            messages.error(request, '两次密码不一致')
        else:
            request.session['install_site'] = {
                'site_name': site_name,
                'admin_user': admin_user,
                'admin_pass': admin_pass,
            }
            if use_mysql:
                db_host = request.POST.get('db_host', '127.0.0.1').strip()
                db_port = request.POST.get('db_port', '3306').strip() or '3306'
                db_name = request.POST.get('db_name', '').strip()
                db_user = request.POST.get('db_user', '').strip()
                db_pass = request.POST.get('db_pass', '')
                if not db_name or not db_user:
                    messages.error(request, '使用 MySQL 时请填写数据库名和用户名')
                elif not db_port.isdigit() or not (1 <= int(db_port) <= 65535):
                    messages.error(request, '端口号须为 1–65535')
                else:
                    try:
                        import MySQLdb
                    except ImportError:
                        try:
                            import pymysql
                            pymysql.install_as_MySQLdb()
                        except ImportError:
                            messages.error(request, 'MySQL 驱动未安装')
                            return render(request, 'install/step2_setup.html', {
                                'step': step, 'checks': checks, 'all_ok': all_ok,
                                'site_name': site_name, 'admin_user': admin_user,
                                'use_mysql': True,
                            })
                    try:
                        import MySQLdb as db
                        conn = db.connect(
                            host=db_host, port=int(db_port), user=db_user,
                            password=db_pass, database=db_name, charset='utf8mb4'
                        )
                        conn.close()
                        request.session['install_db'] = {
                            'host': db_host, 'port': db_port, 'dbname': db_name,
                            'username': db_user, 'password': db_pass, 'charset': 'utf8mb4',
                        }
                        return redirect('install_step3')
                    except Exception as e:
                        messages.error(request, f'数据库连接失败：{e}')
            else:
                request.session.pop('install_db', None)
                return redirect('install_step3')

        return render(request, 'install/step2_setup.html', {
            'step': step, 'checks': checks, 'all_ok': all_ok,
            'site_name': site_name, 'admin_user': admin_user,
            'use_mysql': use_mysql,
        })

    return render(request, 'install/step2_setup.html', {
        'step': step, 'checks': checks, 'all_ok': all_ok,
        'site_name': 'VaelorCMS', 'admin_user': '', 'use_mysql': False,
    })


# ========== 步骤 3：执行安装 ==========
@require_http_methods(['GET', 'POST'])
def step_install(request):
    if _is_installed():
        return redirect('home')
    if 'install_site' not in request.session:
        return redirect('install_step1')
    step = 3
    use_mysql = 'install_db' in request.session

    if request.method == 'GET':
        return render(request, 'install/step3_install.html', {
            'step': step, 'use_mysql': use_mysql,
        })

    site = request.session['install_site']
    errors = []

    if use_mysql:
        db = request.session['install_db']
        config_content = f"""[database]
host = {_sanitize_ini_value(db['host'])}
port = {db['port']}
dbname = {_sanitize_ini_value(db['dbname'])}
username = {_sanitize_ini_value(db['username'])}
password = {_sanitize_ini_value(db['password'])}
charset = {db['charset']}

[site]
name = {_sanitize_ini_value(site['site_name'])}
base_url = /
session_name = VAELOR_SESS
cookie_path = /

[security]
login_max_attempts = 5
login_lockout_minutes = 15
force_https = 0

[limits]
list_limit = 100
bbs_content_max = 50000
code_snippet_max = 100000

[pagination]
per_page = 20

[cache]
driver =
host = 127.0.0.1
port = 6379
prefix = vaelor_
"""
        config_path = BASE_DIR / 'config' / 'config.ini'
        try:
            config_path.write_text(config_content, encoding='utf-8')
        except Exception as e:
            errors.append(f'写入配置失败：{e}')

    if not errors:
        try:
            subprocess.run(
                [sys.executable, 'main.py', 'migrate', '--noinput', '--skip-checks'],
                cwd=BASE_DIR,
                capture_output=True,
                text=True,
                timeout=120,
            )
        except subprocess.TimeoutExpired:
            errors.append('数据库迁移超时')
        except Exception as e:
            errors.append(f'数据库迁移失败：{e}')

    if not errors:
        env = os.environ.copy()
        env['INSTALL_ADMIN_USER'] = site['admin_user']
        env['INSTALL_ADMIN_PASS'] = site['admin_pass']
        try:
            subprocess.run(
                [sys.executable, 'main.py', 'create_install_admin'],
                cwd=BASE_DIR,
                env=env,
                capture_output=True,
                text=True,
                timeout=30,
            )
        except Exception as e:
            errors.append(f'创建管理员失败：{e}')

    if not errors:
        try:
            from site_app.settings_service import set
            set('site_name', site['site_name'], 'general')
        except Exception:
            pass

    if not errors:
        lock_path = BASE_DIR / 'config' / 'installed.lock'
        try:
            lock_path.write_text('installed', encoding='utf-8')
        except Exception as e:
            errors.append(f'创建安装锁失败：{e}')

    if errors:
        return render(request, 'install/step3_install.html', {
            'step': step, 'use_mysql': use_mysql, 'errors': errors,
        })

    request.session.pop('install_db', None)
    request.session.pop('install_site', None)
    if use_mysql:
        messages.success(request, '安装完成！请重启应用服务器以使 MySQL 配置生效，然后使用管理员账号登录。')
    else:
        messages.success(request, '安装完成！请使用管理员账号登录后台，在「站点设置」中修改站点信息。')
    return redirect('home')
