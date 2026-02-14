from django.contrib import admin
from .models import CodeSnippet


@admin.register(CodeSnippet)
class CodeSnippetAdmin(admin.ModelAdmin):
    list_display = ('id', 'title', 'language', 'user', 'created_at')
    search_fields = ('title', 'code')
