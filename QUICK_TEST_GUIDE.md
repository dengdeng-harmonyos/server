# Push Kit 快速测试指南

## 准备工作

1. 确保已配置服务账号密钥文件 `config/service-account.json`
2. 确保Docker服务正在运行

## 启动服务

### 方式一：使用脚本（推荐）
```bash
./rebuild.sh
```

### 方式二：手动启动
```bash
# 停止旧容器
docker-compose -f docker-compose.single.yml down

# 重新构建并启动
docker-compose -f docker-compose.single.yml up --build -d

# 查看日志
docker-compose -f docker-compose.single.yml logs -f dangdangdang-all
```

## 测试流程

### 1. 注册设备（获取device_key）

使用华为设备上的应用获取 Push Token，然后注册到服务器：

```bash
curl -X POST http://localhost:8080/api/v1/device/register \
  -H "Content-Type: application/json" \
  -d '{
    "push_token": "YOUR_HUAWEI_PUSH_TOKEN",
    "device_type": "phone",
    "os_version": "HarmonyOS 5.0",
    "app_version": "1.0.0"
  }'
```

**响应示例**:
```json
{
  "success": true,
  "device_key": "abc123def456",
  "message": "Device registered successfully"
}
```

**保存 device_key，后续测试需要使用！**

### 2. 测试通知消息

#### 基础通知
```bash
curl "http://localhost:8080/api/v1/push/notification?device_key=abc123def456&title=测试通知&body=这是一条测试消息"
```

#### 带数据的通知
```bash
curl "http://localhost:8080/api/v1/push/notification?device_key=abc123def456&title=新消息&body=您收到一条新消息&data={\"type\":\"chat\",\"sender\":\"张三\",\"messageId\":\"msg_001\"}"
```

**预期结果**:
- 华为设备收到通知
- 通知栏显示标题和内容
- 点击通知打开应用

### 3. 测试卡片刷新

**前提**: 应用已添加服务卡片到桌面

```bash
curl "http://localhost:8080/api/v1/push/form?device_key=abc123def456&form_id=12345&form_data={\"title\":\"更新的标题\",\"content\":\"更新的内容\",\"count\":42}"
```

**注意**:
- `form_id` 是卡片实例ID（从 onAddForm 获取）
- `form_data` 的 key 需要与卡片配置中的变量对应

**预期结果**:
- 桌面卡片内容实时更新

### 4. 测试后台消息

```bash
curl "http://localhost:8080/api/v1/push/background?device_key=abc123def456&data={\"action\":\"sync\",\"timestamp\":1704816000,\"items\":[\"item1\",\"item2\"]}"
```

**预期结果**:
- 不显示通知
- 应用在前台时收到数据
- 应用在后台时数据被缓存

### 5. 测试批量推送

```bash
# 首先注册多个设备，获取多个 device_key
# 然后批量发送

curl "http://localhost:8080/api/v1/push/batch?device_keys=key1,key2,key3&title=系统通知&body=系统将在30分钟后维护"
```

**预期结果**:
- 所有指定设备收到相同的通知

## 查看日志

### 实时查看所有日志
```bash
./view-logs.sh
# 选择选项 1: 查看所有日志（实时）
```

### 查看错误日志
```bash
docker-compose -f docker-compose.single.yml logs dangdangdang-all | grep "ERROR"
```

### 查看推送相关日志
```bash
docker-compose -f docker-compose.single.yml logs dangdangdang-all | grep "Push"
```

## 常见问题排查

### 问题 1: 推送失败，提示 Token 无效

**原因**: Push Token 过期或无效

**解决**:
1. 在设备上重新获取 Push Token
2. 调用更新接口更新 Token
```bash
curl -X PUT http://localhost:8080/api/v1/device/YOUR_DEVICE_KEY/token \
  -H "Content-Type: application/json" \
  -d '{"push_token": "NEW_PUSH_TOKEN"}'
```

### 问题 2: 推送失败，提示参数无效

**原因**: 请求参数格式错误

**检查**:
1. JSON 格式是否正确
2. URL 编码是否正确
3. 必填参数是否都提供了

