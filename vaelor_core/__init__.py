"""
Vaelor Core - VaelorCMS 内核
提供配置、缓存、钩子、安全等核心抽象与工具
"""
__version__ = '1.0.0'

from .config import get_config, set_config
from .security import safe_str, validate_input, validate_slug, validate_email

# hooks 与 vaelor.plugin 独立，此处不强制导入
def run_hook(hook_name, request=None, context=None):
    try:
        from vaelor.plugin import run_hook as _run
        return _run(hook_name, request, context)
    except ImportError:
        return []

def register_hook(hook_name, callback, priority=10):
    try:
        from vaelor.plugin import add_hook
        add_hook(hook_name, callback, priority)
    except ImportError:
        pass

__all__ = [
    'get_config', 'set_config',
    'run_hook', 'register_hook',
    'safe_str', 'validate_input', 'validate_slug', 'validate_email',
]
