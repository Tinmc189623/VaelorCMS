"""
密码校验 - 统一使用管理员配置的密码策略
"""
from django.core.exceptions import ValidationError


def _special_chars():
    """允许的特殊字符"""
    return set('!@#$%^&*()_+-=[]{}|;:,.<>?/~`')


def validate_password(password):
    """
    根据站点配置校验密码，不通过则抛出 ValidationError
    配置项：min_password_length, require_strong_password, require_password_special
    """
    if not password:
        return
    try:
        from site_app.settings_service import get
        min_len = int(get('min_password_length', '8'))
        strong = get('require_strong_password', '0') == '1'
        special = get('require_password_special', '0') == '1'
    except (ImportError, ValueError):
        min_len, strong, special = 8, False, False

    if len(password) < min_len:
        raise ValidationError(f'密码至少 {min_len} 位')

    if strong:
        if not any(c.isalpha() for c in password) or not any(c.isdigit() for c in password):
            raise ValidationError('密码须包含字母与数字')

    if special:
        if not any(c in _special_chars() for c in password):
            raise ValidationError('密码须包含至少一个特殊字符（如 !@#$%^&*）')
