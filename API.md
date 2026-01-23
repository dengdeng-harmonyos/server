# DengDeng Push Service API 文档

## 基础信息

- **基础URL**: `http://your-domain.com`
- **内容类型**: `application/json`
- **字符编码**: `UTF-8`
- **API版本**: `v1`

## 统一响应格式

所有API接口都遵循统一的响应格式：

```json
{
  "code": 0,
  "msg": "success",
  "data": {}
}
```

### 响应字段说明

| 字段 | 类型 | 说明 |
|-----|------|------|
| code | int | 响应码，0表示成功，非0表示错误 |
| msg | string | 响应消息，成功时为"success"，失败时为错误描述 |
| data | object/array/null | 响应数据，根据接口不同返回不同类型的数据 |

## 错误码定义

### 系统错误 (1xxx)

| 错误码 | 说明 |
|-------|------|
| 1001 | 系统内部错误 |
| 1002 | 数据库错误 |
| 1003 | 网络错误 |
| 1004 | 服务不可用 |

### 认证错误 (2xxx)

| 错误码 | 说明 |
|-------|------|
| 2001 | 未授权 |
| 2002 | 认证失败 |
| 2003 | Token无效 |
| 2004 | Token过期 |

### 业务错误 (3xxx)

| 错误码 | 说明 |
|-------|------|
| 3001 | 设备未注册 |
| 3002 | 设备已存在 |
| 3003 | 推送失败 |
| 3004 | 消息不存在 |
| 3005 | 操作失败 |

### 数据错误 (4xxx)

| 错误码 | 说明 |
|-------|------|
| 4001 | 参数错误 |
| 4002 | 参数缺失 |
| 4003 | 参数格式错误 |
| 4004 | 数据不存在 |

## HTTP 状态码说明

| 状态码 | 说明 |
|-------|------|
| 200 | 请求成功 |
| 400 | 请求参数错误 |
| 401 | 未授权 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

---

## API 接口详细文档

### 1. 设备注册

#### 接口地址
```
POST /api/v1/device/register
```

#### 请求参数

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| device_id | string | 是 | 设备唯一标识符 |
| push_token | string | 是 | 华为推送Token |
| device_type | string | 否 | 设备类型（如：HarmonyOS） |
| os_version | string | 否 | 操作系统版本 |
| app_version | string | 否 | 应用版本 |

#### 请求示例

```bash
curl -X POST http://your-domain.com/api/v1/device/register \
  -H "Content-Type: application/json" \
  -d '{
    "device_id": "device123",
    "push_token": "token456",
    "device_type": "HarmonyOS",
    "os_version": "4.0",
    "app_version": "1.0.0"
  }'
```

#### 成功响应示例

```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "device_id": "device123",
    "push_token": "token456",
    "status": "active",
    "created_at": "2026-01-23T10:00:00Z"
  }
}
```

#### 错误响应示例

```json
{
  "code": 4001,
  "msg": "device_id is required",
  "data": null
}
```

---

### 2. 更新设备Token

#### 接口地址
```
PUT /api/v1/device/update-token
```

#### 请求参数

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| device_id | string | 是 | 设备唯一标识符 |
| push_token | string | 是 | 新的华为推送Token |

#### 请求示例

```bash
curl -X PUT http://your-domain.com/api/v1/device/update-token \
  -H "Content-Type: application/json" \
  -d '{
    "device_id": "device123",
    "push_token": "new_token789"
  }'
```

#### 成功响应示例

```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "device_id": "device123",
    "push_token": "new_token789",
    "updated_at": "2026-01-23T10:05:00Z"
  }
}
```

#### 错误响应示例

```json
{
  "code": 3001,
  "msg": "device not found",
  "data": null
}
```

---

### 3. 删除设备

#### 接口地址
```
DELETE /api/v1/device/delete
```

#### 请求参数

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| device_id | string | 是 | 设备唯一标识符 |

#### 请求示例

```bash
curl -X DELETE http://your-domain.com/api/v1/device/delete \
  -H "Content-Type: application/json" \
  -d '{
    "device_id": "device123"
  }'
```

#### 成功响应示例

```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "device_id": "device123",
    "deleted_at": "2026-01-23T10:10:00Z"
  }
}
```

#### 错误响应示例

```json
{
  "code": 3001,
  "msg": "device not found",
  "data": null
}
```

#### 注意事项
- 此接口会同时删除设备记录和该设备的所有待接收消息（pending_messages）
- 删除操作不可逆，请谨慎使用

---

### 4. 推送通知

