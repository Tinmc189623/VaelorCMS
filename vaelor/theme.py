"""
VaelorCMS 主题系统
支持多主题切换，主题目录：themes/<theme_id>/
"""
from pathlib import Path


def get_themes_dir():
    """主题根目录"""
    base = Path(__file__).resolve().parent.parent
    return base / 'themes'


def list_themes():
    """列出所有可用主题"""
    themes_dir = get_themes_dir()
    if not themes_dir.exists():
        return []
    result = []
    for d in themes_dir.iterdir():
        if d.is_dir() and not d.name.startswith('_'):
            info = _read_theme_info(d)
            info['id'] = d.name
            result.append(info)
    return result


def _read_theme_info(theme_dir):
    """读取主题 info.ini 或默认信息"""
    info = {'name': theme_dir.name, 'description': '', 'version': '1.0', 'author': ''}
    ini_path = theme_dir / 'info.ini'
    if ini_path.exists():
        try:
            import configparser
            cp = configparser.ConfigParser()
            cp.read(ini_path, encoding='utf-8')
            if cp.has_section('theme'):
                info['name'] = cp.get('theme', 'name', fallback=info['name'])
                info['description'] = cp.get('theme', 'description', fallback='')
                info['version'] = cp.get('theme', 'version', fallback='1.0')
                info['author'] = cp.get('theme', 'author', fallback='')
        except Exception:
            pass
    return info


def get_theme_css_path(theme_id):
    """获取主题 CSS 路径，供 static 或模板引用"""
    if not theme_id or theme_id == 'default':
        return None
    themes_dir = get_themes_dir()
    css_path = themes_dir / theme_id / 'style.css'
    if css_path.exists():
        return f'themes/{theme_id}/style.css'
    return None


def get_theme_static_path(theme_id, filename):
    """获取主题内静态文件路径"""
    if not theme_id or theme_id == 'default':
        return None
    themes_dir = get_themes_dir()
    f = themes_dir / theme_id / filename
    if f.exists():
        return f'themes/{theme_id}/{filename}'
    return None
