# CSS 模块化架构

> VaelorCMS 采用模块化 CSS 架构，便于定制与扩展。

## 模块结构

```
assets/
├── theme.css          # 主入口，按序导入各模块
└── css/
    ├── _variables.css # 变量（颜色、圆角、Aero 参数）
    ├── _base.css      # 重置、字体、body
    ├── _aero.css      # Aero 玻璃特效（body.aero 时增强）
    ├── _layout.css    # 顶栏、主区、页脚
    ├── _components.css# 卡片、按钮、表单、列表
    └── _admin.css     # 管理后台、FAQ、工具类
```

## 定制方式

### 1. 主题覆盖

主题目录 `themes/<id>/style.css` 在 `theme.css` 之后加载，可覆盖任意变量或类。

### 2. 管理员设置

站点设置 → 主题：

- **Aero 玻璃特效**：启用/关闭站点全局玻璃效果
- **Aero 模糊强度**：0–40px，默认 16
- **强调色覆盖**：如 `#8b5a2b`，覆盖 `--accent`
- **自定义 CSS**：注入到 `<head>` 的 CSS 片段

### 3. 变量覆盖

在主题或自定义 CSS 中覆盖 `:root` 变量：

```css
:root {
  --accent: #2563eb;
  --aero-blur: 20px;
}
```

## Aero 特效

当 `aero_enabled=1` 时，`body` 添加 `aero` 类，启用：

- 增强的 `backdrop-filter` 模糊
- 顶栏、卡片、页脚玻璃效果
- 背景径向渐变叠加层

## 扩展新模块

1. 在 `assets/css/` 下新建 `_mymodule.css`
2. 在 `theme.css` 末尾添加 `@import url('css/_mymodule.css');`
