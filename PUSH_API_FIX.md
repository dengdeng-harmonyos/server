# Push Kit API 修正说明

## 修正日期
2026-01-14

## 修正内容

根据华为 Push Kit HarmonyOS 场景化消息 REST API 官方文档，对推送接口进行了全面修正。

### 1. 修正的关键问题

#### 问题 1: target 字段名称错误
- **错误**: `"pushToken": ["xxx"]`
- **正确**: `"token": ["xxx"]`
- **影响**: 华为服务器无法识别 pushToken 字段，导致推送失败

#### 问题 2: AlertPayload 结构不符合规范
- **错误结构**:
```json
{
  "alert": {
    "title": "标题",
    "body": "内容"
  }
}
```

- **正确结构**:
```json
{
  "notification": {
    "category": "IM",
    "title": "标题",
    "body": "内容",
    "clickAction": {
      "actionType": 0
    }
  }
}
```

#### 问题 3: 缺少必填字段
- `category`: 通知消息类型（必填）
- `clickAction`: 点击行为（必填）

#### 问题 4: FormUpdatePayload 字段不完整
- **缺少的必填字段**:
  - `version`: 卡片版本号
  - `moduleName`: 模块名称
  - `formName`: 卡片名称
  - `abilityName`: Ability 名称

#### 问题 5: FormID 类型错误
- **错误**: `string` 类型
- **正确**: `int64` 类型（范围 0 到 2^63-1）

#### 问题 6: BackgroundPayload 字段名称错误
- **错误**: `"data"`
- **正确**: `"extraData"`

#### 问题 7: PushOptions 字段不完整
- **缺少的重要字段**:
  - `testMessage`: 测试消息标识
  - `biTag`: 批量任务消息标识

### 2. 代码修改摘要

#### 修改的文件

1. **internal/service/huawei_push.go**
   - 修正 `V3Target` 结构体：`PushToken` → `Token`
   - 重构 `AlertPayload`：添加 `Notification` 结构体
   - 添加 `ClickAction` 结构体
   - 重构 `FormUpdatePayload`：添加完整字段
   - 添加 `FormImage` 结构体
   - 重构 `BackgroundPayload`：`Data` → `ExtraData`
   - 添加 `ExtensionPayload`（语音播报）
   - 添加 `VoIPCallPayload`（应用内通话）
   - 重构 `PushOptions`：添加 `testMessage`, `biTag` 等字段
   - 修正 `SendNotification` 方法：构建完整的通知结构
   - 修正 `SendFormUpdate` 方法：添加完整参数
   - 添加 `SendFormUpdateSimple` 方法：简化版卡片刷新
   - 添加 `SendVoIPCall` 方法：应用内通话消息

2. **internal/handler/push.go**
   - 添加 `fmt` 包导入
   - 修正 `SendFormUpdate` 处理器：转换 FormID 为 int64
   - 调用 `SendFormUpdateSimple` 方法

### 3. 新增功能

#### 3.1 应用内通话消息支持
```go
err := pushService.SendVoIPCall(pushToken, extraData)
```

#### 3.2 完整的卡片刷新支持
```go
err := pushService.SendFormUpdate(
    pushToken,
    formID,      // int64
    version,     // int
    moduleName,  // string
    formName,    // string
    abilityName, // string
    formData,    // map[string]interface{}
)
```

#### 3.3 测试消息模式
默认启用测试消息模式（`testMessage: true`），限制：
- 单个项目最多 1000 条测试消息
- 每次推送最多 10 个 Token

### 4. API 请求结构对比

#### 4.1 通知消息（push-type=0）

**修正前**:
```json
{
  "payload": {
    "alert": {
      "title": "标题",
      "body": "内容"
    }
  },
  "target": {
    "pushToken": ["xxx"]
  }
}
```

**修正后**:
```json
{
  "payload": {
    "notification": {
      "category": "IM",
      "title": "标题",
      "body": "内容",
      "clickAction": {
        "actionType": 0
      }
    }
  },
  "target": {
    "token": ["xxx"]
  },
  "pushOptions": {
    "testMessage": true,
    "ttl": 86400
  }
}
```

#### 4.2 卡片刷新消息（push-type=1）

**修正前**:
```json
{
  "payload": {
    "formId": "123",
    "formData": {
      "key": "value"
    }
  },
  "target": {
    "pushToken": ["xxx"]
  }
}
```

**修正后**:
```json
{
  "payload": {
    "formId": 123,
    "version": 922337203,
    "moduleName": "entry",
    "formName": "widget",
    "abilityName": "EntryAbility",
    "formData": {
      "key": "value"
    }
  },
  "target": {
    "token": ["xxx"]
  },
  "pushOptions": {
    "ttl": 666
  }
}
```

#### 4.3 后台消息（push-type=6）

**修正前**:
```json
{
  "payload": {
    "data": "携带的数据"
  },
  "target": {
    "pushToken": ["xxx"]
  }
}
```

**修正后**:
```json
{
  "payload": {
    "extraData": "携带的数据"
  },
  "target": {
    "token": ["xxx"]
  }
}
```

### 5. 兼容性说明

#### 向后兼容
- 保留了原有的 API 端点
- `SendFormUpdateSimple` 提供简化调用方式
- 自动处理 FormID 类型转换

#### 不兼容的变更
- FormID 必须是数字（之前可以是任意字符串）
- 卡片刷新需要提供更多参数（通过简化方法使用默认值）

### 6. 测试建议

1. **通知消息测试**:
```bash
curl "http://localhost:8080/api/v1/push/notification?device_key=YOUR_KEY&title=测试&body=这是一条测试消息"
```

2. **卡片刷新测试**:
```bash
curl "http://localhost:8080/api/v1/push/form?device_key=YOUR_KEY&form_id=12345&form_data={\"content\":\"测试内容\"}"
```

3. **后台消息测试**:
```bash
curl "http://localhost:8080/api/v1/push/background?device_key=YOUR_KEY&data={\"action\":\"sync\"}"
```

### 7. 注意事项

1. **category 选择**:
   - 根据实际业务场景选择合适的 category
   - 不同的 category 影响消息展示和提醒方式
   - 当前默认使用 "IM"（即时聊天）

2. **clickAction 配置**:
   - actionType=0: 打开应用首页（默认）
   - actionType=1: 打开自定义页面（需配置 action 或 uri）
   - actionType=3: 清除通知

3. **测试消息限制**:
   - 生产环境应设置 `testMessage: false`
   - 测试环境使用 `testMessage: true`

4. **卡片版本号管理**:
   - 每次刷新必须递增版本号
   - 版本号为 0 时表示重置

### 8. 下一步工作

- [ ] 配置文件中添加 category 配置项
- [ ] 支持自定义 clickAction
- [ ] 添加图片推送支持
- [ ] 添加角标管理
- [ ] 实现实况窗消息（push-type=7）
- [ ] 添加消息撤回功能
- [ ] 完善错误处理和重试机制

### 9. 参考文档

- [场景化消息请求体结构说明](https://developer.huawei.com/consumer/cn/doc/harmonyos-references/push-scenariozed-api-request-struct)
- [请求体参数说明](https://developer.huawei.com/consumer/cn/doc/harmonyos-references/push-scenariozed-api-request-param)
- [请求示例](https://developer.huawei.com/consumer/cn/doc/harmonyos-references/push-scenariozed-api-request-example)
