from django.contrib import admin
from .models import Article, ArticleComment


@admin.register(Article)
class ArticleAdmin(admin.ModelAdmin):
    list_display = ('id', 'title', 'user', 'status', 'category', 'view_count', 'created_at')
    list_filter = ('status',)
    search_fields = ('title', 'content')


@admin.register(ArticleComment)
class ArticleCommentAdmin(admin.ModelAdmin):
    list_display = ('id', 'article', 'author_display', 'content_short', 'is_approved', 'created_at')
    list_filter = ('is_approved',)
    list_editable = ('is_approved',)
    search_fields = ('content', 'author_name')
    raw_id_fields = ('article', 'user', 'parent')

    def author_display(self, obj):
        return obj.user.username if obj.user else (obj.author_name or '游客')
    author_display.short_description = '作者'

    def content_short(self, obj):
        return (obj.content[:40] + '...') if len(obj.content) > 40 else obj.content
    content_short.short_description = '内容'
