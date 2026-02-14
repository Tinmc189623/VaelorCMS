from django.shortcuts import render, redirect, get_object_or_404
from django.contrib.auth import login, logout, authenticate
from django.contrib.auth.decorators import login_required
from django.contrib.auth.views import PasswordChangeView
from django.contrib import messages
from django.urls import reverse_lazy
from .forms import LoginForm, RegisterForm, ProfileForm, PasswordChangeFormCustom
from .models import UserProfile


def login_view(request):
    if request.user.is_authenticated:
        return redirect('home')
    from .login_throttle import is_locked, record_fail, clear_on_success
    locked, remain_secs = is_locked(request)
    if locked:
        mins = max(1, remain_secs // 60)
        messages.error(request, f'登录失败次数过多，请 {mins} 分钟后再试')
        return render(request, 'auth/login.html', {'form': LoginForm()})
    if request.method == 'POST':
        form = LoginForm(request.POST)
        if form.is_valid():
            user = authenticate(request, username=form.cleaned_data['username'], password=form.cleaned_data['password'])
            if user:
                clear_on_success(request)
                login(request, user)
                return redirect('home')
            record_fail(request)
        messages.error(request, '用户名或密码错误')
    else:
        form = LoginForm()
    return render(request, 'auth/login.html', {'form': form})


def register_view(request):
    from site_app.settings_service import get
    if get('allow_register', '1') != '1':
        from django.contrib import messages
        messages.error(request, '注册已关闭')
        return redirect('home')
    if request.method == 'POST':
        form = RegisterForm(request.POST)
        if form.is_valid():
            form.save()
            messages.success(request, '注册成功，请登录')
            return redirect('login')
    else:
        form = RegisterForm()
    return render(request, 'auth/register.html', {'form': form})


def logout_view(request):
    logout(request)
    return redirect('home')


@login_required
def settings_profile(request):
    """用户设置 - 资料"""
    profile, _ = UserProfile.objects.get_or_create(user=request.user, defaults={})
    if request.method == 'POST':
        form = ProfileForm(request.POST, request.FILES, instance=profile, user=request.user)
        if form.is_valid():
            form.save()
            messages.success(request, '资料已保存')
            return redirect('settings_profile')
    else:
        form = ProfileForm(instance=profile, user=request.user)
    return render(request, 'auth/settings_profile.html', {'form': form})


@login_required
def settings_password(request):
    """用户设置 - 修改密码"""
    if request.method == 'POST':
        form = PasswordChangeFormCustom(user=request.user, data=request.POST)
        if form.is_valid():
            form.save()
            logout(request)
            messages.success(request, '密码已修改，请重新登录')
            return redirect('login')
    else:
        form = PasswordChangeFormCustom(user=request.user)
    return render(request, 'auth/settings_password.html', {'form': form})


@login_required
def settings_account(request):
    """用户设置 - 账户与安全（会话等）"""
    # 当前会话信息
    session_key = request.session.session_key
    return render(request, 'auth/settings_account.html', {'session_key': session_key})


@login_required
def settings_logout_others(request):
    """登出其他设备（仅当前用户的其他会话）"""
    from django.contrib.sessions.models import Session
    from django.utils import timezone
    current_key = request.session.session_key
    uid = str(request.user.id)
    qs = Session.objects.exclude(session_key=current_key).filter(expire_date__gte=timezone.now())
    for s in qs:
        try:
            data = s.get_decoded()
            if data.get('_auth_user_id') == uid:
                s.delete()
        except Exception:
            pass
    messages.success(request, '已登出其他设备')
    return redirect('settings_account')
