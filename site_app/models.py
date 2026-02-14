"""站点模型 - admin_logs、站点设置"""
from django.db import models
from django.conf import settings


class SiteSetting(models.Model):
    """管理员可编辑的站点设置（key-value）"""
    key = models.CharField(max_length=64, unique=True, db_index=True)
    value = models.TextField(default='')
    category = models.CharField(max_length=32, default='general')  # general, security, user, content, maintenance
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        db_table = 'site_settings'
        verbose_name = '站点设置'
        verbose_name_plural = '站点设置'


class AdminLog(models.Model):
    user = models.ForeignKey(settings.AUTH_USER_MODEL, on_delete=models.SET_NULL, null=True)
    username = models.CharField(max_length=64, default='')
    action = models.CharField(max_length=32)
    target = models.CharField(max_length=255, default='')
    target_id = models.PositiveIntegerField(default=0)
    ip_address = models.CharField(max_length=45, default='')
    created_at = models.DateTimeField(auto_now_add=True)

    class Meta:
        db_table = 'admin_logs'
        ordering = ['-id']


class NewsletterSubscriber(models.Model):
    """邮件订阅"""
    email = models.EmailField(unique=True)
    is_active = models.BooleanField(default=True, db_index=True)
    created_at = models.DateTimeField(auto_now_add=True)

    class Meta:
        db_table = 'newsletter_subscribers'
        verbose_name = '邮件订阅'
        verbose_name_plural = '邮件订阅'


class Link(models.Model):
    """友情链接"""
    title = models.CharField(max_length=128)
    url = models.URLField(max_length=512)
    description = models.CharField(max_length=255, default='', blank=True)
    order = models.PositiveSmallIntegerField(default=0)
    is_visible = models.BooleanField(default=True, db_index=True)
    created_at = models.DateTimeField(auto_now_add=True)

    class Meta:
        db_table = 'links'
        ordering = ['order', 'id']
        verbose_name = '友情链接'
        verbose_name_plural = '友情链接'


class Page(models.Model):
    """自定义页面 - 管理员创建，通过 /p/<slug>/ 访问"""
    slug = models.SlugField(max_length=128, unique=True, db_index=True)
    title = models.CharField(max_length=255)
    content = models.TextField()
    content_is_html = models.BooleanField(default=False, help_text='内容为 HTML 时勾选，否则按换行显示')
    is_published = models.BooleanField(default=True)
    show_in_nav = models.BooleanField(default=False, help_text='是否在导航栏显示')
    order = models.PositiveSmallIntegerField(default=0)
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        db_table = 'pages'
        ordering = ['order', 'slug']
        verbose_name = '自定义页面'
        verbose_name_plural = '自定义页面'
