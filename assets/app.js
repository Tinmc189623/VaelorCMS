/**
 * VaelorCMS - 现代 ES6+ 前端脚本
 */
'use strict';

const validateRequired = (el, msg) => {
  if (!el) return true;
  if (el.value.trim() === '') {
    alert(msg);
    el.focus();
    return false;
  }
  return true;
};

const validateBbsPost = () => {
  const title = document.getElementById('bbs_title');
  return validateRequired(title, '请填写标题');
};

const validateCodeSubmit = () => {
  const title = document.getElementById('code_title');
  return validateRequired(title, '请填写标题');
};

const validateLogin = () => {
  const u = document.getElementById('login_username');
  const p = document.getElementById('login_password');
  if (!u || !p) return true;
  if (u.value.trim() === '') {
    alert('请输入用户名');
    u.focus();
    return false;
  }
  if (p.value === '') {
    alert('请输入密码');
    p.focus();
    return false;
  }
  return true;
};

const validateRegister = () => {
  const u = document.getElementById('reg_username');
  const p1 = document.getElementById('reg_password');
  const p2 = document.getElementById('reg_password2');
  if (!u || !p1 || !p2) return true;
  const un = u.value.trim();
  if (un.length < 2) {
    alert('用户名至少 2 个字符');
    u.focus();
    return false;
  }
  if (p1.value.length < 6) {
    alert('密码至少 6 个字符');
    p1.focus();
    return false;
  }
  if (p1.value !== p2.value) {
    alert('两次密码不一致');
    p2.focus();
    return false;
  }
  return true;
};

document.addEventListener('DOMContentLoaded', () => {
  const bbsForm = document.getElementById('bbs_post_form');
  if (bbsForm) bbsForm.addEventListener('submit', e => { if (!validateBbsPost()) e.preventDefault(); });

  const codeForm = document.getElementById('code_submit_form');
  if (codeForm) codeForm.addEventListener('submit', e => { if (!validateCodeSubmit()) e.preventDefault(); });

  const loginForm = document.getElementById('login_form');
  if (loginForm) loginForm.addEventListener('submit', e => { if (!validateLogin()) e.preventDefault(); });

  const registerForm = document.getElementById('register_form');
  if (registerForm) registerForm.addEventListener('submit', e => { if (!validateRegister()) e.preventDefault(); });
});
