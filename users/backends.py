"""自定义认证后端 - 支持用户名/邮箱登录"""
from django.contrib.auth.backends import ModelBackend
from .models import User


class CustomAuthBackend(ModelBackend):
    def authenticate(self, request, username=None, password=None, **kwargs):
        if username is None or password is None:
            return None
        try:
            user = User.objects.get(username=username)
        except User.DoesNotExist:
            return None
        if user.check_password(password):
            return user
        return None
