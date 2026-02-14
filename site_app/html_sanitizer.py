"""
HTML 净化 - 自定义页面 content_is_html 时防 XSS
优先使用 bleach，回退到自研正则净化
"""
import re
import html

# 危险标签（含内容） - script, style, svg 等
DANGEROUS_BLOCK = re.compile(
    r'<(script|style|iframe|object|embed|form|svg|math)[^>]*>.*?</\1>',
    re.IGNORECASE | re.DOTALL
)
# 危险自闭合或空标签（含 template、details 可被滥用）
DANGEROUS_TAG = re.compile(
    r'<(script|iframe|object|embed|input|button|textarea|select|meta|link|base|template|details|dialog)[^>]*/?>',
    re.IGNORECASE
)
# on* 事件属性
ON_EVENT = re.compile(r'\s+on\w+\s*=\s*["\'][^"\']*["\']', re.IGNORECASE)
# javascript:/data:/vbscript: 等危险协议
JS_PROTO = re.compile(r'(href|src)\s*=\s*["\']\s*(javascript|data|vbscript)\s*:[^"\']*["\']', re.IGNORECASE)
# style 内 expression、url(javascript:)
STYLE_DANGER = re.compile(r'expression\s*\(|url\s*\(\s*["\']?\s*javascript\s*:', re.IGNORECASE)
# formaction、poster 等可执行属性
FORMACTION = re.compile(r'\s+formaction\s*=\s*["\'][^"\']*["\']', re.IGNORECASE)


def _sanitize_regex(html_str):
    """自研正则净化（无 bleach 时回退）"""
    s = html_str
    s = DANGEROUS_BLOCK.sub('', s)
    s = DANGEROUS_TAG.sub('', s)
    s = ON_EVENT.sub('', s)
    s = JS_PROTO.sub(r'\1="#"', s)
    s = STYLE_DANGER.sub('', s)
    s = FORMACTION.sub('', s)
    return s


def sanitize_html(html_str):
    """
    净化 HTML，移除危险标签与事件，防止 XSS。
    优先使用 bleach（若已安装），否则使用自研正则。
    若净化失败则返回转义后的纯文本。
    """
    if not html_str or not isinstance(html_str, str):
        return ''
    try:
        try:
            import bleach
            ALLOWED_TAGS = ['p', 'br', 'strong', 'em', 'u', 's', 'a', 'ul', 'ol', 'li', 'h1', 'h2', 'h3', 'h4', 'h5', 'h6', 'blockquote', 'code', 'pre', 'img', 'span', 'div', 'table', 'thead', 'tbody', 'tr', 'th', 'td']
            ALLOWED_ATTRS = {'href', 'title', 'src', 'alt', 'class', 'target', 'rel'}
            return bleach.clean(html_str, tags=ALLOWED_TAGS, attributes=ALLOWED_ATTRS, strip=True)
        except ImportError:
            return _sanitize_regex(html_str)
    except Exception:
        return html.escape(html_str)
