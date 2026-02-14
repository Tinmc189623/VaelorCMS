from django.contrib import admin
from django.contrib.auth.admin import UserAdmin as BaseUserAdmin
from .models import User


@admin.register(User)
class UserAdmin(BaseUserAdmin):
    list_display = ('id', 'username', 'email', 'role', 'created_at')
    list_filter = ('role',)
    search_fields = ('username', 'email')
    ordering = ('-id',)
    fieldsets = (
        (None, {'fields': ('username', 'password')}),
        ('信息', {'fields': ('email', 'role')}),
    )
    add_fieldsets = (
        (None, {'fields': ('username', 'password1', 'password2')}),
        ('信息', {'fields': ('email', 'role')}),
    )
