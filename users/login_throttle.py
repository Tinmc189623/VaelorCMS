"""
登录失败锁定 - 按管理员配置的 login_max_attempts、login_lockout_minutes 执行
使用 Django cache 存储，生产环境建议配置数据库或 Redis 缓存
"""
from django.core.cache import cache
from django.utils import timezone


def _get_ip(request):
    """获取客户端 IP（考虑 X-Forwarded-For）"""
    xff = request.META.get('HTTP_X_FORWARDED_FOR')
    if xff:
        return xff.split(',')[0].strip()
    return request.META.get('REMOTE_ADDR', '0.0.0.0')


def _cache_key(ip):
    return f'login_lockout:{ip}'


def _get_config():
    """读取登录锁定配置，避免多处重复调用"""
    try:
        from site_app.settings_service import get
        return int(get('login_max_attempts', '5')), int(get('login_lockout_minutes', '15'))
    except (ValueError, ImportError):
        return 5, 15


def is_locked(request):
    """检查当前 IP 是否处于锁定状态"""
    max_attempts, lockout_mins = _get_config()
    ip = _get_ip(request)
    key = _cache_key(ip)
    data = cache.get(key)
    if not data:
        return False, 0

    fail_count, locked_until = data
    if locked_until and timezone.now().timestamp() < locked_until:
        return True, int(locked_until - timezone.now().timestamp())
    if fail_count >= max_attempts:
        # 达到上限，开始锁定
        locked_until = timezone.now().timestamp() + lockout_mins * 60
        cache.set(key, (fail_count, locked_until), timeout=lockout_mins * 60 + 60)
        return True, lockout_mins * 60
    return False, 0


def record_fail(request):
    """记录一次登录失败"""
    max_attempts, lockout_mins = _get_config()
    ip = _get_ip(request)
    key = _cache_key(ip)
    data = cache.get(key) or (0, None)
    fail_count, _ = data
    fail_count += 1

    if fail_count >= max_attempts:
        locked_until = timezone.now().timestamp() + lockout_mins * 60
        cache.set(key, (fail_count, locked_until), timeout=lockout_mins * 60 + 60)
    else:
        cache.set(key, (fail_count, None), timeout=3600)


def clear_on_success(request):
    """登录成功后清除失败计数"""
    ip = _get_ip(request)
    cache.delete(_cache_key(ip))
