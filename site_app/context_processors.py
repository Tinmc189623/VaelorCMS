from django.conf import settings
from django.core.cache import cache

try:
    from config.version import __version__ as CMS_VERSION
except ImportError:
    CMS_VERSION = 'Demo-26.02.13.26'

_CTX_KEYS = ('site_name', 'site_description', 'site_keywords', 'theme', 'seo_default_description',
             'seo_canonical_base', 'seo_og_image', 'seo_default_keywords',
             'aero_enabled', 'aero_blur', 'accent_color', 'custom_css', 'dark_mode')
_SITE_SETTINGS_CACHE_KEY = 'vaelor_site_settings_ctx'
_SITE_SETTINGS_CACHE_TIMEOUT = 60


def _get_settings_bulk(keys):
    """批量读取设置，减少 DB 查询"""
    try:
        from .models import SiteSetting
        rows = SiteSetting.objects.filter(key__in=keys).values_list('key', 'value')
        return dict(rows)
    except Exception:
        return {}


def site_settings(request):
    cached = cache.get(_SITE_SETTINGS_CACHE_KEY)
    if cached is not None:
        return cached
    vals = _get_settings_bulk(_CTX_KEYS)
    site_name = (vals.get('site_name') or '').strip() or getattr(settings, 'SITE_NAME', 'VaelorCMS')
    theme_id = (vals.get('theme') or 'default').strip() or 'default'
    base = (vals.get('seo_canonical_base', '') or '').rstrip('/')
    desc = vals.get('site_description', '') or vals.get('seo_default_description', '')
    nav_pages = []
    footer_links = []
    try:
        from .models import Page, Link
        nav_pages = list(Page.objects.filter(show_in_nav=True, is_published=True).order_by('order', 'slug').values('slug', 'title'))
        footer_links = list(Link.objects.filter(is_visible=True).order_by('order', 'id').values('title', 'url'))
    except Exception:
        pass
    aero_enabled = (vals.get('aero_enabled') or '1').strip() == '1'
    aero_blur = (vals.get('aero_blur') or '16').strip()
    _accent = (vals.get('accent_color') or '').strip()
    # 仅允许 #xxx 或 #xxxxxx 格式
    accent_color = _accent if _accent and _accent.startswith('#') and len(_accent) in (4, 7) and all(c in '0123456789abcdefABCDEF#' for c in _accent) else ''
    _raw_css = (vals.get('custom_css') or '').strip()
    # 移除危险 CSS（expression、javascript: 等）
    if _raw_css:
        import re
        _raw_css = re.sub(r'expression\s*\(', '', _raw_css, flags=re.I)
        _raw_css = re.sub(r'javascript\s*:', '', _raw_css, flags=re.I)
        _raw_css = re.sub(r'-moz-binding\s*:', '', _raw_css, flags=re.I)
        _raw_css = re.sub(r'behavior\s*:\s*url\s*\(', '', _raw_css, flags=re.I)
    custom_css = _raw_css
    dark_mode = (vals.get('dark_mode') or '0').strip()
    result = {
        'site_name': site_name,
        'nav_pages': nav_pages,
        'footer_links': footer_links,
        'site_description': vals.get('site_description', ''),
        'site_keywords': vals.get('site_keywords', ''),
        'cms_version': CMS_VERSION,
        'theme_id': theme_id,
        'theme_css': f'themes/{theme_id}/style.css' if theme_id and theme_id != 'default' else None,
        'seo_description': desc,
        'seo_keywords': vals.get('site_keywords', '') or vals.get('seo_default_keywords', ''),
        'seo_canonical_base': base,
        'seo_og_image': vals.get('seo_og_image', ''),
        'aero_enabled': aero_enabled,
        'aero_blur': aero_blur,
        'accent_color': accent_color,
        'custom_css': custom_css,
        'dark_mode': dark_mode,
    }
    cache.set(_SITE_SETTINGS_CACHE_KEY, result, _SITE_SETTINGS_CACHE_TIMEOUT)
    return result




def plugin_hooks(request):
    """插件钩子输出，供模板使用"""
    try:
        from vaelor.plugin import run_hook
        return {
            'head_hooks': run_hook('head_extra', request, {}),
            'body_top_hooks': run_hook('body_top', request, {}),
            'nav_hooks': run_hook('nav_extra', request, {}),
            'footer_hooks': run_hook('footer_extra', request, {}),
            'body_bottom_hooks': run_hook('body_bottom', request, {}),
        }
    except Exception:
        return {
            'head_hooks': [], 'body_top_hooks': [], 'nav_hooks': [],
            'footer_hooks': [], 'body_bottom_hooks': [],
        }


