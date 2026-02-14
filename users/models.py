"""
用户模型 - 匹配原有 users 表结构
"""
from django.db import models
from django.contrib.auth.models import AbstractBaseUser, BaseUserManager, PermissionsMixin


class UserManager(BaseUserManager):
    def create_user(self, username, password=None, email='', **kwargs):
        user = self.model(username=username, email=email or '', **kwargs)
        user.set_password(password)
        user.save(using=self._db)
        return user

    def create_superuser(self, username, password=None, email='', **kwargs):
        kwargs.setdefault('role', 'admin')
        return self.create_user(username, password, email, **kwargs)


class User(AbstractBaseUser, PermissionsMixin):
    """自定义用户 - 对应原 users 表"""
    username = models.CharField(max_length=64, unique=True)
    email = models.CharField(max_length=128, default='', blank=True)
    role = models.CharField(max_length=20, default='user')
    created_at = models.DateTimeField(auto_now_add=True)
    last_login = models.DateTimeField(null=True, blank=True)  # AbstractBaseUser 需要

    objects = UserManager()
    USERNAME_FIELD = 'username'
    REQUIRED_FIELDS = []

    class Meta:
        db_table = 'users'
        managed = True  # 首次运行 migrate 建表；已有库需先执行 sql/alter_users_lastlogin.sql

    @property
    def is_staff(self):
        return self.role == 'admin'

    @property
    def is_superuser(self):
        return self.role == 'admin'

    def has_perm(self, perm, obj=None):
        return self.is_superuser

    def has_module_perms(self, app_label):
        return self.is_superuser


class UserProfile(models.Model):
    """用户扩展资料与设置"""
    user = models.OneToOneField(User, on_delete=models.CASCADE, related_name='profile')
    nickname = models.CharField(max_length=64, default='', blank=True)
    bio = models.TextField(default='', blank=True)
    avatar = models.ImageField(upload_to='avatars/%Y/%m/', blank=True, null=True)
    # 隐私
    show_email = models.BooleanField(default=False, verbose_name='公开邮箱')
    profile_visible = models.CharField(max_length=20, default='all', choices=[
        ('all', '所有人'),
        ('login', '仅登录用户'),
        ('self', '仅自己'),
    ], verbose_name='资料可见性')
    # 通知
    email_on_reply = models.BooleanField(default=True, verbose_name='回复时邮件通知')
    email_on_mention = models.BooleanField(default=True, verbose_name='被@时邮件通知')
    # 偏好
    timezone = models.CharField(max_length=64, default='Asia/Shanghai')
    language = models.CharField(max_length=16, default='zh-hans')
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        db_table = 'user_profiles'
        verbose_name = '用户资料'
        verbose_name_plural = '用户资料'
