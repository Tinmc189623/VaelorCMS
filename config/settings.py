"""
Django 配置 - VaelorCMS
基于 config/config.ini 与原有站点信息
使用相对路径，适配任意部署环境
"""
import os
from pathlib import Path

# 项目根目录（相对 config/settings.py 的上级）
BASE_DIR = Path(__file__).resolve().parent.parent
INSTALLED_LOCK_PATH = BASE_DIR / 'config' / 'installed.lock'

# 读取 config.ini
_config = {}
_config_path = BASE_DIR / 'config' / 'config.ini'
if _config_path.exists():
    import configparser
    cp = configparser.ConfigParser()
    cp.read(_config_path, encoding='utf-8')
    for section in cp.sections():
        _config[section] = dict(cp[section])

def _get(section, key, default=''):
    return _config.get(section, {}).get(key, default)

SECRET_KEY = os.environ.get('DJANGO_SECRET_KEY', 'dev-secret-change-in-production')
# 默认生产模式，本地开发可设 DJANGO_DEBUG=1
DEBUG = os.environ.get('DJANGO_DEBUG', '0') == '1'
ALLOWED_HOSTS = ['*']

# CSRF 可信来源（反向代理/HTTPS 部署时必需）
_csrf_origins = os.environ.get('CSRF_TRUSTED_ORIGINS', '')
CSRF_TRUSTED_ORIGINS = [o.strip() for o in _csrf_origins.split(',') if o.strip()]
# 从 config.ini [site] trusted_origins 补充（逗号分隔）
_extra = _get('site', 'trusted_origins', '')
if _extra:
    CSRF_TRUSTED_ORIGINS.extend([o.strip() for o in _extra.split(',') if o.strip()])

INSTALLED_APPS = [
    'users.apps.UsersConfig',  # 必须在 auth 之前，AUTH_USER_MODEL 依赖
    'django.contrib.admin',
    'django.contrib.auth',
    'django.contrib.contenttypes',
    'django.contrib.sessions',
    'django.contrib.messages',
    'django.contrib.staticfiles',
    'django.contrib.sitemaps',
    'bbs',
    'snippets',
    'articles',
    'site_app',
    'install_app',
]

MIDDLEWARE = [
    'django.middleware.security.SecurityMiddleware',
    'whitenoise.middleware.WhiteNoiseMiddleware',
    'django.contrib.sessions.middleware.SessionMiddleware',
    'install_app.middleware.InstallRedirectMiddleware',
    'site_app.csrf_proxy_middleware.CsrfProxyTrustMiddleware',
    'site_app.api_throttle.api_throttle_middleware',
    'django.middleware.common.CommonMiddleware',
    'django.middleware.csrf.CsrfViewMiddleware',
    'django.contrib.auth.middleware.AuthenticationMiddleware',
    'site_app.maintenance_middleware.MaintenanceMiddleware',
    'site_app.seo_middleware.SEOMiddleware',
    'django.contrib.messages.middleware.MessageMiddleware',
    'django.middleware.clickjacking.XFrameOptionsMiddleware',
    'site_app.security_headers_middleware.SecurityHeadersMiddleware',
]

ROOT_URLCONF = 'config.urls'
WSGI_APPLICATION = 'config.wsgi.application'

# 数据库 - 环境变量可覆盖（容器部署）；config.ini 存在用 MySQL，否则 SQLite
# 容器内 127.0.0.1 指向自身无法连 MySQL，未显式设置 DB_HOST 时自动用 SQLite 进入安装向导
_db_host = os.environ.get('DB_HOST') or (_get('database', 'host', '127.0.0.1') if _config_path.exists() else '')
_in_container = bool(
    Path('/.dockerenv').exists()
    or os.environ.get('KUBERNETES_SERVICE_HOST')
    or os.environ.get('KUBERNETES_SERVICE_PORT')
    or os.environ.get('DOCKER')
    or os.environ.get('container')
)
_use_sqlite = (not _config_path.exists() or
               (_in_container and _db_host in ('127.0.0.1', 'localhost') and 'DB_HOST' not in os.environ))

if not _use_sqlite and _config_path.exists():
    DATABASES = {
        'default': {
            'ENGINE': 'django.db.backends.mysql',
            'HOST': _db_host or _get('database', 'host', '127.0.0.1'),
            'PORT': os.environ.get('DB_PORT', _get('database', 'port', '3306')),
            'NAME': os.environ.get('DB_NAME', _get('database', 'dbname', 'vaelor_cms')),
            'USER': os.environ.get('DB_USER', _get('database', 'username', 'vaelor_user')),
            'PASSWORD': os.environ.get('DB_PASSWORD', _get('database', 'password', '')),
            'CONN_MAX_AGE': 60,
            'OPTIONS': {
                'charset': _get('database', 'charset', 'utf8mb4'),
                'init_command': "SET sql_mode='STRICT_TRANS_TABLES'",
                'connect_timeout': 10,
            },
        }
    }
else:
    DATABASES = {
        'default': {
            'ENGINE': 'django.db.backends.sqlite3',
            'NAME': BASE_DIR / 'config' / 'install_temp.db',
        }
    }

