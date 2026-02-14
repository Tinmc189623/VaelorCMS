# 插件开发指南

> 写给想给 VaelorCMS 做插件的开发者。基于钩子（hooks）的扩展机制，不侵入核心代码。

---

## 一、插件是啥

插件通过**钩子**在固定位置注入自己的逻辑或 HTML。比如：在页脚加统计代码、在导航加一个链接、在文章页加「分享」按钮。

**你能做的：**

- 往 `<head>` 里加脚本、样式
- 往导航、页脚、内容区前后加 HTML
- 注册自己的 URL、视图（需在项目里加，见下）

**你不能做的（不通过钩子）：**

- 改核心模型、改数据库结构（除非自己写 migration）
- 改核心模板（除非 fork）

---

## 二、钩子列表

| 钩子名 | 位置 | 回调参数 | 返回值 |
|--------|------|----------|--------|
| head_extra | `<head>` 末尾 | (request, context) | str，HTML |
| body_top | `<body>` 开头 | (request, context) | str |
| nav_extra | 导航栏 | (request, context) | str |
| content_before | 主内容前 | (request, context) | str |
| content_after | 主内容后 | (request, context) | str |
| footer_extra | 页脚内 | (request, context) | str |
| body_bottom | `</body>` 前 | (request, context) | str |

---

## 三、写一个插件

### 1. 创建插件模块

在项目里新建 `plugins/my_plugin/`（或任意位置），写一个 `__init__.py`：

```python
# plugins/my_plugin/__init__.py
from vaelor.plugin import register_plugin, add_hook

def head_analytics(request, context):
    return '<script async src="https://example.com/analytics.js"></script>'

def init_plugin():
    register_plugin('my_plugin', '我的插件', '1.0', '注入统计代码')
    add_hook('head_extra', head_analytics, priority=10)

init_plugin()
```

### 2. 加载插件

在 Django 启动时加载。推荐在 `config/apps.py` 或 `config/__init__.py` 里：

```python
# config/__init__.py
def ready():
    from site_app.settings_service import get
    enabled = get('plugins_enabled', '').split(',')
    if 'my_plugin' in [p.strip() for p in enabled if p.strip()]:
        import plugins.my_plugin  # 触发 init_plugin
```

或者更简单：在 `main.py` 或 `config/wsgi.py` 里，根据配置 import 插件模块。

### 3. 启用插件

在管理后台「站点设置」→「主题」→「已启用插件」里填：`my_plugin`（多个用逗号分隔）。

---

## 四、API 速查

```python
from vaelor.plugin import register_plugin, add_hook, run_hook, get_plugins

# 注册插件
register_plugin(plugin_id, name, version='1.0', description='')

# 添加钩子
add_hook(hook_name, callback, priority=10)  # priority 越小越先执行

# 手动执行钩子（一般不需要，模板会调）
run_hook('head_extra', request, context)

# 获取已注册插件
get_plugins()
```

---

## 五、注意事项

1. **别抛异常**：钩子回调里尽量 try/except，别让插件拖垮整站
2. **返回值**：返回 `str` 或可转成 HTML 的内容，`None` 会被忽略
3. **request**：可能是 `None`（某些场景），做好判断
4. **性能**：钩子会在每次请求执行，别干重活

---

## 六、扩展 URL

如果插件要加自己的页面，需要在 `config/urls.py` 里手动加：

```python
urlpatterns = [
    ...
    path('my-plugin/', include('plugins.my_plugin.urls')),
]
```

插件本身不自动注册 URL，保持核心路由可控。
