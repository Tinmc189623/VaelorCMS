from django.shortcuts import render, redirect, get_object_or_404
from django.contrib.auth.decorators import login_required
from django.contrib import messages
from django.db.models import Q, F, Prefetch
from django.db.models.functions import TruncMonth
from django.core.cache import cache
from django.views.decorators.http import require_http_methods
from .models import Article, ArticleComment, ArticleLike, ArticleFavorite

_ARTICLE_CACHE_TIMEOUT = 300  # 5 分钟


def _invalidate_article_cache():
    """文章/分类/归档变更时清除缓存"""
    for key in ('vaelor_article_tag_cloud', 'vaelor_article_categories', 'vaelor_article_archives'):
        cache.delete(key)


def _get_tag_cloud(limit=30):
    """从已发布文章中提取标签及出现次数（仅取 tags 字段，减少内存）"""
    cache_key = 'vaelor_article_tag_cloud'
    cached = cache.get(cache_key)
    if cached is not None:
        return cached
    rows = Article.objects.filter(status='published').exclude(tags='').values_list('tags', flat=True)
    tag_counts = {}
    for tags_str in rows:
        for t in (tags_str or '').replace('，', ',').split(','):
            t = t.strip()
            if t:
                tag_counts[t] = tag_counts.get(t, 0) + 1
    result = sorted(tag_counts.items(), key=lambda x: -x[1])[:limit]
    cache.set(cache_key, result, _ARTICLE_CACHE_TIMEOUT)
    return result


def index(request):
    if request.user.is_authenticated:
        qs = Article.objects.filter(Q(status='published') | Q(user=request.user))
    else:
        qs = Article.objects.filter(status='published')
    category = request.GET.get('category', '').strip()
    tag = request.GET.get('tag', '').strip()
    if category:
        qs = qs.filter(category=category)
    if tag:
        qs = qs.filter(tags__icontains=tag)
    month = request.GET.get('month', '').strip()
    if month and len(month) == 7:  # YYYY-MM
        qs = qs.filter(created_at__year=month[:4], created_at__month=month[5:7])
    articles = qs.select_related('user').order_by('-is_pinned', '-id')[:100]
    def _get_categories():
        return list(Article.objects.filter(status='published').exclude(category='').values_list('category', flat=True).distinct()[:20])
    categories = cache.get_or_set('vaelor_article_categories', _get_categories, _ARTICLE_CACHE_TIMEOUT)
    tag_cloud = _get_tag_cloud()
    def _build_archives():
        months = Article.objects.filter(status='published').annotate(month=TruncMonth('created_at')).values_list('month', flat=True).distinct().order_by('-month')[:24]
        return [(m.strftime('%Y-%m'), m.strftime('%Y年%m月')) for m in months if m]
    archives = cache.get_or_set('vaelor_article_archives', _build_archives, _ARTICLE_CACHE_TIMEOUT)
    return render(request, 'articles/index.html', {
        'articles': articles, 'categories': categories, 'current_category': category,
        'tag_cloud': tag_cloud, 'current_tag': tag, 'archives': archives, 'current_month': month,
    })


def detail(request, pk):
    article = get_object_or_404(Article.objects.select_related('user'), pk=pk)
    can_view = article.status == 'published' or (request.user.is_authenticated and (article.user_id == request.user.id or request.user.role == 'admin'))
    if not can_view:
        from django.http import HttpResponseForbidden
        return HttpResponseForbidden('无权访问')
    Article.objects.filter(pk=pk).update(view_count=F('view_count') + 1)
    article.refresh_from_db()
    from site_app.seo_middleware import set_request_seo
    content = article.content or ''
    desc = content[:160] + ('...' if len(content) > 160 else '')
    set_request_seo(request, title=article.title, description=desc, keywords=article.tags)
    comments = []
    show_comments = False
    from site_app.settings_service import get
    if get('article_comment', '0') == '1':
        show_comments = True
        replies_qs = ArticleComment.objects.filter(is_approved=True).select_related('user')
        comments = list(ArticleComment.objects.filter(article=article, is_approved=True, parent__isnull=True)
                        .select_related('user').prefetch_related(Prefetch('replies', queryset=replies_qs))[:200])
    liked = favorited = False
    if request.user.is_authenticated:
        liked = ArticleLike.objects.filter(article=article, user=request.user).exists()
        favorited = ArticleFavorite.objects.filter(article=article, user=request.user).exists()
    return render(request, 'articles/detail.html', {
        'article': article, 'comments': comments, 'show_comments': show_comments,
        'liked': liked, 'favorited': favorited,
    })


