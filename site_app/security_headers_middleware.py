"""
安全响应头中间件 - 增强 XSS、点击劫持、HTTPS 等防护
"""
import re
from django.utils.deprecation import MiddlewareMixin
from django.http import HttpResponsePermanentRedirect


class SecurityHeadersMiddleware(MiddlewareMixin):
    """添加安全相关 HTTP 响应头，支持强制 HTTPS"""

    def process_request(self, request):
        # 强制 HTTPS（当管理员开启 force_https 时）
        if request.is_secure():
            return None
        try:
            from .settings_service import get
            if get('force_https', '0') == '1':
                url = request.build_absolute_uri(request.get_full_path())
                if url.startswith('http://'):
                    return HttpResponsePermanentRedirect(url.replace('http://', 'https://', 1))
        except Exception:
            pass
        return None

    def process_response(self, request, response):
        # 防止 MIME 类型嗅探
        if 'X-Content-Type-Options' not in response:
            response['X-Content-Type-Options'] = 'nosniff'
        # 防止点击劫持（Django 默认已有，此处确保）
        if 'X-Frame-Options' not in response:
            response['X-Frame-Options'] = 'DENY'
        # 控制 Referer 泄露
        if 'Referrer-Policy' not in response:
            response['Referrer-Policy'] = 'strict-origin-when-cross-origin'
        # 禁用部分浏览器特性，减少攻击面
        if 'Permissions-Policy' not in response:
            response['Permissions-Policy'] = 'geolocation=(), microphone=(), camera=(), payment=()'
        # Content-Security-Policy（生产环境启用，可通过 SECURE_CSP=0 关闭）
        if 'Content-Security-Policy' not in response:
            try:
                from django.conf import settings
                if not getattr(settings, 'DEBUG', True) and getattr(settings, 'SECURE_CSP', True):
                    # 宽松策略：允许内联脚本/样式（主题、阅读进度等），限制 frame 防点击劫持
                    csp = "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https: blob:; font-src 'self'; frame-ancestors 'none'; base-uri 'self'"
                    response['Content-Security-Policy'] = csp
            except Exception:
                pass
        # HTTPS 时添加 HSTS（1 年，含子域名）
        if request.is_secure():
            if 'Strict-Transport-Security' not in response:
                response['Strict-Transport-Security'] = 'max-age=31536000; includeSubDomains'
        return response
