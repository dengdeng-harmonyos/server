# API 接口文档

本文档详细描述了噔噔推送服务的所有 API 接口。

## 基础信息

- **基础URL**: `http://your-server:8080`
- **内容类型**: `application/json` (POST 请求)
- **字符编码**: `UTF-8`

## 接口列表

### 1. 设备注册

注册设备并获取 Device Key，用于后续的推送操作。

**接口地址**：`POST /api/v1/device/register`

**请求头**：
```
Content-Type: application/json
```

**请求参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| push_token | string | 是 | 华为推送服务返回的 Push Token |
| public_key | string | 否 | RSA 公钥(PEM格式)，用于客户端消息加密 |
| device_type | string | 否 | 设备类型：phone/tablet/watch |
| os_version | string | 否 | HarmonyOS 版本号 |
| app_version | string | 否 | 应用版本号 |

**请求示例**：

```bash
curl -X POST http://localhost:8080/api/v1/device/register \
  -H "Content-Type: application/json" \
  -d '{
    "push_token": "your_huawei_push_token",
    "public_key": "-----BEGIN PUBLIC KEY-----\n...\n-----END PUBLIC KEY-----",
    "device_type": "phone",
    "os_version": "5.0.0",
    "app_version": "1.0.0"
  }'
```

**响应示例**：

```json
{
  "success": true,
  "device_key": "abc123def456...",
  "message": "Device registered successfully"
}
```

**响应字段说明**：

| 字段名 | 类型 | 说明 |
|--------|------|------|
| success | boolean | 注册是否成功 |
| device_key | string | 设备密钥，用于后续推送请求 |
| message | string | 响应消息 |

---

### 2. 删除设备

删除设备并清除相关的推送配置和待发送消息。

**接口地址**：`DELETE /api/v1/device/delete`

**请求参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| device_key | string | 是 | 设备密钥（注册时返回） |

**请求示例**：

```bash
curl -X DELETE "http://localhost:8080/api/v1/device/delete?device_key=abc123def456..."
```

**响应示例**：

```json
{
  "success": true,
  "message": "Device deleted successfully"
}
```

**响应字段说明**：

| 字段名 | 类型 | 说明 |
|--------|------|------|
| success | boolean | 删除是否成功 |
| message | string | 响应消息 |

**说明**：
- 会自动级联删除该设备的所有待发送消息（pending_messages）
- 删除操作不可逆，请谨慎操作
- 注意：由于华为 Push Kit 不提供 Token 删除接口，设备的 Push Token 在华为侧仍然有效，但本地已无法使用该 device_key 进行推送

---

### 3. 推送通知消息

发送通知栏消息到指定设备。

**接口地址**：`GET /api/v1/push/notification`

**请求参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| device_key | string | 是 | 设备密钥（注册时返回） |
| title | string | 是 | 通知标题 |
| body | string | 是 | 通知内容 |
| data | string | 否 | 附加数据（JSON字符串） |

**请求示例**：

```bash
curl "http://localhost:8080/api/push/notification?device_key=abc123&title=Hello&body=World&data=%7B%22key%22%3A%22value%22%7D"
```

**响应示例**：

```json
{
  "success": true,
  "message": "Push notification sent successfully",
  "request_id": "huawei_request_id_xxx"
}
```

**响应字段说明**：

| 字段名 | 类型 | 说明 |
|--------|------|------|
| success | boolean | 推送是否成功 |
| message | string | 响应消息 |
| request_id | string | 华为推送服务返回的请求ID |

---

### 4. 卡片刷新

刷新 HarmonyOS 卡片(Form)内容。

**接口地址**：`GET /api/v1/push/form`

**请求参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| device_key | string | 是 | 设备密钥（注册时返回） |
| form_id | string | 是 | 卡片ID |
| form_data | string | 是 | 卡片数据（JSON字符串） |

**请求示例**：

```bash
curl "http://localhost:8080/api/push/form?device_key=abc123&form_id=form_001&form_data=%7B%22temperature%22%3A%2225%C2%B0C%22%7D"
```

**响应示例**：

```json
{
  "success": true,
  "message": "Form update sent successfully",
  "request_id": "huawei_request_id_xxx"
}
```

---

### 5. 后台推送

发送后台数据推送，应用在后台也能接收。

**接口地址**：`GET /api/v1/push/background`

**请求参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| device_key | string | 是 | 设备密钥（注册时返回） |
| data | string | 是 | 推送数据（JSON字符串） |

**请求示例**：

```bash
curl "http://localhost:8080/api/push/background?device_key=abc123&data=%7B%22action%22%3A%22sync%22%2C%22timestamp%22%3A1234567890%7D"
```

**响应示例**：

```json
{
  "success": true,
  "message": "Background push sent successfully",
  "request_id": "huawei_request_id_xxx"
}
```

---

### 6. 批量推送

同时向多个设备推送通知消息。

**接口地址**：`GET /api/v1/push/batch`

**请求参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| device_keys | string | 是 | 设备密钥列表（逗号分隔） |
| title | string | 是 | 通知标题 |
| body | string | 是 | 通知内容 |
| data | string | 否 | 附加数据（JSON字符串） |

**请求示例**：

```bash
curl "http://localhost:8080/api/push/batch?device_keys=abc123,def456,ghi789&title=群发消息&body=这是一条群发消息"
```

**响应示例**：

```json
{
  "success": true,
  "message": "Batch push completed",
  "results": [
    {
      "device_key": "abc123",
      "success": true,
      "request_id": "huawei_request_id_xxx"
    },
    {
      "device_key": "def456",
      "success": true,
      "request_id": "huawei_request_id_yyy"
    },
    {
      "device_key": "ghi789",
      "success": false,
      "error": "Device not found"
    }
  ]
}
```

