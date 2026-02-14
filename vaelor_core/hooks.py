"""
Vaelor Core - 钩子系统
插件与主题可注册钩子，在模板渲染时注入内容
"""
_REGISTRY = {}


def register_hook(hook_name, callback, priority=10):
    """注册钩子回调，priority 越小越先执行"""
    if hook_name not in _REGISTRY:
        _REGISTRY[hook_name] = []
    _REGISTRY[hook_name].append((priority, callback))
    _REGISTRY[hook_name].sort(key=lambda x: x[0])


def run_hook(hook_name, request=None, context=None):
    """执行钩子，返回所有回调结果的列表"""
    context = context or {}
    results = []
    for _, cb in _REGISTRY.get(hook_name, []):
        try:
            r = cb(request, context) if callable(cb) else []
            if r:
                results.extend(r if isinstance(r, (list, tuple)) else [r])
        except Exception:
            pass
    return results
