# VaelorCMS API 文档

## 概述

RESTful JSON 接口，供移动端、第三方、爬虫调用。入口为 `/api/v1/`。**对开发者友好**：无鉴权（公开数据）、UTF-8、标准 JSON。

## 基础

- **Content-Type**：`application/json; charset=utf-8`
- **编码**：UTF-8
- **方法**：GET（只读，无副作用）
- **鉴权**：无，均为公开接口
- **限流**：按 IP 限制请求频率，默认 60 次/分钟，超限返回 429；健康检查接口不限流

## 接口一览

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/stats/ | 获取统计（users、bbs、code、articles） |
| GET | /api/v1/search/?q=关键词 | 搜索（返回 bbs、code、articles） |
| GET | /api/v1/articles/?category=分类 | 文章列表（可选 category 筛选） |
| GET | /api/v1/upgrade/ | 升级检查（预留） |
| GET | /api/v1/health/ | 健康检查（负载均衡/监控探针） |
| GET | /sitemap.xml | 站点地图（SEO） |
| GET | /robots.txt | 爬虫规则（含 sitemap 链接） |

## 示例

### 统计

```
GET /api/v1/stats/
```

响应：
```json
{"users": 10, "bbs": 25, "code": 8}
```

### 搜索

```
GET /api/v1/search/?q=python
```

响应：
```json
{
  "bbs": [{"id": 1, "title": "...", "created_at": "2024-01-01 12:00"}],
  "code": [{"id": 1, "title": "...", "language": "python", "created_at": "2024-01-01 12:00"}]
}
```

### 升级检查

```
GET /api/v1/upgrade/
```

响应：
```json
{"version": "Demo-26.02.13.26", "upgrade_available": false, "message": "当前已是最新版本"}
```

### 健康检查

```
GET /api/v1/health/
```

响应（正常）：
```json
{"status": "ok", "version": "Demo-26.02.13.26", "database": "connected"}
```

响应（数据库异常）：
```json
{"status": "degraded", "version": "Demo-26.02.13.26", "database": "disconnected"}
```
