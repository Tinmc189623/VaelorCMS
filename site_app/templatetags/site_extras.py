"""站点模板标签 - SEO、主题、插件、阅读体验"""
from django import template
from django.utils.safestring import mark_safe

register = template.Library()


@register.filter
def split(value, arg=','):
    """按分隔符拆分字符串为列表"""
    if not value:
        return []
    return [s.strip() for s in str(value).replace('，', arg).split(arg) if s.strip()]


@register.filter
def reading_time(text, wpm=200):
    """估算阅读时间（分钟），按 200 字/分钟"""
    if not text:
        return 0
    length = len(str(text).strip())
    return max(1, round(length / wpm))


@register.simple_tag(takes_context=True)
def page_title(context, default='首页'):
    """当前页标题（不含站点后缀），供 title 与 og:title 复用"""
    return context.get('page_title', default)


@register.simple_tag(takes_context=True)
def page_description(context):
    """当前页描述，供 meta description 与 og:description"""
    return context.get('page_description', '')


@register.simple_tag(takes_context=True)
def page_keywords(context):
    """当前页关键词"""
    return context.get('page_keywords', '')


@register.simple_tag(takes_context=True)
def canonical_url(context):
    """规范 URL"""
    base = context.get('seo_canonical_base', '')
    path = context.get('request')
    if base and path:
        return base.rstrip('/') + (path.path or '/')
    return ''


@register.simple_tag(takes_context=True)
def breadcrumb(context, items):
    """面包屑：items 为 [(url, label), ...]"""
    from django.utils.safestring import mark_safe
    parts = []
    for i, (url, label) in enumerate(items):
        if url and i < len(items) - 1:
            parts.append(f'<a href="{url}">{label}</a>')
        else:
            parts.append(f'<span>{label}</span>')
    return mark_safe(' &rsaquo; '.join(parts))


@register.simple_tag
def run_plugin_hook(hook_name):
    """执行插件钩子并返回拼接的 HTML（在模板中调用）"""
    from vaelor.plugin import run_hook
    from django.template.context import RequestContext
    # 需要在模板中传入 request，这里简化处理
    return mark_safe(''.join(run_hook(hook_name, None, {})))
