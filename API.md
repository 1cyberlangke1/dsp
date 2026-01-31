# API 文档

本文档详细介绍 DS2API 提供的所有 API 端点。

## 基础信息

- **Base URL**: `https://your-domain.com` 或 `http://localhost:5001`
- **认证方式**: Bearer Token
- **响应格式**: JSON

## OpenAI 兼容接口

### 获取模型列表

```http
GET /v1/models
```

**响应示例**:

```json
{
  "object": "list",
  "data": [
    {"id": "deepseek-chat", "object": "model", "owned_by": "deepseek"},
    {"id": "deepseek-reasoner", "object": "model", "owned_by": "deepseek"},
    {"id": "deepseek-chat-search", "object": "model", "owned_by": "deepseek"},
    {"id": "deepseek-reasoner-search", "object": "model", "owned_by": "deepseek"}
  ]
}
```

---

### 对话补全

```http
POST /v1/chat/completions
Authorization: Bearer your-api-key
Content-Type: application/json
```

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|-----|------|------|------|
| `model` | string | ✅ | 模型名称 |
| `messages` | array | ✅ | 对话消息列表 |
| `stream` | boolean | ❌ | 是否流式输出，默认 `false` |
| `temperature` | number | ❌ | 温度参数，0-2 |
| `max_tokens` | number | ❌ | 最大输出 token 数 |

**支持的模型**:

| 模型 | thinking_enabled | search_enabled | 说明 |
|-----|-----------------|----------------|------|
| `deepseek-chat` | ❌ | ❌ | 标准对话 |
| `deepseek-reasoner` | ✅ | ❌ | 深度思考（R1 推理） |
| `deepseek-chat-search` | ❌ | ✅ | 联网搜索 |
| `deepseek-reasoner-search` | ✅ | ✅ | 深度思考 + 联网搜索 |

**请求示例**:

```json
{
  "model": "deepseek-reasoner-search",
  "messages": [
    {"role": "system", "content": "你是一个有帮助的助手。"},
    {"role": "user", "content": "今天有什么重要新闻？"}
  ],
  "stream": true
}
```

**流式响应格式** (`stream: true`):

```
data: {"id":"...","object":"chat.completion.chunk","choices":[{"delta":{"role":"assistant"},"index":0}]}

data: {"id":"...","object":"chat.completion.chunk","choices":[{"delta":{"reasoning_content":"让我思考一下..."},"index":0}]}

data: {"id":"...","object":"chat.completion.chunk","choices":[{"delta":{"content":"根据搜索结果..."},"index":0}]}

data: {"id":"...","object":"chat.completion.chunk","choices":[{"index":0,"finish_reason":"stop"}]}

data: [DONE]
```

**非流式响应格式** (`stream: false`):

```json
{
  "id": "chatcmpl-xxx",
  "object": "chat.completion",
  "created": 1699000000,
  "model": "deepseek-reasoner",
  "choices": [{
    "index": 0,
    "message": {
      "role": "assistant",
      "content": "回复内容",
      "reasoning_content": "思考过程（仅 reasoner 模型）"
    },
    "finish_reason": "stop"
  }],
  "usage": {
    "prompt_tokens": 10,
    "completion_tokens": 50,
    "total_tokens": 60
  }
}
```

---

## 管理接口

所有管理接口需要在请求头中携带 `Authorization: Bearer <DS2API_ADMIN_KEY>`。

### 登录

```http
POST /admin/login
Content-Type: application/json
```

**请求体**:

```json
{
  "key": "your-admin-key"
}
```

**响应**:

```json
{
  "success": true,
  "token": "jwt-token",
  "expires_in": 86400
}
```

---

### 获取配置

```http
GET /admin/config
Authorization: Bearer <jwt-token>
```

**响应**:

```json
{
  "keys": ["api-key-1", "api-key-2"],
  "accounts": [
    {
      "email": "user@example.com",
      "password": "***",
      "token": "session-token"
    }
  ]
}
```

---

### 添加账号

```http
POST /admin/accounts
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**请求体**:

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

---

### 批量导入账号

```http
POST /admin/accounts/batch
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**请求体**:

```json
{
  "accounts": [
    {"email": "user1@example.com", "password": "pass1"},
    {"email": "user2@example.com", "password": "pass2"}
  ]
}
```

---

### 测试账号

```http
POST /admin/accounts/test
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**请求体**:

```json
{
  "email": "user@example.com"
}
```

---

### 批量测试所有账号

```http
POST /admin/accounts/test-all
Authorization: Bearer <jwt-token>
```

---

### 获取队列状态

```http
GET /admin/queue/status
Authorization: Bearer <jwt-token>
```

**响应**:

```json
{
  "total_accounts": 5,
  "healthy_accounts": 4,
  "queue_size": 10,
  "accounts": [
    {
      "email": "user@example.com",
      "status": "healthy",
      "last_used": "2026-02-01T06:00:00Z"
    }
  ]
}
```

---

### 同步到 Vercel

```http
POST /admin/vercel/sync
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**请求体**:

```json
{
  "vercel_token": "your-vercel-token",
  "project_id": "your-project-id"
}
```

---

## 错误处理

所有错误响应遵循以下格式：

```json
{
  "error": {
    "message": "错误描述",
    "type": "error_type",
    "code": "error_code"
  }
}
```

**常见错误码**:

| HTTP 状态码 | 错误类型 | 说明 |
|------------|---------|------|
| 400 | `invalid_request_error` | 请求参数错误 |
| 401 | `authentication_error` | API Key 无效或未提供 |
| 403 | `permission_denied` | 权限不足 |
| 429 | `rate_limit_error` | 请求过于频繁 |
| 500 | `internal_error` | 服务器内部错误 |
| 503 | `service_unavailable` | 无可用账号 |

---

## 使用示例

### Python

```python
import openai

client = openai.OpenAI(
    api_key="your-api-key",
    base_url="https://your-domain.com/v1"
)

response = client.chat.completions.create(
    model="deepseek-reasoner",
    messages=[{"role": "user", "content": "你好"}],
    stream=True
)

for chunk in response:
    if chunk.choices[0].delta.content:
        print(chunk.choices[0].delta.content, end="")
```

### cURL

```bash
curl https://your-domain.com/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-api-key" \
  -d '{
    "model": "deepseek-chat",
    "messages": [{"role": "user", "content": "你好"}]
  }'
```

### JavaScript

```javascript
const response = await fetch('https://your-domain.com/v1/chat/completions', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': 'Bearer your-api-key'
  },
  body: JSON.stringify({
    model: 'deepseek-chat-search',
    messages: [{ role: 'user', content: '今天有什么新闻？' }],
    stream: true
  })
});

const reader = response.body.getReader();
const decoder = new TextDecoder();

while (true) {
  const { done, value } = await reader.read();
  if (done) break;
  console.log(decoder.decode(value));
}
```
