"""
站点设置服务 - 读写 SiteSetting，提供默认值
"""
from django.conf import settings as django_settings

DEFAULTS = {
    # general
    'site_name': ('VaelorCMS', 'general'),
    'site_description': ('', 'general'),
    'site_keywords': ('', 'general'),
    'timezone': ('Asia/Shanghai', 'general'),
    'language': ('zh-hans', 'general'),
    # security
    'login_max_attempts': ('5', 'security'),
    'login_lockout_minutes': ('15', 'security'),
    'api_rate_limit': ('60', 'security'),  # 每分钟每 IP 最大请求数，0 表示不限
    'min_password_length': ('8', 'security'),
    'require_strong_password': ('0', 'security'),  # 1=字母+数字
    'require_password_special': ('0', 'security'),  # 1=须含特殊字符 !@#$%^&* 等
    'session_timeout_days': ('7', 'security'),
    'force_https': ('0', 'security'),
    # user
    'allow_register': ('1', 'user'),
    'require_email_verify': ('0', 'user'),
    'default_user_role': ('user', 'user'),
    # content
    'bbs_moderation': ('0', 'content'),
    'bbs_guest_post': ('0', 'content'),
    'article_comment': ('1', 'content'),
    # maintenance
    'maintenance_mode': ('0', 'maintenance'),
    'maintenance_message': ('站点维护中，请稍后再试。', 'maintenance'),
    # theme & plugin
    'theme': ('default', 'theme'),
    'plugins_enabled': ('', 'theme'),  # 逗号分隔插件 ID
    'aero_enabled': ('1', 'theme'),  # 1=启用站点全局 Aero 玻璃特效
    'aero_blur': ('16', 'theme'),  # 玻璃模糊强度 px
    'accent_color': ('', 'theme'),  # 强调色覆盖，如 #8b5a2b，空则用主题默认
    'custom_css': ('', 'theme'),  # 自定义 CSS 片段，注入到 head
    'dark_mode': ('0', 'theme'),  # 1=暗色模式，0=跟随系统，2=亮色
    # seo
    'seo_title_suffix': (' - {site_name}', 'seo'),
    'seo_default_description': ('', 'seo'),
    'seo_default_keywords': ('', 'seo'),
    'seo_og_image': ('', 'seo'),
    'seo_twitter_card': ('summary', 'seo'),
    'seo_canonical_base': ('', 'seo'),  # 如 https://example.com
}


def get(key, default=None):
    from .models import SiteSetting
    try:
        obj = SiteSetting.objects.get(key=key)
        return obj.value
    except SiteSetting.DoesNotExist:
        if default is not None:
            return default
        return DEFAULTS.get(key, ('', 'general'))[0]


def set(key, value, category='general'):
    from .models import SiteSetting
    from django.core.cache import cache
    SiteSetting.objects.update_or_create(key=key, defaults={'value': str(value), 'category': category})
    cache.delete('vaelor_site_settings_ctx')


def get_all_by_category():
    try:
        from .models import SiteSetting
        vals = dict(SiteSetting.objects.values_list('key', 'value'))
    except Exception:
        vals = {}
    result = {}
    for key, (default, cat) in DEFAULTS.items():
        if cat not in result:
            result[cat] = {}
        result[cat][key] = vals.get(key, default)
    return result


def seed_defaults():
    """首次使用时写入默认值"""
    from .models import SiteSetting
    for key, (default, category) in DEFAULTS.items():
        SiteSetting.objects.get_or_create(key=key, defaults={'value': default, 'category': category})


def restore_defaults(category=None):
    """
    一键恢复默认配置
    category: 指定分类则只恢复该分类，None 则恢复全部
    返回恢复的键数量
    """
    from .models import SiteSetting
    count = 0
    for key, (default, cat) in DEFAULTS.items():
        if category is None or cat == category:
            SiteSetting.objects.update_or_create(key=key, defaults={'value': default, 'category': cat})
            count += 1
    from django.core.cache import cache
    cache.delete('vaelor_site_settings_ctx')
    return count


def export_backup():
    """导出当前配置为 JSON 字符串，用于备份"""
    import json
    try:
        from .models import SiteSetting
        rows = list(SiteSetting.objects.values_list('key', 'value', 'category'))
        return json.dumps([{'key': k, 'value': v, 'category': c} for k, v, c in rows], ensure_ascii=False, indent=2)
    except Exception:
        return '[]'


def import_backup(json_str):
    """
    从 JSON 备份恢复配置
    返回 (成功数, 失败数)
    """
    import json
    from .models import SiteSetting
    from django.core.cache import cache
    ok, fail = 0, 0
    try:
        data = json.loads(json_str)
        if not isinstance(data, list):
            return 0, 1
        for item in data:
            if isinstance(item, dict) and 'key' in item:
                try:
                    SiteSetting.objects.update_or_create(
                        key=str(item['key'])[:64],
                        defaults={'value': str(item.get('value', '')), 'category': str(item.get('category', 'general'))[:32]}
                    )
                    ok += 1
                except Exception:
                    fail += 1
        cache.delete('vaelor_site_settings_ctx')
    except (json.JSONDecodeError, TypeError):
        return 0, 1
    return ok, fail