#### 接口地址
```
GET /api/v1/push/notification
```

#### 请求参数

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| device_id | string | 是 | 目标设备ID |
| title | string | 是 | 通知标题 |
| body | string | 是 | 通知内容 |
| click_action | string | 否 | 点击动作 |
| badge_count | int | 否 | 角标数量 |

#### 请求示例

```bash
curl -X GET "http://your-domain.com/api/v1/push/notification?device_id=device123&title=Hello&body=World&badge_count=1"
```

#### 成功响应示例

```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "message_id": "msg_123456",
    "device_id": "device123",
    "push_result": {
      "code": "80000000",
      "msg": "Success",
      "request_id": "req_789"
    },
    "sent_at": "2026-01-23T10:15:00Z"
  }
}
```

#### 错误响应示例

```json
{
  "code": 3003,
  "msg": "push failed: invalid token",
  "data": null
}
```

---

### 5. 推送卡片刷新

#### 接口地址
```
GET /api/v1/push/form
```

#### 请求参数

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| device_id | string | 是 | 目标设备ID |
| form_id | string | 是 | 卡片ID |
| data | string | 是 | 卡片数据（JSON字符串） |

#### 请求示例

```bash
curl -X GET "http://your-domain.com/api/v1/push/form?device_id=device123&form_id=form_001&data=%7B%22title%22%3A%22New%20Content%22%7D"
```

#### 成功响应示例

```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "message_id": "msg_123457",
    "device_id": "device123",
    "form_id": "form_001",
    "push_result": {
      "code": "80000000",
      "msg": "Success",
      "request_id": "req_790"
    },
    "sent_at": "2026-01-23T10:20:00Z"
  }
}
```

#### 错误响应示例

```json
{
  "code": 4003,
  "msg": "invalid form data format",
  "data": null
}
```

---

### 6. 推送后台消息

#### 接口地址
```
GET /api/v1/push/background
```

#### 请求参数

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| device_id | string | 是 | 目标设备ID |
| data | string | 是 | 消息数据（JSON字符串） |

#### 请求示例

```bash
curl -X GET "http://your-domain.com/api/v1/push/background?device_id=device123&data=%7B%22action%22%3A%22sync%22%7D"
```

#### 成功响应示例

```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "message_id": "msg_123458",
    "device_id": "device123",
    "push_result": {
      "code": "80000000",
      "msg": "Success",
      "request_id": "req_791"
    },
    "sent_at": "2026-01-23T10:25:00Z"
  }
}
```

#### 错误响应示例

```json
{
  "code": 3001,
  "msg": "device not found",
  "data": null
}
```

---

### 7. 批量推送

#### 接口地址
```
GET /api/v1/push/batch
```

#### 请求参数

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| device_ids | string | 是 | 设备ID列表（逗号分隔） |
| title | string | 是 | 通知标题 |
| body | string | 是 | 通知内容 |
| badge_count | int | 否 | 角标数量 |

#### 请求示例

```bash
curl -X GET "http://your-domain.com/api/v1/push/batch?device_ids=device123,device456,device789&title=Batch%20Push&body=Hello%20Everyone"
```

#### 成功响应示例

```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "total": 3,
    "success": 2,
    "failed": 1,
    "results": [
      {
        "device_id": "device123",
        "status": "success",
        "message_id": "msg_123459"
      },
      {
        "device_id": "device456",
        "status": "success",
        "message_id": "msg_123460"
      },
      {
        "device_id": "device789",
        "status": "failed",
        "error": "device not found"
      }
    ],
    "sent_at": "2026-01-23T10:30:00Z"
  }
}
```

#### 错误响应示例

```json
{
  "code": 4002,
  "msg": "device_ids is required",
  "data": null
}
```

---

### 8. 推送统计

#### 接口地址
```
GET /api/v1/push/statistics
```

#### 请求参数

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| device_id | string | 否 | 设备ID（不填则返回全局统计） |
| start_time | string | 否 | 开始时间（ISO 8601格式） |
| end_time | string | 否 | 结束时间（ISO 8601格式） |

#### 请求示例

```bash
curl -X GET "http://your-domain.com/api/v1/push/statistics?device_id=device123&start_time=2026-01-20T00:00:00Z&end_time=2026-01-23T23:59:59Z"
```

#### 成功响应示例

```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "device_id": "device123",
    "period": {
      "start": "2026-01-20T00:00:00Z",
      "end": "2026-01-23T23:59:59Z"
    },
    "statistics": {
      "total_sent": 150,
      "total_delivered": 145,
      "total_failed": 5,
      "by_type": {
        "notification": 100,
        "form": 30,
        "background": 20
      },
      "delivery_rate": 96.67
    }
  }
}
```

