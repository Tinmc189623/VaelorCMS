from django.contrib import admin
from django.utils.html import format_html
from django.core.cache import cache
from .models import Page, SiteSetting, Link, NewsletterSubscriber


def _invalidate_site_settings_cache():
    cache.delete('vaelor_site_settings_ctx')


@admin.register(Page)
class PageAdmin(admin.ModelAdmin):
    list_display = ('slug', 'title', 'content_is_html', 'is_published', 'show_in_nav', 'order', 'updated_at', 'view_link')
    list_editable = ('is_published', 'show_in_nav', 'order')
    search_fields = ('slug', 'title', 'content')
    prepopulated_fields = {'slug': ('title',)}

    def save_model(self, request, obj, form, change):
        super().save_model(request, obj, form, change)
        _invalidate_site_settings_cache()

    def delete_model(self, request, obj):
        super().delete_model(request, obj)
        _invalidate_site_settings_cache()

    def view_link(self, obj):
        from django.urls import reverse
        url = reverse('page_detail', args=[obj.slug])
        return format_html('<a href="{}" target="_blank">预览</a>', url)
    view_link.short_description = '预览'


@admin.register(SiteSetting)
class SiteSettingAdmin(admin.ModelAdmin):
    list_display = ('key', 'value_short', 'category', 'updated_at')
    list_filter = ('category',)
    search_fields = ('key', 'value')
    ordering = ('category', 'key')

    def value_short(self, obj):
        return (obj.value[:50] + '...') if len(obj.value) > 50 else obj.value
    value_short.short_description = '值'


@admin.register(NewsletterSubscriber)
class NewsletterSubscriberAdmin(admin.ModelAdmin):
    list_display = ('email', 'is_active', 'created_at')
    list_filter = ('is_active',)
    search_fields = ('email',)


@admin.register(Link)
class LinkAdmin(admin.ModelAdmin):
    list_display = ('title', 'url', 'order', 'is_visible', 'created_at')
    list_editable = ('order', 'is_visible')
    list_filter = ('is_visible',)
    search_fields = ('title', 'url')

    def save_model(self, request, obj, form, change):
        super().save_model(request, obj, form, change)
        _invalidate_site_settings_cache()

    def delete_model(self, request, obj):
        super().delete_model(request, obj)
        _invalidate_site_settings_cache()
