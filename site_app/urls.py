from django.urls import path
from . import views

urlpatterns = [
    path('robots.txt', views.robots_txt),
    path('', views.home, name='home'),
    path('search/', views.search, name='search'),
    path('about/', views.about, name='about'),
    path('help/', views.help_page, name='help'),
    path('faq/', views.faq, name='faq'),
    path('p/<slug:slug>/', views.page_detail, name='page_detail'),
    path('games/', views.games, name='games'),
    path('newsletter/subscribe/', views.newsletter_subscribe, name='newsletter_subscribe'),
    path('profile/', views.profile, name='profile'),
    path('admin-panel/', views.admin_panel, name='admin_panel'),
    path('admin-panel/users/', views.admin_users, name='admin_users'),
    path('admin-panel/bbs/', views.admin_bbs, name='admin_bbs'),
    path('admin-panel/bbs/<int:pk>/approve/', views.admin_bbs_approve, name='admin_bbs_approve'),
    path('admin-panel/code/', views.admin_code, name='admin_code'),
    path('admin-panel/logs/', views.admin_logs, name='admin_logs'),
    path('admin-panel/security/', views.admin_security, name='admin_security'),
    path('admin-panel/settings/', views.admin_settings, name='admin_settings'),
    path('admin-panel/rescue/', views.admin_rescue, name='admin_rescue'),
]
