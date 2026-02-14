"""
VaelorCMS 插件系统
基于钩子（hooks）的扩展机制，插件可注册到各扩展点
"""
_plugins = []
_hooks = {}  # hook_name -> [callable, ...]


def register_plugin(plugin_id, name, version='1.0', description=''):
    """注册插件元信息"""
    _plugins.append({
        'id': plugin_id,
        'name': name,
        'version': version,
        'description': description,
    })


def add_hook(hook_name, callback, priority=10):
    """
    向钩子添加回调
    hook_name: 钩子名称，如 'head_extra', 'footer_extra', 'nav_extra'
    callback: 可调用对象，接收 (request, context) 返回 str 或 dict
    priority: 数字越小越先执行，默认 10
    """
    if hook_name not in _hooks:
        _hooks[hook_name] = []
    _hooks[hook_name].append((priority, callback))
    _hooks[hook_name].sort(key=lambda x: x[0])


def run_hook(hook_name, request, context=None):
    """执行钩子，返回所有回调结果的列表"""
    context = context or {}
    results = []
    for _, cb in _hooks.get(hook_name, []):
        try:
            r = cb(request, context)
            if r is not None:
                results.append(r)
        except Exception:
            pass
    return results


def get_plugins():
    """获取已注册插件列表"""
    return list(_plugins)


# 预定义钩子（供开发者参考）
HOOKS = [
    'head_extra',      # <head> 内追加 HTML
    'body_top',        # <body> 开头
    'nav_extra',       # 导航栏追加项
    'content_before',  # 主内容前
    'content_after',   # 主内容后
    'footer_extra',    # 页脚追加
    'body_bottom',     # </body> 前
]