# 自定义用户模型（匹配原有 users 表）
AUTH_USER_MODEL = 'users.User'
AUTHENTICATION_BACKENDS = ['users.backends.CustomAuthBackend']

# 密码
AUTH_PASSWORD_VALIDATORS = [
    {'NAME': 'django.contrib.auth.password_validation.MinimumLengthValidator', 'OPTIONS': {'min_length': 6}},
]

# 国际化
LANGUAGE_CODE = 'zh-hans'
TIME_ZONE = 'Asia/Shanghai'
USE_I18N = True
USE_TZ = True

# 静态文件
STATIC_URL = '/static/'
STATICFILES_DIRS = [
    BASE_DIR / 'assets',
    ('themes', BASE_DIR / 'themes'),  # themes/<id>/style.css -> /static/themes/<id>/style.css
]
STATIC_ROOT = BASE_DIR / 'staticfiles'

# 媒体文件
MEDIA_URL = '/media/'
MEDIA_ROOT = BASE_DIR / 'storage' / 'uploads'

# 模板
TEMPLATES = [
    {
        'BACKEND': 'django.template.backends.django.DjangoTemplates',
        'DIRS': [BASE_DIR / 'templates'],
        'APP_DIRS': True,
        'OPTIONS': {
            'context_processors': [
                'django.template.context_processors.request',
                'django.contrib.auth.context_processors.auth',
                'django.contrib.messages.context_processors.messages',
                'site_app.context_processors.site_settings',
                'site_app.context_processors.plugin_hooks',
            ],
        },
    },
]

# CMS 名称
CMS_NAME = 'VaelorCMS'

# 站点信息（保留原有配置）
SITE_NAME = _get('site', 'name', CMS_NAME)
SESSION_COOKIE_NAME = _get('site', 'session_name', 'VAELOR_SESS')
SESSION_COOKIE_PATH = _get('site', 'cookie_path', '/')

# 默认主键
DEFAULT_AUTO_FIELD = 'django.db.models.BigAutoField'

# 缓存 - 从 config.ini [cache] 读取；driver=redis 时用 Redis，否则 LocMem
_cache_driver = (_get('cache', 'driver') or '').strip().lower()
if _config_path.exists() and _cache_driver == 'redis':
    _redis_host = os.environ.get('REDIS_HOST', _get('cache', 'host', '127.0.0.1'))
    _redis_port = os.environ.get('REDIS_PORT', _get('cache', 'port', '6379'))
    _redis_pwd = os.environ.get('REDIS_PASSWORD', _get('cache', 'password', ''))
    _redis_prefix = _get('cache', 'prefix', 'vaelor_')
    _redis_loc = f'redis://:{_redis_pwd}@{_redis_host}:{_redis_port}/1' if _redis_pwd else f'redis://{_redis_host}:{_redis_port}/1'
    CACHES = {
        'default': {
            'BACKEND': 'django.core.cache.backends.redis.RedisCache',
            'LOCATION': _redis_loc,
            'KEY_PREFIX': _redis_prefix,
        }
    }
else:
    CACHES = {
        'default': {
            'BACKEND': 'django.core.cache.backends.locmem.LocMemCache',
            'OPTIONS': {'MAX_ENTRIES': 256},
        }
    }

# 会话安全
SESSION_COOKIE_HTTPONLY = True
SESSION_COOKIE_SAMESITE = 'Lax'
SESSION_SAVE_EVERY_REQUEST = False  # 减少写库，按需可改为 True
# HTTPS 下启用 Secure：环境变量或反向代理（X-Forwarded-Proto: https）时自动启用
_use_secure = os.environ.get('DJANGO_HTTPS', '').lower() in ('1', 'true', 'yes') or not DEBUG
if _use_secure:
    SESSION_COOKIE_SECURE = True

# CSRF Cookie 安全
CSRF_COOKIE_HTTPONLY = True
CSRF_COOKIE_SAMESITE = 'Lax'
if _use_secure:
    CSRF_COOKIE_SECURE = True

# 反向代理下自动补充 CSRF 可信来源（HTTPS 时根据 Host 推断），设 0 可关闭
CSRF_TRUSTED_ORIGINS_AUTO = os.environ.get('CSRF_TRUSTED_ORIGINS_AUTO', '1') == '1'

# 反向代理 HTTPS（Nginx 等设置 X-Forwarded-Proto: https 时）
SECURE_PROXY_SSL_HEADER = ('HTTP_X_FORWARDED_PROTO', 'https')

# 生产安全（DEBUG=False 时自动启用）
if not DEBUG:
    SECURE_CONTENT_TYPE_NOSNIFF = True
    X_FRAME_OPTIONS = 'DENY'
    # 禁止弱 SECRET_KEY（仅已安装站点时提醒，安装向导阶段不打扰）
    if INSTALLED_LOCK_PATH.exists() and SECRET_KEY in ('dev-secret-change-in-production', 'changeme', 'secret', ''):
        import warnings
        warnings.warn('SECRET_KEY 过于简单，生产环境请设置 DJANGO_SECRET_KEY 环境变量')