#### 错误响应示例

```json
{
  "code": 4003,
  "msg": "invalid time format",
  "data": null
}
```

---

### 9. 获取待接收消息

#### 接口地址
```
GET /api/v1/messages/pending
```

#### 请求参数

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| device_id | string | 是 | 设备ID |
| limit | int | 否 | 返回数量限制（默认：20，最大：100） |

#### 请求示例

```bash
curl -X GET "http://your-domain.com/api/v1/messages/pending?device_id=device123&limit=10"
```

#### 成功响应示例

```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "device_id": "device123",
    "total": 5,
    "messages": [
      {
        "message_id": "msg_001",
        "type": "notification",
        "payload": {
          "title": "New Message",
          "body": "You have a new notification"
        },
        "created_at": "2026-01-23T09:00:00Z"
      },
      {
        "message_id": "msg_002",
        "type": "background",
        "payload": {
          "action": "sync",
          "data": "update_required"
        },
        "created_at": "2026-01-23T09:05:00Z"
      }
    ]
  }
}
```

#### 错误响应示例

```json
{
  "code": 4002,
  "msg": "device_id is required",
  "data": null
}
```

---

### 10. 确认消息

#### 接口地址
```
POST /api/v1/messages/confirm
```

#### 请求参数

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| device_id | string | 是 | 设备ID |
| message_ids | array | 是 | 消息ID数组 |

#### 请求示例

```bash
curl -X POST http://your-domain.com/api/v1/messages/confirm \
  -H "Content-Type: application/json" \
  -d '{
    "device_id": "device123",
    "message_ids": ["msg_001", "msg_002"]
  }'
```

#### 成功响应示例

```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "device_id": "device123",
    "confirmed": 2,
    "message_ids": ["msg_001", "msg_002"],
    "confirmed_at": "2026-01-23T10:40:00Z"
  }
}
```

#### 错误响应示例

```json
{
  "code": 3004,
  "msg": "message not found",
  "data": {
    "invalid_ids": ["msg_003"]
  }
}
```

---

### 11. 健康检查

#### 接口地址
```
GET /health
```

#### 请求参数

无

#### 请求示例

```bash
curl -X GET http://your-domain.com/health
```

#### 成功响应示例

```json
{
  "status": "healthy",
  "timestamp": "2026-01-23T10:45:00Z",
  "version": "1.0.0",
  "services": {
    "database": "connected",
    "push_service": "available"
  }
}
```

#### 注意事项
- 此接口不遵循统一响应格式，保持独立的健康检查格式
- 用于监控系统运行状态

---

## 使用示例代码

### JavaScript (Fetch API)

```javascript
// 统一的响应处理函数
async function callAPI(url, options = {}) {
  try {
    const response = await fetch(url, options);
    const result = await response.json();
    
    if (result.code === 0) {
      console.log('Success:', result.msg);
      return result.data;
    } else {
      console.error('Error:', result.msg);
      throw new Error(`API Error ${result.code}: ${result.msg}`);
    }
  } catch (error) {
    console.error('Request failed:', error);
    throw error;
  }
}

// 设备注册示例
async function registerDevice() {
  const data = await callAPI('http://your-domain.com/api/v1/device/register', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      device_id: 'device123',
      push_token: 'token456',
      device_type: 'HarmonyOS',
      os_version: '4.0',
      app_version: '1.0.0'
    })
  });
  
  console.log('Device registered:', data);
  return data;
}

// 推送通知示例
async function sendNotification(deviceId, title, body) {
  const params = new URLSearchParams({
    device_id: deviceId,
    title: title,
    body: body,
    badge_count: '1'
  });
  
  const data = await callAPI(`http://your-domain.com/api/v1/push/notification?${params}`);
  console.log('Notification sent:', data);
  return data;
}

// 获取待接收消息示例
async function getPendingMessages(deviceId) {
  const params = new URLSearchParams({
    device_id: deviceId,
    limit: '20'
  });
  
  const data = await callAPI(`http://your-domain.com/api/v1/messages/pending?${params}`);
  console.log('Pending messages:', data);
  return data;
}
```

### Python (Requests)

```python
import requests
import json

