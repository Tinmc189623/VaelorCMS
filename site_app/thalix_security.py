"""
Thalix - VaelorCMS 安全审计模块
自研组件，用于扫描配置与运行环境中的安全隐患
"""
from django.conf import settings
from pathlib import Path


class SecurityCheckResult:
    """单次检查结果"""
    def __init__(self, check_id, name, passed, message, severity='warning'):
        self.check_id = check_id
        self.name = name
        self.passed = passed
        self.message = message
        self.severity = severity  # info, warning, critical


def run_security_audit():
    """
    执行完整安全审计，返回检查结果列表
    """
    results = []
    
    # 1. SECRET_KEY 强度
    sk = getattr(settings, 'SECRET_KEY', '')
    if not sk or sk in ('dev-secret-change-in-production', 'changeme', 'secret', 'key', 'django-insecure'):
        results.append(SecurityCheckResult(
            'secret_key', 'SECRET_KEY 强度',
            False, 'SECRET_KEY 过于简单，生产环境必须设置强随机值（环境变量 DJANGO_SECRET_KEY）',
            'critical'
        ))
    elif len(sk) < 32:
        results.append(SecurityCheckResult(
            'secret_key', 'SECRET_KEY 强度',
            False, 'SECRET_KEY 长度建议至少 50 字符',
            'warning'
        ))
    else:
        results.append(SecurityCheckResult('secret_key', 'SECRET_KEY 强度', True, '已配置', 'info'))
    
    # 2. DEBUG 模式
    if settings.DEBUG:
        results.append(SecurityCheckResult(
            'debug_mode', 'DEBUG 模式',
            False, '生产环境必须设置 DEBUG=False',
            'critical'
        ))
    else:
        results.append(SecurityCheckResult('debug_mode', 'DEBUG 模式', True, '已关闭', 'info'))
    
    # 3. ALLOWED_HOSTS
    hosts = getattr(settings, 'ALLOWED_HOSTS', [])
    if not hosts or '*' in hosts:
        results.append(SecurityCheckResult(
            'allowed_hosts', 'ALLOWED_HOSTS',
            False, "ALLOWED_HOSTS 不应使用 '*'，应明确指定域名",
            'warning'
        ))
    else:
        results.append(SecurityCheckResult('allowed_hosts', 'ALLOWED_HOSTS', True, '已配置', 'info'))
    
    # 4. CSRF
    if not getattr(settings, 'CSRF_COOKIE_HTTPONLY', False):
        results.append(SecurityCheckResult(
            'csrf_httponly', 'CSRF Cookie HttpOnly',
            False, '建议启用 CSRF_COOKIE_HTTPONLY',
            'warning'
        ))
    else:
        results.append(SecurityCheckResult('csrf_httponly', 'CSRF Cookie', True, '已加固', 'info'))
    
    # 5. Session 安全
    if not getattr(settings, 'SESSION_COOKIE_HTTPONLY', True):
        results.append(SecurityCheckResult(
            'session_httponly', 'Session Cookie HttpOnly',
            False, '建议启用 SESSION_COOKIE_HTTPONLY',
            'warning'
        ))
    else:
        results.append(SecurityCheckResult('session_httponly', 'Session Cookie', True, '已加固', 'info'))
    
    # 6. HTTPS
    if not getattr(settings, 'SESSION_COOKIE_SECURE', False) and not settings.DEBUG:
        results.append(SecurityCheckResult(
            'https', 'HTTPS / Secure Cookie',
            False, '生产环境建议启用 HTTPS 与 SESSION_COOKIE_SECURE',
            'warning'
        ))
    else:
        results.append(SecurityCheckResult('https', 'HTTPS', True, '已配置或开发模式', 'info'))
    
    # 7. X-Content-Type-Options
    results.append(SecurityCheckResult(
        'x_content_type', 'X-Content-Type-Options',
        True, '由 SecurityHeadersMiddleware 处理', 'info'
    ))
    
    # 8. 安装锁
    lock_path = Path(settings.BASE_DIR) / 'config' / 'installed.lock'
    if not lock_path.exists() and not settings.DEBUG:
        results.append(SecurityCheckResult(
            'install_lock', '安装向导',
            False, '未检测到 installed.lock，请完成安装向导',
            'warning'
        ))
    
    return results


def get_audit_summary():
    """返回审计摘要：通过数、警告数、严重数"""
    results = run_security_audit()
    passed = sum(1 for r in results if r.passed)
    warnings = sum(1 for r in results if not r.passed and r.severity == 'warning')
    critical = sum(1 for r in results if not r.passed and r.severity == 'critical')
    return {'passed': passed, 'warnings': warnings, 'critical': critical, 'total': len(results)}
