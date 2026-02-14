"""VaelorCMS 主路由"""
from django.contrib import admin

admin.site.site_header = 'VaelorCMS 管理'
admin.site.site_title = 'VaelorCMS 管理'
admin.site.index_title = '站点管理'
from django.contrib.sitemaps.views import sitemap
from django.urls import path, include
from django.conf import settings
from django.conf.urls.static import static
from site_app.sitemaps import StaticSitemap, ArticleSitemap, BbsSitemap, CodeSitemap, PageSitemap

sitemaps = {
    'static': StaticSitemap,
    'articles': ArticleSitemap,
    'bbs': BbsSitemap,
    'code': CodeSitemap,
    'pages': PageSitemap,
}

handler404 = 'site_app.views.page_404'
handler500 = 'site_app.views.page_500'

urlpatterns = [
    path('sitemap.xml', sitemap, {'sitemaps': sitemaps}, name='django.contrib.sitemaps.views.sitemap'),
    path('install/', include('install_app.urls')),
    path('admin/', admin.site.urls),
    path('api/v1/', include('site_app.api_urls')),
    path('', include('site_app.urls')),
    path('bbs/', include('bbs.urls')),
    path('code/', include('snippets.urls')),
    path('articles/', include('articles.urls')),
    path('users/', include('users.urls')),
]

if settings.DEBUG:
    urlpatterns += static(settings.MEDIA_URL, document_root=settings.MEDIA_ROOT)
    urlpatterns += static(settings.STATIC_URL, document_root=settings.STATICFILES_DIRS[0])
