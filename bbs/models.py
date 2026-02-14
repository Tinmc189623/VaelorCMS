"""论坛帖子 - 对应原 bbs_posts 表"""
from django.db import models
from django.conf import settings


class BbsPost(models.Model):
    user = models.ForeignKey(settings.AUTH_USER_MODEL, on_delete=models.SET_NULL, null=True, related_name='bbs_posts')
    title = models.CharField(max_length=255)
    content = models.TextField()
    approved = models.BooleanField(default=True, db_index=True)  # 审核通过才显示，管理员可改
    created_at = models.DateTimeField(auto_now_add=True)

    class Meta:
        db_table = 'bbs_posts'
        ordering = ['-id']
