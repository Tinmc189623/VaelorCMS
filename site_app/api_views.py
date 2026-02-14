from django.http import JsonResponse
from django.db.models import Q
from bbs.models import BbsPost
from snippets.models import CodeSnippet
from articles.models import Article
from users.models import User

try:
    from config.version import __version__ as CMS_VERSION
except ImportError:
    CMS_VERSION = 'Demo-26.02.13.26'


def stats(request):
    data = {
        'users': User.objects.count(),
        'bbs': BbsPost.objects.count(),
        'code': CodeSnippet.objects.count(),
        'articles': Article.objects.count(),
    }
    return JsonResponse(data)


def search(request):
    q = request.GET.get('q', '').strip()
    results = {}
    if q:
        bbs = BbsPost.objects.filter(Q(title__icontains=q) | Q(content__icontains=q))[:10]
        results['bbs'] = [{'id': p.id, 'title': p.title, 'created_at': p.created_at.strftime('%Y-%m-%d %H:%M')} for p in bbs]
        code = CodeSnippet.objects.filter(Q(title__icontains=q) | Q(code__icontains=q))[:10]
        results['code'] = [{'id': s.id, 'title': s.title, 'language': s.language, 'created_at': s.created_at.strftime('%Y-%m-%d %H:%M')} for s in code]
        articles = Article.objects.filter(Q(title__icontains=q) | Q(content__icontains=q), status='published')[:10]
        results['articles'] = [{'id': a.id, 'title': a.title, 'created_at': a.created_at.strftime('%Y-%m-%d %H:%M')} for a in articles]
    return JsonResponse(results)


def upgrade(request):
    """
    升级接口 - 预留
    返回当前版本，后续可扩展：检查更新、执行升级脚本等
    """
    return JsonResponse({
        'version': CMS_VERSION,
        'upgrade_available': False,
        'message': '当前已是最新版本',
    })


def articles_list(request):
    """文章列表 API - 最近已发布文章"""
    from articles.models import Article
    category = request.GET.get('category', '').strip()
    qs = Article.objects.filter(status='published').order_by('-created_at')[:50]
    if category:
        qs = qs.filter(category=category)
    items = [{'id': a.id, 'title': a.title, 'category': a.category, 'created_at': a.created_at.strftime('%Y-%m-%d %H:%M')} for a in qs]
    return JsonResponse({'articles': items})


def health(request):
    """
    健康检查接口 - 用于负载均衡、监控探针
    """
    from django.db import connection
    try:
        connection.ensure_connection()
        db_ok = True
    except Exception:
        db_ok = False
    return JsonResponse({
        'status': 'ok' if db_ok else 'degraded',
        'version': CMS_VERSION,
        'database': 'connected' if db_ok else 'disconnected',
    })
