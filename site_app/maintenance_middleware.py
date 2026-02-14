"""维护模式中间件 - 非管理员访问时显示维护页"""
from django.shortcuts import render
from django.http import HttpResponse


def get_maintenance_message():
    try:
        from .models import SiteSetting
        obj = SiteSetting.objects.filter(key='maintenance_mode').first()
        if obj and obj.value == '1':
            msg_obj = SiteSetting.objects.filter(key='maintenance_message').first()
            return msg_obj.value if msg_obj else '站点维护中，请稍后再试。'
    except Exception:
        pass
    return None


class MaintenanceMiddleware:
    def __init__(self, get_response):
        self.get_response = get_response

    def __call__(self, request):
        if request.path.startswith('/install/'):
            return self.get_response(request)
        if request.user.is_authenticated and getattr(request.user, 'role', '') == 'admin':
            return self.get_response(request)
        msg = get_maintenance_message()
        if msg:
            return render(request, 'maintenance.html', {'message': msg})
        return self.get_response(request)