class DengDengAPIClient:
    def __init__(self, base_url):
        self.base_url = base_url
    
    def _handle_response(self, response):
        """统一处理API响应"""
        result = response.json()
        
        if result['code'] == 0:
            print(f"Success: {result['msg']}")
            return result['data']
        else:
            error_msg = f"API Error {result['code']}: {result['msg']}"
            print(error_msg)
            raise Exception(error_msg)
    
    def register_device(self, device_id, push_token, device_type=None, os_version=None, app_version=None):
        """设备注册"""
        url = f"{self.base_url}/api/v1/device/register"
        payload = {
            'device_id': device_id,
            'push_token': push_token
        }
        
        if device_type:
            payload['device_type'] = device_type
        if os_version:
            payload['os_version'] = os_version
        if app_version:
            payload['app_version'] = app_version
        
        response = requests.post(url, json=payload)
        return self._handle_response(response)
    
    def send_notification(self, device_id, title, body, badge_count=None):
        """推送通知"""
        url = f"{self.base_url}/api/v1/push/notification"
        params = {
            'device_id': device_id,
            'title': title,
            'body': body
        }
        
        if badge_count is not None:
            params['badge_count'] = badge_count
        
        response = requests.get(url, params=params)
        return self._handle_response(response)
    
    def get_pending_messages(self, device_id, limit=20):
        """获取待接收消息"""
        url = f"{self.base_url}/api/v1/messages/pending"
        params = {
            'device_id': device_id,
            'limit': limit
        }
        
        response = requests.get(url, params=params)
        return self._handle_response(response)
    
    def confirm_messages(self, device_id, message_ids):
        """确认消息"""
        url = f"{self.base_url}/api/v1/messages/confirm"
        payload = {
            'device_id': device_id,
            'message_ids': message_ids
        }
        
        response = requests.post(url, json=payload)
        return self._handle_response(response)

# 使用示例
if __name__ == '__main__':
    client = DengDengAPIClient('http://your-domain.com')
    
    try:
        # 注册设备
        device_data = client.register_device(
            device_id='device123',
            push_token='token456',
            device_type='HarmonyOS',
            os_version='4.0',
            app_version='1.0.0'
        )
        print(f"Device registered: {device_data}")
        
        # 发送通知
        push_data = client.send_notification(
            device_id='device123',
            title='Hello',
            body='World',
            badge_count=1
        )
        print(f"Notification sent: {push_data}")
        
        # 获取待接收消息
        messages = client.get_pending_messages(device_id='device123', limit=10)
        print(f"Pending messages: {messages}")
        
        # 确认消息
        if messages and messages['messages']:
            message_ids = [msg['message_id'] for msg in messages['messages']]
            confirm_data = client.confirm_messages(device_id='device123', message_ids=message_ids)
            print(f"Messages confirmed: {confirm_data}")
            
    except Exception as e:
        print(f"Error occurred: {e}")
