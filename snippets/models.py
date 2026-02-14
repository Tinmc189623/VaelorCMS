"""代码片段 - 对应原 code_snippets 表"""
from django.db import models
from django.conf import settings


class CodeSnippet(models.Model):
    user = models.ForeignKey(settings.AUTH_USER_MODEL, on_delete=models.SET_NULL, null=True, related_name='code_snippets')
    title = models.CharField(max_length=255)
    code = models.TextField()
    language = models.CharField(max_length=32, default='text')
    created_at = models.DateTimeField(auto_now_add=True)

    class Meta:
        db_table = 'code_snippets'
        ordering = ['-id']
