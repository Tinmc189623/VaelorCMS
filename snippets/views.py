from django.shortcuts import render, redirect, get_object_or_404
from django.contrib.auth.decorators import login_required
from django.contrib import messages
from .models import CodeSnippet


def index(request):
    items = CodeSnippet.objects.select_related('user').all()[:100]
    return render(request, 'code/index.html', {'items': items})


def detail(request, pk):
    item = get_object_or_404(CodeSnippet.objects.select_related('user'), pk=pk)
    return render(request, 'code/detail.html', {'item': item})


@login_required
def submit(request):
    if request.method == 'POST':
        title = request.POST.get('title', '').strip()
        code = request.POST.get('code', '').strip()[:100000]
        language = request.POST.get('language', 'text').strip() or 'text'
        if not title:
            messages.error(request, '请填写标题')
        else:
            CodeSnippet.objects.create(user=request.user, title=title, code=code, language=language)
            messages.success(request, '提交成功')
            return redirect('code_index')
    return redirect('code_index')
