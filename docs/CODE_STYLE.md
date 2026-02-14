# 代码规范

VaelorCMS 项目代码风格约定。

## Python

- 遵循 PEP 8
- 行宽建议 100 字符以内
- 导入顺序：标准库 → 第三方 → 本地
- 使用 `pathlib.Path` 处理路径
- 异常处理：具体异常优先，避免裸 `except`

## 模板

- 使用 Django 模板语法，避免复杂逻辑
- 区块命名：`{% block title %}`、`{% block content %}`
- 变量使用 `{{ variable }}`，必要时 `|default:""`

## 安全

- 用户输入必须过滤/转义
- 敏感操作需 `@login_required` 或权限检查
- 密码、密钥不得硬编码

## 文档

- 模块、函数、类需 docstring
- 复杂逻辑加注释
