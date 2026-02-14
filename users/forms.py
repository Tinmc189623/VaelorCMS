from django import forms
from django.contrib.auth.forms import UserCreationForm, PasswordChangeForm
from django.core.exceptions import ValidationError
from .models import User, UserProfile


class LoginForm(forms.Form):
    username = forms.CharField(max_length=64, label='用户名')
    password = forms.CharField(widget=forms.PasswordInput, label='密码')


class RegisterForm(UserCreationForm):
    email = forms.EmailField(required=False, label='邮箱')

    class Meta:
        model = User
        fields = ('username', 'email', 'password1', 'password2')

    def clean_password1(self):
        p1 = self.cleaned_data.get('password1')
        if p1:
            from .password_validators import validate_password
            try:
                validate_password(p1)
            except ValidationError:
                raise
        return p1

    def save(self, commit=True):
        user = super().save(commit=False)
        user.email = self.cleaned_data.get('email', '')
        user.role = 'user'
        if commit:
            user.save()
        return user


class ProfileForm(forms.ModelForm):
    email = forms.EmailField(required=False, label='邮箱')

    class Meta:
        model = UserProfile
        fields = ('nickname', 'bio', 'avatar', 'show_email', 'profile_visible',
                  'email_on_reply', 'email_on_mention', 'timezone', 'language')
        widgets = {
            'bio': forms.Textarea(attrs={'rows': 4}),
        }

    def __init__(self, *args, **kwargs):
        self.user = kwargs.pop('user', None)
        super().__init__(*args, **kwargs)
        if self.user:
            self.fields['email'].initial = self.user.email

    def save(self, commit=True):
        profile = super().save(commit=commit)
        if self.user and 'email' in self.cleaned_data:
            self.user.email = self.cleaned_data.get('email', '')
            self.user.save(update_fields=['email'])
        return profile


class PasswordChangeFormCustom(PasswordChangeForm):
    """继承 Django 密码修改表单，可扩展验证"""

    def clean_new_password1(self):
        p1 = self.cleaned_data.get('new_password1')
        if p1:
            from .password_validators import validate_password
            try:
                validate_password(p1)
            except ValidationError:
                raise
        return p1
