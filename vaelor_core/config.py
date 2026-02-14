"""
Vaelor Core - 配置抽象层
统一访问站点配置，支持缓存与默认值
"""
from django.core.cache import cache


def get_config(key, default=None):
    """获取配置项，优先从缓存/数据库读取"""
    try:
        from site_app.settings_service import get
        return get(key, default)
    except Exception:
        return default


def set_config(key, value, category='general'):
    """设置配置项并清除缓存"""
    try:
        from site_app.settings_service import set
        set(key, value, category)
    except Exception:
        pass
