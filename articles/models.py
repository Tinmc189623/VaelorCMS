"""文章 - 对应原 articles 表"""
from django.db import models
from django.conf import settings


class Article(models.Model):
    STATUS_CHOICES = [('draft', '草稿'), ('published', '已发布')]
    user = models.ForeignKey(settings.AUTH_USER_MODEL, on_delete=models.SET_NULL, null=True, related_name='articles')
    title = models.CharField(max_length=255)
    content = models.TextField()
    status = models.CharField(max_length=16, default='draft', choices=STATUS_CHOICES)
    tags = models.CharField(max_length=255, default='', blank=True)
    category = models.CharField(max_length=128, default='', blank=True)
    scheduled_at = models.DateTimeField(null=True, blank=True)
    view_count = models.PositiveIntegerField(default=0, db_index=True)
    like_count = models.PositiveIntegerField(default=0, db_index=True)
    is_pinned = models.BooleanField(default=False, db_index=True)
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True, null=True)

    class Meta:
        db_table = 'articles'
        ordering = ['-is_pinned', '-id']


class ArticleLike(models.Model):
    """文章点赞"""
    article = models.ForeignKey(Article, on_delete=models.CASCADE, related_name='likes')
    user = models.ForeignKey(settings.AUTH_USER_MODEL, on_delete=models.CASCADE)
    created_at = models.DateTimeField(auto_now_add=True)

    class Meta:
        db_table = 'article_likes'
        unique_together = [('article', 'user')]


class ArticleFavorite(models.Model):
    """文章收藏"""
    article = models.ForeignKey(Article, on_delete=models.CASCADE, related_name='favorites')
    user = models.ForeignKey(settings.AUTH_USER_MODEL, on_delete=models.CASCADE)
    created_at = models.DateTimeField(auto_now_add=True)

    class Meta:
        db_table = 'article_favorites'
        unique_together = [('article', 'user')]


class ArticleComment(models.Model):
    """文章评论"""
    article = models.ForeignKey(Article, on_delete=models.CASCADE, related_name='comments')
    user = models.ForeignKey(settings.AUTH_USER_MODEL, on_delete=models.CASCADE, null=True, blank=True)
    author_name = models.CharField(max_length=64, default='', blank=True)  # 游客昵称
    content = models.TextField()
    parent = models.ForeignKey('self', on_delete=models.CASCADE, null=True, blank=True, related_name='replies')
    is_approved = models.BooleanField(default=True, db_index=True)
    created_at = models.DateTimeField(auto_now_add=True)

    class Meta:
        db_table = 'article_comments'
        ordering = ['id']
