from django.urls import path
from . import views

urlpatterns = [
    path('login/', views.login_view, name='login'),
    path('register/', views.register_view, name='register'),
    path('logout/', views.logout_view, name='logout'),
    path('settings/profile/', views.settings_profile, name='settings_profile'),
    path('settings/password/', views.settings_password, name='settings_password'),
    path('settings/account/', views.settings_account, name='settings_account'),
    path('settings/logout-others/', views.settings_logout_others, name='settings_logout_others'),
]