**示例**: JSON 数据需要 URL 编码
```bash
# 正确的做法
DATA=$(echo '{"type":"chat"}' | jq -sRr @uri)
curl "http://localhost:8080/api/v1/push/notification?device_key=abc&title=test&body=test&data=$DATA"
```

### 问题 3: 服务启动失败

**检查清单**:
1. 查看日志: `docker-compose -f docker-compose.single.yml logs`
2. 检查服务账号文件是否存在且格式正确
3. 检查端口是否被占用: `lsof -i :8080`

### 问题 4: 收不到推送

**排查步骤**:
1. 检查设备 Token 是否正确
2. 查看服务端日志确认推送是否发送成功
3. 检查设备网络连接
4. 确认应用通知权限已开启
5. 检查是否是测试消息超出限制（最多1000条）

## 性能测试

### 单设备推送性能
```bash
# 发送100条通知
for i in {1..100}; do
  curl -s "http://localhost:8080/api/v1/push/notification?device_key=abc&title=测试$i&body=内容$i" &
done
wait
```

### 批量推送性能
```bash
# 准备1000个设备key（逗号分隔）
KEYS="key1,key2,key3,...,key1000"
time curl "http://localhost:8080/api/v1/push/batch?device_keys=$KEYS&title=压力测试&body=批量推送测试"
```

## 监控指标

### 查看推送统计
```bash
curl http://localhost:8080/api/v1/stats/push
```

**响应示例**:
```json
{
  "total_pushes": 1000,
  "success_count": 980,
  "failed_count": 20,
  "success_rate": 98.0,
  "by_type": {
    "notification": 800,
    "form": 150,
    "background": 50
  }
}
```

### 查看活跃设备数
```bash
curl http://localhost:8080/api/v1/stats/devices
```

## 集成到应用

### Android/HarmonyOS 端集成

```kotlin
// 1. 获取 Push Token
PushClient.getToken()
    .addOnSuccessListener { token ->
        // 2. 注册到服务器
        registerDevice(token)
    }

fun registerDevice(pushToken: String) {
    val request = DeviceRegisterRequest(
        pushToken = pushToken,
        deviceType = "phone",
        osVersion = "HarmonyOS 5.0",
        appVersion = BuildConfig.VERSION_NAME
    )
    
    api.registerDevice(request)
        .enqueue(object : Callback<DeviceRegisterResponse> {
            override fun onResponse(call: Call<DeviceRegisterResponse>, response: Response<DeviceRegisterResponse>) {
                if (response.isSuccessful) {
                    val deviceKey = response.body()?.deviceKey
                    // 保存 deviceKey 到本地
                    saveDeviceKey(deviceKey)
                }
            }
            
            override fun onFailure(call: Call<DeviceRegisterResponse>, t: Throwable) {
                // 处理错误
            }
        })
}
```

### 接收推送消息

```kotlin
class MyPushService : HmsMessageService() {
    override fun onNewToken(token: String) {
        // Token 刷新时更新
        updateDeviceToken(token)
    }
    
    override fun onMessageReceived(message: RemoteMessage) {
        // 接收推送消息
        when (message.messageType) {
            MessageType.NOTIFICATION -> {
                // 通知消息
                showNotification(message)
            }
            MessageType.DATA -> {
                // 后台消息
                handleBackgroundMessage(message.data)
            }
        }
    }
}
```

## 生产环境部署

### 1. 修改配置

编辑 `.env` 文件:
```bash
# 关闭测试消息模式
TEST_MESSAGE_MODE=false

# 设置生产环境URL
PUSH_API_URL=https://push-api.cloud.huawei.com

# 设置日志级别
LOG_LEVEL=info
```

### 2. 使用正式服务账号
确保使用生产环境的服务账号密钥文件

### 3. 配置HTTPS
建议使用反向代理（Nginx/Caddy）配置HTTPS

### 4. 监控和告警
- 监控推送成功率
- 设置失败告警阈值
- 记录推送日志便于排查

## 下一步

- [ ] 阅读 [API_USAGE.md](API_USAGE.md) 了解完整API文档
- [ ] 阅读 [PUSH_API_FIX.md](PUSH_API_FIX.md) 了解API修正详情
- [ ] 查看 [华为 Push Kit 官方文档](https://developer.huawei.com/consumer/cn/doc/harmonyos-references/push-api)