@login_required
def create(request):
    return render(request, 'articles/form.html', {'article': None})


@login_required
def edit(request, pk):
    article = get_object_or_404(Article, pk=pk)
    if article.user_id != request.user.id and request.user.role != 'admin':
        from django.http import HttpResponseForbidden
        return HttpResponseForbidden('无权编辑')
    return render(request, 'articles/form.html', {'article': article})


@login_required
def save_article(request):
    if request.method != 'POST':
        return redirect('articles_index')
    pk = request.POST.get('id')
    title = request.POST.get('title', '').strip()
    content = request.POST.get('content', '').strip()
    action = request.POST.get('action', 'save')
    status = 'published' if action == 'publish' else request.POST.get('status', 'draft')
    tags = request.POST.get('tags', '').strip()
    category = request.POST.get('category', '').strip()
    is_pinned = request.POST.get('is_pinned') == '1'
    if not title or not content:
        messages.error(request, '标题和内容必填')
        if pk:
            return redirect('article_edit', pk=pk)
        return redirect('article_create')
    if pk:
        article = get_object_or_404(Article, pk=pk)
        if article.user_id != request.user.id and request.user.role != 'admin':
            from django.http import HttpResponseForbidden
            return HttpResponseForbidden('无权编辑')
        article.title, article.content, article.status, article.tags, article.category, article.is_pinned = title, content, status, tags, category, is_pinned
        article.save()
        _invalidate_article_cache()
        return redirect('article_detail', pk=pk)
    else:
        article = Article.objects.create(user=request.user, title=title, content=content, status=status, tags=tags, category=category, is_pinned=is_pinned)
        _invalidate_article_cache()
        return redirect('article_detail', pk=article.pk)


@require_http_methods(['POST'])
def add_comment(request, pk):
    from site_app.settings_service import get
    if get('article_comment', '0') != '1':
        messages.error(request, '评论功能已关闭')
        return redirect('article_detail', pk=pk)
    article = get_object_or_404(Article, pk=pk)
    if article.status != 'published':
        messages.error(request, '该文章未发布')
        return redirect('article_detail', pk=pk)
    content = request.POST.get('content', '').strip()
    if not content or len(content) > 2000:
        messages.error(request, '评论内容 1–2000 字')
        return redirect('article_detail', pk=pk)
    parent_id = request.POST.get('parent_id')
    parent = None
    if parent_id:
        try:
            parent = ArticleComment.objects.get(pk=int(parent_id), article=article)
        except (ValueError, ArticleComment.DoesNotExist):
            pass
    if request.user.is_authenticated:
        ArticleComment.objects.create(article=article, user=request.user, content=content, parent=parent)
    else:
        author_name = request.POST.get('author_name', '').strip() or '游客'
        ArticleComment.objects.create(article=article, author_name=author_name[:64], content=content, parent=parent)
    messages.success(request, '评论已提交')
    return redirect('article_detail', pk=pk)


@login_required
@require_http_methods(['POST'])
def like_article(request, pk):
    article = get_object_or_404(Article, pk=pk, status='published')
    obj, created = ArticleLike.objects.get_or_create(article=article, user=request.user)
    if created:
        Article.objects.filter(pk=pk).update(like_count=F('like_count') + 1)
    else:
        obj.delete()
        Article.objects.filter(pk=pk, like_count__gt=0).update(like_count=F('like_count') - 1)
    return redirect('article_detail', pk=pk)


@login_required
@require_http_methods(['POST'])
def favorite_article(request, pk):
    article = get_object_or_404(Article, pk=pk, status='published')
    obj, created = ArticleFavorite.objects.get_or_create(article=article, user=request.user)
    if not created:
        obj.delete()
    messages.success(request, '已收藏' if created else '已取消收藏')
    return redirect('article_detail', pk=pk)
