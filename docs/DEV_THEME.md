# 主题开发指南

> 写给想给 VaelorCMS 做主题的开发者。看完这篇，你就能做出自己的主题。

---

## 一、主题是啥

主题就是一套 CSS（外加可选的静态资源），用来改变站点的外观。VaelorCMS 的主题系统很简单：把主题放到 `themes/<主题ID>/` 目录，在管理后台选一下，就能生效。

**你能改的：**

- 颜色、字体、圆角、阴影
- 布局微调（通过覆盖类）
- 深色/浅色、各种风格

**你不能改的（目前）：**

- HTML 结构（模板是固定的）
- 路由、业务逻辑

---

## 二、主题目录结构

```
themes/
└── my_theme/
    ├── info.ini      # 主题信息（必填）
    ├── style.css     # 主样式（必填）
    └── (可选) 其他图片、字体等
```

### info.ini

```ini
[theme]
name = 我的主题
description = 一个 xxx 风格的主题
version = 1.0
author = 你的名字
```

### style.css

主题的样式文件。会**追加**在默认 `theme.css` 之后加载，所以你可以用覆盖的方式改样式。

**示例（深色主题）：**

```css
:root {
  --fg: #e8e6e3;
  --fg-muted: #a8a29e;
  --bg: #1c1917;
  --bg-card: #292524;
  --border: #44403c;
  --accent: #d6d3d1;
  --accent-hover: #fafaf9;
  --aero-bg: rgba(41, 37, 36, 0.9);
  --aero-border: rgba(68, 64, 60, 0.6);
}
body {
  background: linear-gradient(135deg, #1c1917 0%, #292524 100%);
}
```

---

## 三、可用的 CSS 变量

默认主题（`assets/theme.css`）定义了这些变量，你可以在自己的主题里覆盖：

| 变量 | 说明 | 默认值 |
|------|------|--------|
| --fg | 主文字色 | #2a2520 |
| --fg-muted | 次要文字色 | #6b635c |
| --bg | 页面背景 | #f6f3ef |
| --bg-card | 卡片背景 | #fefcf8 |
| --border | 边框色 | #ddd8d0 |
| --accent | 强调色（链接、按钮） | #8b5a2b |
| --accent-hover | 悬停强调色 | #a07040 |
| --danger | 错误/危险色 | #a63d3d |
| --ok | 成功色 | #4a6b4a |
| --radius | 圆角 | 8px |
| --radius-sm | 小圆角 | 4px |
| --aero-bg | Aero 玻璃背景 | rgba(254,252,248,0.85) |
| --aero-border | Aero 玻璃边框 | rgba(221,216,208,0.6) |
| --font | 字体族 | Noto Serif SC, ... |

---

## 四、主要类名（方便你覆盖）

| 类名 | 用途 |
|------|------|
| .app-bar | 顶栏 |
| .nav-pills | 导航按钮区 |
| .main | 主内容区 |
| .card | 卡片容器 |
| .card-title | 卡片标题 |
| .card-body | 卡片内容 |
| .tile | 首页瓦片链接 |
| .list-item | 列表项 |
| .btn | 按钮 |
| .input-group | 表单组 |
| .alert-error / .alert-ok | 提示框 |

---

## 五、发布主题

1. 把主题目录打包成 zip
2. 别人解压到 `themes/` 下
3. 在管理后台「站点设置」→「主题」里选择

没有应用商店，没有自动更新，就是文件拷贝。简单直接。

---

## 六、参考

- 内置深色主题：`themes/dark/`
- 默认主题变量：`assets/theme.css`
