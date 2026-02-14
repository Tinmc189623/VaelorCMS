# 主题开发规范

> 正式规范文档，供主题开发者遵循，确保主题可维护性与扩展性。

## 1. 目录结构规范

```
themes/
└── <theme_id>/           # 主题 ID：小写字母、数字、下划线
    ├── info.ini          # 必填，主题元信息
    ├── style.css         # 必填，主样式
    └── (可选) 其他静态资源
```

## 2. info.ini 规范

| 字段 | 必填 | 说明 | 示例 |
|------|------|------|------|
| name | 是 | 主题显示名称 | 深色主题 |
| description | 是 | 主题描述（简短） | 适合夜间阅读的深色主题 |
| version | 是 | 语义化版本 x.y.z | 1.0.0 |
| author | 否 | 作者名 | 张三 |
| min_cms_version | 否 | 最低 CMS 版本要求 | Demo-26.02.13.20 |

### 示例

```ini
[theme]
name = 深色主题
description = 适合夜间阅读的深色主题
version = 1.0.0
author = VaelorCMS
min_cms_version = Demo-26.02.13.20
```

## 3. CSS 变量规范

主题必须通过覆盖 `:root` 下的 CSS 变量实现样式，不得直接修改核心类名的结构相关属性（如 display、position、flex 布局）。

### 必须支持的变量

| 变量 | 类型 | 说明 |
|------|------|------|
| --fg | color | 主文字色 |
| --fg-muted | color | 次要文字色 |
| --bg | color | 页面背景 |
| --bg-card | color | 卡片背景 |
| --border | color | 边框色 |
| --accent | color | 强调色 |
| --accent-hover | color | 悬停强调色 |

### 可选覆盖变量

| 变量 | 默认值 |
|------|--------|
| --danger | #a63d3d |
| --ok | #4a6b4a |
| --radius | 8px |
| --radius-sm | 4px |
| --aero-bg | rgba(...) |
| --aero-border | rgba(...) |
| --font | 字体族 |

## 4. 禁止事项

- 不得使用 `!important` 覆盖核心布局（颜色、间距除外）
- 不得修改 `.app-bar`、`.main` 等核心容器的 display/flex 结构
- 不得引入外部 CDN 脚本（仅允许字体、图标等只读资源）
- 不得包含恶意代码或追踪脚本

## 5. 目标环境

- 面向现代浏览器（Chrome、Firefox、Safari、Edge 最新版）
- 支持 backdrop-filter、CSS 变量等现代特性

## 6. 版本号

- 遵循语义化版本：主版本.次版本.修订号
- 仅样式微调：修订号 +1
- 新增变量或结构变更：次版本 +1
- 不兼容变更：主版本 +1

## 7. 参考

- [DEV_THEME.md](DEV_THEME.md) - 快速上手
- 内置主题：`themes/default/`、`themes/dark/`