**响应字段说明**：

| 字段名 | 类型 | 说明 |
|--------|------|------|
| success | boolean | 批量推送是否全部成功 |
| message | string | 响应消息 |
| results | array | 每个设备的推送结果 |
| results[].device_key | string | 设备密钥 |
| results[].success | boolean | 该设备推送是否成功 |
| results[].request_id | string | 华为推送请求ID（成功时） |
| results[].error | string | 错误信息（失败时） |

---

### 7. 健康检查

检查服务是否正常运行。

**接口地址**：`GET /health`

**请求示例**：

```bash
curl http://localhost:8080/health
```

**响应示例**：

```json
{
  "status": "ok",
  "database": "connected",
  "timestamp": "2026-01-21T10:30:00Z"
}
```

---

## 错误码说明

### HTTP 状态码

| 状态码 | 说明 |
|--------|------|
| 200 | 请求成功 |
| 400 | 请求参数错误 |
| 404 | 资源不存在 |
| 429 | 请求过于频繁，触发速率限制 |
| 500 | 服务器内部错误 |

### 业务错误码

错误响应格式：

```json
{
  "success": false,
  "error": "错误描述",
  "code": "ERROR_CODE"
}
```

常见错误码：

| 错误码 | 说明 |
|--------|------|
| INVALID_DEVICE_KEY | Device Key 无效或已过期 |
| DEVICE_NOT_FOUND | 设备未找到 |
| PUSH_TOKEN_INVALID | Push Token 无效 |
| RATE_LIMIT_EXCEEDED | 超出速率限制 |
| ENCRYPTION_ERROR | 加密/解密错误 |
| HUAWEI_PUSH_ERROR | 华为推送服务错误 |
| DATABASE_ERROR | 数据库错误 |

---

## 使用示例

### JavaScript/TypeScript

```javascript
// 注册设备
async function registerDevice(pushToken) {
  const response = await fetch('http://localhost:8080/api/device/register', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      push_token: pushToken,
      device_type: 'phone',
      os_version: '5.0.0',
      app_version: '1.0.0'
    })
  });
  const data = await response.json();
  return data.device_key;
}

// 推送通知
async function sendNotification(deviceKey, title, body) {
  const params = new URLSearchParams({
    device_key: deviceKey,
    title: title,
    body: body
  });
  const response = await fetch(`http://localhost:8080/api/push/notification?${params}`);
  return await response.json();
}
```

### Python

```python
import requests
import json

# 注册设备
def register_device(push_token):
    response = requests.post(
        'http://localhost:8080/api/device/register',
        json={
            'push_token': push_token,
            'device_type': 'phone',
            'os_version': '5.0.0',
            'app_version': '1.0.0'
        }
    )
    return response.json()['device_key']

# 推送通知
def send_notification(device_key, title, body):
    params = {
        'device_key': device_key,
        'title': title,
        'body': body
    }
    response = requests.get(
        'http://localhost:8080/api/push/notification',
        params=params
    )
    return response.json()
```

### Go

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
)

type RegisterRequest struct {
    PushToken  string `json:"push_token"`
    DeviceType string `json:"device_type"`
    OSVersion  string `json:"os_version"`
    AppVersion string `json:"app_version"`
}

type RegisterResponse struct {
    Success   bool   `json:"success"`
    DeviceKey string `json:"device_key"`
    Message   string `json:"message"`
}

// 注册设备
func registerDevice(pushToken string) (string, error) {
    req := RegisterRequest{
        PushToken:  pushToken,
        DeviceType: "phone",
        OSVersion:  "5.0.0",
        AppVersion: "1.0.0",
    }
    
    jsonData, _ := json.Marshal(req)
    resp, err := http.Post(
        "http://localhost:8080/api/device/register",
        "application/json",
        bytes.NewBuffer(jsonData),
    )
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    
    var result RegisterResponse
    json.NewDecoder(resp.Body).Decode(&result)
    return result.DeviceKey, nil
}

// 推送通知
func sendNotification(deviceKey, title, body string) error {
    params := url.Values{}
    params.Add("device_key", deviceKey)
    params.Add("title", title)
    params.Add("body", body)
    
    resp, err := http.Get(
        "http://localhost:8080/api/push/notification?" + params.Encode(),
    )
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    return nil
}
```

---

## 注意事项

1. **安全性**
   - 生产环境请使用 HTTPS
   - 妥善保管 Device Key，不要泄露
   - 建议实现 API 鉴权机制

2. **速率限制**
   - 默认每设备每日最大推送数为 1000
   - 超出限制将返回 429 错误
   - 可通过环境变量 `MAX_DAILY_PUSH_PER_DEVICE` 调整

3. **数据格式**
   - GET 请求参数中的 JSON 数据需要 URL 编码
   - POST 请求使用标准 JSON 格式
   - 所有字符串使用 UTF-8 编码

4. **Device Key 有效期**
   - 默认有效期为 1 年
   - 过期后需要重新注册
   - 可通过环境变量 `DEVICE_KEY_TTL` 调整

5. **错误处理**
   - 所有接口都应检查 `success` 字段
   - 失败时查看 `error` 字段获取详细信息
   - 建议实现重试机制

---

## 更新日志

### v1.0.0 (2026-01-21)
- 初始版本发布
- 支持基本的设备注册和推送功能
- 支持通知、卡片、后台推送
- 支持批量推送

---

如有问题或建议，请访问 [GitHub Issues](https://github.com/dengdeng-harmenyos/server/issues)。
