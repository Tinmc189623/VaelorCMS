from django.shortcuts import render, redirect, get_object_or_404
from django.contrib.auth.decorators import login_required
from django.contrib import messages
from .models import BbsPost


def _can_view_unapproved(request):
    return request.user.is_authenticated and getattr(request.user, 'role', '') == 'admin'


def index(request):
    from site_app.settings_service import get
    qs = BbsPost.objects.select_related('user')
    if not _can_view_unapproved(request):
        qs = qs.filter(approved=True)
    posts = qs.order_by('-id')[:100]
    can_post = request.user.is_authenticated or get('bbs_guest_post', '0') == '1'
    return render(request, 'bbs/index.html', {'posts': posts, 'can_post': can_post})


def detail(request, pk):
    post = get_object_or_404(BbsPost.objects.select_related('user'), pk=pk)
    if not post.approved and not _can_view_unapproved(request):
        from django.http import HttpResponseForbidden
        return HttpResponseForbidden('帖子待审核')
    return render(request, 'bbs/detail.html', {'post': post})


def post(request):
    from site_app.settings_service import get
    guest_ok = get('bbs_guest_post', '0') == '1'
    need_login = not request.user.is_authenticated and not guest_ok
    if need_login:
        messages.error(request, '请先登录后发帖')
        return redirect('bbs_index')
    if request.method == 'POST':
        title = request.POST.get('title', '').strip()
        content = request.POST.get('content', '').strip()[:50000]
        if not title:
            messages.error(request, '请填写标题')
        else:
            need_mod = get('bbs_moderation', '0') == '1'
            BbsPost.objects.create(
                user=request.user if request.user.is_authenticated else None,
                title=title,
                content=content,
                approved=not need_mod,
            )
            messages.success(request, '发帖成功' + ('，待审核后显示' if need_mod else ''))
            return redirect('bbs_index')
    return redirect('bbs_index')
