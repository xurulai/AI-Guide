# API接口文档

## 对话接口

### POST /api/v1/chat

阻塞式对话接口

**请求示例：**
```json
{
  "session_id": "session-uuid-123",
  "query": "我想买一辆自行车",
  "user_id": "user-123"
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "session_id": "session-uuid-123",
    "response": "好的！为您推荐几款适合的自行车...",
    "tool_used": "PRODUCT_RECOMMENDATION_MODULE",
    "recommended_products": [],
    "metadata": {}
  }
}
```

### POST /api/v1/chat/stream

流式对话接口

## 会话接口

### POST /api/v1/sessions

创建会话

### GET /api/v1/sessions/:session_id

获取会话详情

### GET /api/v1/sessions

获取用户会话列表

### DELETE /api/v1/sessions/:session_id

删除会话

## 内部接口

### POST /internal/products/search

商品检索（供Dify调用）

## 管理接口

### GET /health

健康检查

### GET /metrics

监控指标

### GET /admin/dify/workflows/status

Dify工作流状态

