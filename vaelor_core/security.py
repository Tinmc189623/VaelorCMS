"""
Vaelor Core - 安全工具
输入验证、安全字符串处理
"""
import re
import html


def safe_str(value, max_length=0):
    """安全转换为字符串，移除控制字符"""
    if value is None:
        return ''
    s = str(value).strip()
    # 移除 ASCII 控制字符
    s = re.sub(r'[\x00-\x08\x0b\x0c\x0e-\x1f\x7f]', '', s)
    if max_length > 0 and len(s) > max_length:
        s = s[:max_length]
    return s


def validate_input(value, allowed_pattern=None, max_length=1000):
    """
    验证用户输入
    allowed_pattern: 允许的正则模式，如 r'^[\\w\\s\\-\\.]+$'
    返回 (is_valid, sanitized_value)
    """
    s = safe_str(value, max_length)
    if not s:
        return True, ''
    if allowed_pattern and not re.match(allowed_pattern, s):
        return False, html.escape(s)
    return True, s


def validate_slug(slug):
    """验证 URL slug（字母数字、连字符、下划线）"""
    if not slug or not isinstance(slug, str):
        return False, ''
    s = slug.strip().lower()
    if not re.match(r'^[a-z0-9\-_]+$', s):
        return False, ''
    return True, s[:128]


def validate_email(email):
    """简单邮箱格式验证"""
    if not email or not isinstance(email, str):
        return False, ''
    s = email.strip().lower()
    if len(s) > 254:
        return False, ''
    if not re.match(r'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$', s):
        return False, ''
    return True, s
