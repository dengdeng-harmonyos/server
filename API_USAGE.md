# Push Server API 使用说明

本文档基于华为 Push Kit HarmonyOS 场景化消息 REST API 规范编写。

## 目录
- [认证方式](#认证方式)
- [消息类型](#消息类型)
- [API 接口](#api-接口)
- [请求示例](#请求示例)

## 认证方式

服务端使用 **JWT 认证**方式（基于服务账号密钥）：
- 算法：PS256 (SHA256withRSA/PSS)
- 有效期：3600秒（1小时）
- Header 包含：`kid`（密钥ID）、`typ`（JWT）、`alg`（PS256）

## 消息类型

| push-type | 类型 | 说明 |
|-----------|------|------|
| 0 | Alert消息 | 通知消息，展示在通知栏 |
| 1 | 卡片刷新消息 | 用于更新服务卡片内容 |
| 2 | 语音播报消息 | 语音播报类型通知 |
| 6 | 后台消息 | 透传消息，不展示通知 |
| 7 | 实况窗消息 | 实时更新的动态通知 |
| 10 | 应用内通话消息 | VoIP 通话消息 |

## API 接口

### 1. 发送通知消息

**端点**: `GET /api/v1/push/notification`

**参数**:
- `device_key` (必填): 设备唯一标识
- `title` (必填): 通知标题
- `body` (必填): 通知内容
- `data` (可选): JSON 格式的额外数据

**请求体结构** (发送到华为服务器):
```json
{
  "payload": {
    "notification": {
      "category": "IM",
      "title": "消息标题",
      "body": "消息内容",
      "clickAction": {
        "actionType": 0
      }
    }
  },
  "target": {
    "token": ["MAMzL*******"]
  },
  "pushOptions": {
    "testMessage": true,
    "ttl": 86400
  }
}
```

**通知消息类型** (category):
- `IM`: 即时聊天
- `VOIP`: 语音通话邀请、视频通话邀请
- `MISS_CALL`: 未接通话消息提醒
- `SUBSCRIPTION`: 订阅
- `TRAVEL`: 出行
- `HEALTH`: 健康
- `WORK`: 工作事项提醒
- `ACCOUNT`: 账号动态
- `EXPRESS`: 订单&物流
- `FINANCE`: 财务
- `DEVICE_REMINDER`: 设备提醒
- `MAIL`: 邮件

**点击行为** (clickAction.actionType):
- `0`: 打开应用首页
- `1`: 打开应用自定义页面（需配合 action 或 uri）
- `3`: 清除通知
- `5`: 打开拨号界面（需配合 data.tel）

### 2. 发送卡片刷新消息

**端点**: `GET /api/v1/push/form`

**参数**:
- `device_key` (必填): 设备唯一标识
- `form_id` (必填): 卡片实例ID（数字）
- `form_data` (必填): JSON 格式的卡片数据

**请求体结构** (发送到华为服务器):
```json
{
  "payload": {
    "formId": 0,
    "version": 922337203,
    "moduleName": "entry",
    "formName": "widget",
    "abilityName": "EntryAbility",
    "formData": {
      "key": "value"
    },
    "images": [
      {
        "keyName": "icon",
        "url": "https://example.com/image.png",
        "require": 1
      }
    ]
  },
  "target": {
    "token": ["MAMzL*******"]
  },
  "pushOptions": {
    "ttl": 666
  }
}
```

**说明**:
- `formId`: 卡片实例ID，范围 [0, 2^63-1]
- `version`: 卡片版本号，新版本号必须大于旧版本号
- `moduleName`: 模块名称（module.json5 中的 module.name）
- `formName`: 卡片名称（form 配置中的 name）
- `abilityName`: 卡片 Ability 名称
- `formData`: 待刷新的卡片数据，key-value 格式
- `images`: 可选，卡片图片数组

### 3. 发送后台消息

**端点**: `GET /api/v1/push/background`

**参数**:
- `device_key` (必填): 设备唯一标识
- `data` (必填): 传递给应用的数据

**请求体结构** (发送到华为服务器):
```json
{
  "payload": {
    "extraData": "携带的数据"
  },
  "target": {
    "token": ["MAMzL*******"]
  }
}
```

**说明**:
- 后台消息不会展示通知，数据直接传递给应用
- 应用在前台时直接传递，不在前台时缓存或静默写入
- 可选 `proxyData: "ENABLE"` 开启数据代理

### 4. 发送应用内通话消息

**说明**: 需要调用新增的 `SendVoIPCall` 方法

```go
err := pushService.SendVoIPCall(pushToken, extraData)
```

**请求体结构** (发送到华为服务器):
```json
{
  "payload": {
    "extraData": "传递给应用的数据"
  },
  "target": {
    "token": ["MAMzL*******"]
  },
  "pushOptions": {
    "ttl": 30
  }
}
```

**说明**:
- 用于应用内通话场景
- 建议 TTL 设置为 30-60 秒
- 应用根据 extraData 自行处理通话逻辑

### 5. 批量发送通知消息

**端点**: `GET /api/v1/push/batch`

**参数**:
- `device_keys` (必填): 逗号分隔的设备标识列表
- `title` (必填): 通知标题
- `body` (必填): 通知内容
- `data` (可选): JSON 格式的额外数据

**限制**:
- 单次最多 1000 个 token
- 卡片刷新场景单次只允许 1 个 token

## 请求示例

### cURL 示例

#### 1. 发送通知消息
```bash
curl "http://localhost:8080/api/v1/push/notification?device_key=abc123&title=新消息&body=您有一条新消息&data={\"type\":\"chat\",\"id\":\"123\"}"
```

#### 2. 发送卡片刷新
```bash
curl "http://localhost:8080/api/v1/push/form?device_key=abc123&form_id=12345&form_data={\"content\":\"更新的内容\"}"
```

#### 3. 发送后台消息
```bash
curl "http://localhost:8080/api/v1/push/background?device_key=abc123&data={\"action\":\"sync\",\"timestamp\":1234567890}"
```

#### 4. 批量发送
```bash
curl "http://localhost:8080/api/v1/push/batch?device_keys=abc123,def456,ghi789&title=系统通知&body=系统维护通知"
```

### Postman 示例

已包含在 `postman/Dengdeng_Push_API.postman_collection.json` 文件中。

## 响应格式

### 成功响应
```json
{
  "success": true,
  "message": "Notification sent successfully"
}
```

### 错误响应
```json
{
  "success": false,
  "error": "错误描述"
}
```

## 华为 Push Kit 错误码

| 错误码 | 说明 |
|--------|------|
| 80000000 | 成功 |
| 80100000 | 部分成功 |
| 80100013 | Token 无效 |
| 80300002 | 参数无效 |
| 80300007 | TTL 参数超出范围 |
| 80300010 | Token 数量超出限制 |

## 注意事项

1. **测试消息限制**: 
   - 单个项目最多推送 1000 条测试消息（testMessage=true）
   - 每次推送携带 Token 数不超过 10 个

2. **消息大小限制**:
   - 整体消息体不超过 4KB
   - 图片 URL 长度限制

3. **频率控制**:
   - 建议合理控制推送频率
   - 避免短时间内大量推送

4. **Token 管理**:
   - Token 可能会过期或失效
   - 建议定期刷新 Token

5. **卡片刷新**:
   - 单次只能携带 1 个 Token
   - 版本号必须递增

## 参考文档

- [华为 Push Kit 场景化消息 REST API](https://developer.huawei.com/consumer/cn/doc/harmonyos-references/push-scenariozed-api-request-struct)
- [请求体参数说明](https://developer.huawei.com/consumer/cn/doc/harmonyos-references/push-scenariozed-api-request-param)
- [请求示例](https://developer.huawei.com/consumer/cn/doc/harmonyos-references/push-scenariozed-api-request-example)