```

### Go

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
)

// 统一响应结构
type APIResponse struct {
    Code int             `json:"code"`
    Msg  string          `json:"msg"`
    Data json.RawMessage `json:"data"`
}

type DengDengClient struct {
    BaseURL string
    Client  *http.Client
}

// 统一处理响应
func (c *DengDengClient) handleResponse(resp *http.Response, result interface{}) error {
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return fmt.Errorf("failed to read response: %w", err)
    }
    
    var apiResp APIResponse
    if err := json.Unmarshal(body, &apiResp); err != nil {
        return fmt.Errorf("failed to parse response: %w", err)
    }
    
    if apiResp.Code != 0 {
        return fmt.Errorf("API error %d: %s", apiResp.Code, apiResp.Msg)
    }
    
    if result != nil && apiResp.Data != nil {
        if err := json.Unmarshal(apiResp.Data, result); err != nil {
            return fmt.Errorf("failed to parse data: %w", err)
        }
    }
    
    return nil
}

// 设备注册
type RegisterDeviceRequest struct {
    DeviceID   string `json:"device_id"`
    PushToken  string `json:"push_token"`
    DeviceType string `json:"device_type,omitempty"`
    OSVersion  string `json:"os_version,omitempty"`
    AppVersion string `json:"app_version,omitempty"`
}

type DeviceInfo struct {
    DeviceID  string `json:"device_id"`
    PushToken string `json:"push_token"`
    Status    string `json:"status"`
    CreatedAt string `json:"created_at"`
}

func (c *DengDengClient) RegisterDevice(req RegisterDeviceRequest) (*DeviceInfo, error) {
    payload, err := json.Marshal(req)
    if err != nil {
        return nil, err
    }
    
    resp, err := c.Client.Post(
        c.BaseURL+"/api/v1/device/register",
        "application/json",
        bytes.NewBuffer(payload),
    )
    if err != nil {
        return nil, err
    }
    
    var device DeviceInfo
    if err := c.handleResponse(resp, &device); err != nil {
        return nil, err
    }
    
    return &device, nil
}

// 推送通知
func (c *DengDengClient) SendNotification(deviceID, title, body string, badgeCount int) (map[string]interface{}, error) {
    params := url.Values{}
    params.Add("device_id", deviceID)
    params.Add("title", title)
    params.Add("body", body)
    if badgeCount > 0 {
        params.Add("badge_count", fmt.Sprintf("%d", badgeCount))
    }
    
    resp, err := c.Client.Get(c.BaseURL + "/api/v1/push/notification?" + params.Encode())
    if err != nil {
        return nil, err
    }
    
    var result map[string]interface{}
    if err := c.handleResponse(resp, &result); err != nil {
        return nil, err
    }
    
    return result, nil
}

// 获取待接收消息
type PendingMessagesResponse struct {
    DeviceID string          `json:"device_id"`
    Total    int             `json:"total"`
    Messages []PendingMessage `json:"messages"`
}

type PendingMessage struct {
    MessageID string                 `json:"message_id"`
    Type      string                 `json:"type"`
    Payload   map[string]interface{} `json:"payload"`
    CreatedAt string                 `json:"created_at"`
}

func (c *DengDengClient) GetPendingMessages(deviceID string, limit int) (*PendingMessagesResponse, error) {
    params := url.Values{}
    params.Add("device_id", deviceID)
    params.Add("limit", fmt.Sprintf("%d", limit))
    
    resp, err := c.Client.Get(c.BaseURL + "/api/v1/messages/pending?" + params.Encode())
    if err != nil {
        return nil, err
    }
    
    var result PendingMessagesResponse
    if err := c.handleResponse(resp, &result); err != nil {
        return nil, err
    }
    
    return &result, nil
}

func main() {
    client := &DengDengClient{
        BaseURL: "http://your-domain.com",
        Client:  &http.Client{},
    }
    
    // 注册设备
    device, err := client.RegisterDevice(RegisterDeviceRequest{
        DeviceID:   "device123",
        PushToken:  "token456",
        DeviceType: "HarmonyOS",
        OSVersion:  "4.0",
        AppVersion: "1.0.0",
    })
    if err != nil {
        fmt.Printf("Register failed: %v\n", err)
        return
    }
    fmt.Printf("Device registered: %+v\n", device)
    
    // 发送通知
    pushResult, err := client.SendNotification("device123", "Hello", "World", 1)
    if err != nil {
        fmt.Printf("Push failed: %v\n", err)
        return
    }
    fmt.Printf("Notification sent: %+v\n", pushResult)
    
    // 获取待接收消息
    messages, err := client.GetPendingMessages("device123", 10)
    if err != nil {
        fmt.Printf("Get messages failed: %v\n", err)
        return
    }
    fmt.Printf("Pending messages: %+v\n", messages)
}
```

---

## 注意事项

1. **统一响应格式**
   - 所有API接口（除健康检查外）都遵循统一的响应格式：`{code, msg, data}`
   - `code` 为 0 表示成功，非 0 表示错误
   - 客户端应始终检查 `code` 字段来判断请求是否成功

2. **设备删除**
   - 删除设备时会自动删除该设备的所有待接收消息
   - 删除操作不可逆，请谨慎使用

3. **消息确认机制**
   - 客户端获取待接收消息后，应及时调用确认接口
   - 未确认的消息会保留在服务器，直到被确认或过期

4. **批量推送**
   - 批量推送接口会返回每个设备的推送结果
   - 部分失败不会影响其他设备的推送

5. **错误处理**
   - 建议在客户端实现统一的错误处理逻辑
   - 根据错误码进行相应的错误处理和用户提示

6. **性能优化**
   - 批量操作优先使用批量接口，避免循环调用单个接口
   - 合理设置 `limit` 参数，避免一次获取过多数据

7. **安全性**
   - 生产环境应使用 HTTPS 协议
   - 建议实现 API 认证机制（如 JWT Token）
   - 敏感数据应进行加密传输

8. **推送Token管理**
   - 推送Token可能会变更，应及时调用更新Token接口
   - 定期检查设备状态，及时清理无效设备

9. **时区处理**
   - 所有时间字段均使用 ISO 8601 格式的 UTC 时间
   - 客户端需要根据本地时区进行转换

10. **API版本**
    - 当前版本为 v1
    - 后续版本更新会保持向后兼容，或通过URL版本号区分
