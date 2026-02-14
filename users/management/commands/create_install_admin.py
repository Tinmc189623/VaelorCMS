"""
安装向导用：根据环境变量创建管理员账号
"""
import os
from django.core.management.base import BaseCommand
from users.models import User


class Command(BaseCommand):
    help = '安装向导：根据 INSTALL_ADMIN_USER、INSTALL_ADMIN_PASS 创建管理员'

    def handle(self, *args, **options):
        username = os.environ.get('INSTALL_ADMIN_USER', '').strip()
        password = os.environ.get('INSTALL_ADMIN_PASS', '')
        if not username or not password:
            self.stderr.write('缺少 INSTALL_ADMIN_USER 或 INSTALL_ADMIN_PASS')
            return
        if User.objects.filter(username=username).exists():
            user = User.objects.get(username=username)
            user.set_password(password)
            user.role = 'admin'
            user.save()
            self.stdout.write(f'已更新管理员 {username}')
        else:
            User.objects.create_superuser(username, password, '')
            self.stdout.write(f'已创建管理员 {username}')
