from django.apps import AppConfig


class InstallAppConfig(AppConfig):
    default_auto_field = 'django.db.models.BigAutoField'
    name = 'install_app'
    verbose_name = '安装向导'
