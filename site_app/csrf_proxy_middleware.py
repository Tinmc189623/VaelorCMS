"""
反向代理下自动补充 CSRF 可信来源

当请求经 HTTPS 反向代理（X-Forwarded-Proto: https）时，
根据 Host / X-Forwarded-Host 自动添加可信来源，解决云平台部署 403。
"""
from django.conf import settings


class CsrfProxyTrustMiddleware:
    """根据请求自动补充 CSRF_TRUSTED_ORIGINS"""

    def __init__(self, get_response):
        self.get_response = get_response

    def __call__(self, request):
        if getattr(settings, 'CSRF_TRUSTED_ORIGINS_AUTO', True):
            proto = request.META.get('HTTP_X_FORWARDED_PROTO', '').strip().lower()
            if proto == 'https':
                host = (request.META.get('HTTP_X_FORWARDED_HOST') or '').strip() or request.get_host()
                if host:
                    origin = f'https://{host}'
                    trusted = list(settings.CSRF_TRUSTED_ORIGINS)
                    if origin not in trusted:
                        trusted.append(origin)
                        settings.CSRF_TRUSTED_ORIGINS = trusted
        return self.get_response(request)
