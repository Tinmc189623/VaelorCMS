"""
安装重定向中间件：未安装时将所有请求重定向到安装向导
"""
from django.http import HttpResponseRedirect
from django.urls import reverse
def is_installed():
    from django.conf import settings
    lock_path = getattr(settings, 'INSTALLED_LOCK_PATH', None)
    if not lock_path:
        return True  # 无锁路径则视为已安装
    return lock_path.exists()


class InstallRedirectMiddleware:
    """未安装时重定向到 /install/"""

    def __init__(self, get_response):
        self.get_response = get_response

    def __call__(self, request):
        if is_installed():
            # 已安装：禁止访问安装页
            if request.path.startswith('/install'):
                from django.http import HttpResponseNotFound
                return HttpResponseNotFound('安装已完成')
            return self.get_response(request)

        # 未安装：仅允许访问安装相关路径和静态资源
        if request.path.startswith('/install') or request.path.startswith('/static'):
            return self.get_response(request)
        return HttpResponseRedirect(reverse('install_step1'))
