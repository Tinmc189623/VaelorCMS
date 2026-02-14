"""
API 限流 - 按 IP 限制 /api/v1/ 请求频率
使用 Django cache，配置项 api_rate_limit（每分钟请求数，0 表示不限）
"""
from django.core.cache import cache
from django.http import JsonResponse
from django.utils import timezone


def _get_ip(request):
    xff = request.META.get('HTTP_X_FORWARDED_FOR')
    if xff:
        return xff.split(',')[0].strip()
    return request.META.get('REMOTE_ADDR', '0.0.0.0')


def _get_limit():
    try:
        from .settings_service import get
        return int(get('api_rate_limit', '60'))
    except (ValueError, ImportError):
        return 60


def api_throttle_middleware(get_response):
    """API 限流中间件，仅对 /api/v1/ 生效，健康检查接口不限流"""
    def middleware(request):
        if not request.path.startswith('/api/v1/'):
            return get_response(request)
        if request.path == '/api/v1/health/':
            return get_response(request)
        limit = _get_limit()
        if limit <= 0:
            return get_response(request)
        ip = _get_ip(request)
        key = f'api_throttle:{ip}'
        now = timezone.now()
        minute_key = now.strftime('%Y%m%d%H%M')
        cache_key = f'{key}:{minute_key}'
        count = cache.get(cache_key, 0)
        if count >= limit:
            return JsonResponse({
                'error': 'rate_limit_exceeded',
                'message': f'请求过于频繁，请稍后再试（每分钟 {limit} 次）',
            }, status=429)
        cache.set(cache_key, count + 1, timeout=65)
        return get_response(request)
    return middleware
