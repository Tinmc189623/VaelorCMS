from django.shortcuts import render, redirect, get_object_or_404
from django.contrib import messages
from django.contrib.auth.decorators import login_required
from django.contrib.admin.views.decorators import staff_member_required
from django.db.models import Q
from bbs.models import BbsPost
from snippets.models import CodeSnippet
from articles.models import Article
from users.models import User
from .models import AdminLog


def home(request):
    return render(request, 'home.html')


def robots_txt(request):
    """动态 robots.txt，注入 sitemap URL"""
    from django.http import HttpResponse
    from .settings_service import get
    base = get('seo_canonical_base', '').rstrip('/')
    sitemap_line = f'Sitemap: {base}/sitemap.xml\n' if base else '# Sitemap: https://your-domain.com/sitemap.xml\n'
    content = f"""User-agent: *
Allow: /
Disallow: /admin/
Disallow: /install/

{sitemap_line}"""
    return HttpResponse(content, content_type='text/plain')


def search(request):
    q = request.GET.get('q', '').strip()
    results = {}
    if q:
        results['bbs'] = BbsPost.objects.filter(Q(title__icontains=q) | Q(content__icontains=q)).select_related('user')[:20]
        results['code'] = CodeSnippet.objects.filter(Q(title__icontains=q) | Q(code__icontains=q)).select_related('user')[:20]
        results['articles'] = Article.objects.filter(
            Q(title__icontains=q) | Q(content__icontains=q),
            status='published'
        ).select_related('user')[:20]
    return render(request, 'search.html', {'q': q, 'results': results})


def about(request):
    return render(request, 'about.html')


def help_page(request):
    return render(request, 'help.html')


def faq(request):
    return render(request, 'faq.html')


def _safe_redirect_back(request, default='home'):
    """安全回跳：仅当 Referer 为本站时使用，防开放重定向"""
    ref = request.META.get('HTTP_REFERER', '').strip()
    if not ref:
        return redirect(default)
    try:
        from urllib.parse import urlparse
        from django.conf import settings
        parsed = urlparse(ref)
        if parsed.scheme not in ('http', 'https', ''):
            return redirect(default)
        req_host = request.get_host().split(':')[0].lower()
        ref_host = (parsed.netloc or '').split(':')[0].lower()
        if ref_host and ref_host == req_host:
            return redirect(ref)
    except Exception:
        pass
    return redirect(default)


def newsletter_subscribe(request):
    if request.method != 'POST':
        return redirect('home')
    email = request.POST.get('email', '').strip()
    if not email or '@' not in email or len(email) > 254:
        messages.error(request, '请输入有效邮箱')
        return _safe_redirect_back(request)
    try:
        from .models import NewsletterSubscriber
        NewsletterSubscriber.objects.get_or_create(email=email, defaults={'is_active': True})
        messages.success(request, '订阅成功，感谢关注！')
    except Exception:
        messages.success(request, '订阅成功，感谢关注！')
    return _safe_redirect_back(request)


def page_404(request, exception=None):
    return render(request, '404.html', status=404)


def page_500(request):
    return render(request, '500.html', status=500)


def page_detail(request, slug):
    """自定义页面 - 管理员创建，通过 /p/<slug>/ 访问"""
    from .models import Page
    from django.http import Http404
    from django.utils.safestring import mark_safe
    from .html_sanitizer import sanitize_html
    page = Page.objects.filter(slug=slug, is_published=True).first()
    if not page:
        raise Http404('页面不存在')
    from site_app.seo_middleware import set_request_seo
    desc = (page.content or '')[:160] + ('...' if len(page.content or '') > 160 else '')
    set_request_seo(request, title=page.title, description=desc)
    # HTML 内容经净化后输出，防 XSS
    content_safe = mark_safe(sanitize_html(page.content)) if page.content_is_html else None
    return render(request, 'page.html', {'page': page, 'content_safe': content_safe})


def games(request):
    game = request.GET.get('g', '')
    if game not in ('snake', 'memory', 'guess'):
        game = ''
    return render(request, 'games.html', {'game': game})


@login_required
def profile(request):
    posts = BbsPost.objects.filter(user=request.user).order_by('-id')[:50]
    snippets = CodeSnippet.objects.filter(user=request.user).order_by('-id')[:50]
    return render(request, 'profile.html', {'posts': posts, 'snippets': snippets})


@staff_member_required
def admin_panel(request):
    stats = {
        'users': User.objects.count(),
        'bbs': BbsPost.objects.count(),
        'code': CodeSnippet.objects.count(),
        'articles': Article.objects.count(),
    }
    return render(request, 'admin/index.html', {'stats': stats})


@staff_member_required
def admin_users(request):
    users = User.objects.all().order_by('-id')[:200]
    return render(request, 'admin/users.html', {'users': users})


@staff_member_required
def admin_bbs(request):
    items = BbsPost.objects.select_related('user').order_by('-id')[:200]
    return render(request, 'admin/bbs.html', {'items': items})


@staff_member_required
def admin_bbs_approve(request, pk):
    if request.method != 'POST':
        return redirect('admin_bbs')
    post = get_object_or_404(BbsPost, pk=pk)
    post.approved = True
    post.save()
    from django.contrib import messages
    messages.success(request, '帖子已通过审核')
    return redirect('admin_bbs')


@staff_member_required
def admin_code(request):
    items = CodeSnippet.objects.select_related('user').order_by('-id')[:200]
    return render(request, 'admin/code.html', {'items': items})


@staff_member_required
def admin_logs(request):
    logs = AdminLog.objects.all().order_by('-id')[:100]
    return render(request, 'admin/logs.html', {'logs': logs})


