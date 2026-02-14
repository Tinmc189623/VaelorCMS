"""SEO 中间件 - 为 request 设置默认 page_title，供模板使用"""
_URL_NAME_TITLES = {
    'home': '首页',
    'bbs_index': '论坛',
    'bbs_detail': '帖子详情',
    'code_index': '代码分享',
    'code_detail': '代码详情',
    'articles_index': '文章',
    'article_detail': '文章详情',
    'search': '搜索',
    'help': '帮助',
    'about': '关于',
    'games': '小游戏',
    'login': '登录',
    'register': '注册',
    'profile': '用户中心',
    'admin_panel': '管理后台',
}


class SEOMiddleware:
    def __init__(self, get_response):
        self.get_response = get_response

    def __call__(self, request):
        request.page_title = '首页'
        request.page_description = ''
        request.page_keywords = ''
        try:
            from django.urls import resolve
            match = resolve(request.path)
            if match and match.url_name:
                request.page_title = _URL_NAME_TITLES.get(match.url_name, request.page_title)
        except Exception:
            pass
        response = self.get_response(request)
        return response


def set_request_seo(request, title=None, description=None, keywords=None):
    """在视图中调用，设置当前页 SEO 信息"""
    if title is not None:
        request.page_title = title
    if description is not None:
        request.page_description = description
    if keywords is not None:
        request.page_keywords = keywords
