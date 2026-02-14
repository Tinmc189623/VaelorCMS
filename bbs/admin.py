from django.contrib import admin
from .models import BbsPost


@admin.register(BbsPost)
class BbsPostAdmin(admin.ModelAdmin):
    list_display = ('id', 'title', 'user', 'created_at')
    search_fields = ('title', 'content')