@staff_member_required
def admin_rescue(request):
    """站点急救 - 清除缓存、恢复配置、修复 config.ini"""
    from django.contrib import messages
    from django.http import HttpResponse
    from .settings_service import restore_defaults, export_backup, import_backup
    if request.method != 'POST':
        return render(request, 'admin/rescue.html')
    action = request.POST.get('action', '')
    if action == 'clear_cache':
        from django.core.cache import cache
        cache.clear()
        messages.success(request, '缓存已清除')
    elif action == 'clear_login_lockout':
        from django.core.cache import cache
        try:
            if hasattr(cache, 'delete_pattern'):
                cache.delete_pattern('*login_lockout*')
                messages.success(request, '登录锁定已清除')
            else:
                cache.clear()
                messages.success(request, '已清除全部缓存（含登录锁定）')
        except Exception:
            cache.clear()
            messages.success(request, '已清除全部缓存（含登录锁定）')
    elif action == 'disable_maintenance':
        from .settings_service import set
        set('maintenance_mode', '0', 'maintenance')
        messages.success(request, '维护模式已关闭')
    elif action == 'restore_all':
        count = restore_defaults(category=None)
        messages.success(request, f'已恢复全部 {count} 项配置为默认值')
    elif action == 'fix_config_ini':
        import shutil
        from pathlib import Path
        from django.conf import settings
        sample = Path(settings.BASE_DIR) / 'config' / 'config.ini.sample'
        target = Path(settings.BASE_DIR) / 'config' / 'config.ini'
        if sample.exists():
            shutil.copy2(sample, target)
            messages.success(request, '已从 config.ini.sample 恢复 config.ini，请重启服务生效')
        else:
            messages.error(request, 'config.ini.sample 不存在')
    elif action == 'export_backup':
        data = export_backup()
        resp = HttpResponse(data, content_type='application/json; charset=utf-8')
        resp['Content-Disposition'] = 'attachment; filename="vaelor_config_backup.json"'
        return resp
    elif action == 'import_backup':
        f = request.FILES.get('backup_file')
        if f and f.name.endswith('.json'):
            try:
                content = f.read().decode('utf-8')
                ok, fail = import_backup(content)
                messages.success(request, f'已恢复 {ok} 项配置' + (f'，{fail} 项失败' if fail else ''))
            except Exception as e:
                messages.error(request, f'恢复失败：{e}')
        else:
            messages.error(request, '请上传 .json 备份文件')
    return redirect('admin_rescue')


@staff_member_required
def admin_security(request):
    """Thalix 安全审计 - 检查配置与运行环境"""
    from .thalix_security import run_security_audit, get_audit_summary
    results = run_security_audit()
    summary = get_audit_summary()
    return render(request, 'admin/security.html', {'results': results, 'summary': summary})


@staff_member_required
def admin_settings(request):
    """管理员设置 - 站点、安全、用户、内容、维护、高级"""
    from .settings_service import get_all_by_category, set, seed_defaults, restore_defaults, DEFAULTS
    from .models import SiteSetting
    from django.contrib import messages
    seed_defaults()
    tab = request.GET.get('tab', 'general')
    if request.method == 'POST':
        tab = request.POST.get('tab', 'general')
        if request.POST.get('action') == 'restore':
            if tab == 'advanced':
                messages.error(request, '高级模式下请使用「站点急救」恢复全部')
            else:
                count = restore_defaults(category=tab)
                messages.success(request, f'已恢复 {count} 项配置为默认值')
            return redirect(f'{request.path}?tab={tab}')
        if tab == 'advanced':
            import re
            for k, v in request.POST.items():
                if k.startswith('adv_key_'):
                    sid = k.replace('adv_key_', '')
                    if request.POST.get(f'adv_del_{sid}') == '1':
                        SiteSetting.objects.filter(pk=sid).delete()
                        continue
                    key = (v or '').strip()[:64]
                    val = request.POST.get(f'adv_value_{sid}', '')
                    cat = (request.POST.get(f'adv_cat_{sid}', 'general') or 'general')[:32]
                    if key and re.match(r'^[a-zA-Z0-9_]+$', key):
                        try:
                            obj = SiteSetting.objects.get(pk=sid)
                            obj.key, obj.value, obj.category = key, str(val), cat
                            obj.save()
                        except (SiteSetting.DoesNotExist, Exception):
                            pass
            nk = (request.POST.get('new_key') or '').strip()[:64]
            nv = request.POST.get('new_value', '')
            nc = (request.POST.get('new_category') or 'general')[:32]
            if nk and re.match(r'^[a-zA-Z0-9_]+$', nk):
                SiteSetting.objects.update_or_create(key=nk, defaults={'value': str(nv), 'category': nc})
            from django.core.cache import cache
            cache.delete('vaelor_site_settings_ctx')
            messages.success(request, '高级配置已保存')
        else:
            for key, (default, cat) in DEFAULTS.items():
                if cat == tab:
                    val = request.POST.get(key, default)
                    set(key, val, cat)
            messages.success(request, '设置已保存')
        return redirect(f'{request.path}?tab={tab}')
    data = get_all_by_category()
    themes = [{'id': 'default', 'name': '默认（暖色纸本）'}]
    try:
        from vaelor.theme import list_themes
        extra = list_themes()
        if extra:
            themes = themes + [t for t in extra if t.get('id') != 'default']
    except Exception:
        pass
    advanced_items = []
    if tab == 'advanced':
        advanced_items = list(SiteSetting.objects.all().order_by('category', 'key').values('id', 'key', 'value', 'category'))
    return render(request, 'admin/settings.html', {'data': data, 'tab': tab, 'themes': themes, 'advanced_items': advanced_items})
