# 快速开始指南

## 第一步：配置环境

### 1. 生成加密密钥

```bash
# 运行密钥生成脚本
./scripts/generate-keys.sh

# 或手动生成
openssl rand -base64 32
```

### 2. 配置数据库

```bash
# 创建数据库
createdb push_server

# 或使用PostgreSQL客户端
psql -U postgres
CREATE DATABASE push_server;
\q
```

### 3. 配置华为Push Kit

从 [AppGallery Connect](https://developer.huawei.com/consumer/cn/service/josp/agc/index.html) 下载服务账号文件：

1. 登录并选择项目
2. 项目设置 → 常规 → 我的应用
3. 下载 `agconnect-services.json` 文件
4. 将文件保存到 `config/agconnect-services.json`

### 4. 编辑配置文件

编辑 `.env`:

```bash
# 数据库
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=push_server

# 华为Push Kit
HUAWEI_PROJECT_ID=101653523863440882  # 从agconnect-services.json的client.project_id获取
HUAWEI_SERVICE_ACCOUNT_FILE=./config/agconnect-services.json

# 加密密钥（从generate-keys.sh获取）
PUSH_TOKEN_ENCRYPTION_KEY=your_generated_key_here
```

## 第二步：启动服务

```bash
# 方式1：直接运行
go run cmd/server/main.go

# 方式2：编译后运行
go build -o bin/push-server cmd/server/main.go
./bin/push-server
```

看到以下输出说明启动成功：

```
Server starting on port 8080
Push API: https://push-api.cloud.huawei.com/v3/1234567890
```

## 第三步：测试接口

### 1. 健康检查

```bash
curl http://localhost:8080/health
```

预期响应：

```json
{
    "status": "ok",
    "version": "1.0.0",
    "service": "Dangdangdang Push Server (Huawei Push Kit v3)"
}
```

### 2. 注册设备

```bash
curl -X POST http://localhost:8080/api/v1/device/register \
  -H "Content-Type: application/json" \
  -d '{
    "push_token": "APA91bHun4MxP5egoKMwt2KZFBaFUH...",
    "device_type": "phone",
    "os_version": "HarmonyOS 5.0",
    "app_version": "1.0.0"
  }'
```

预期响应：

```json
{
    "success": true,
    "device_key": "550e8400-e29b-41d4-a716-446655440000",
    "message": "Device registered successfully"
}
```

**重要**: 保存返回的 `device_key`！

### 3. 发送推送消息

```bash
# 使用device_key发送推送
curl "http://localhost:8080/api/v1/push/notification?device_key=550e8400-e29b-41d4-a716-446655440000&title=测试消息&body=这是一条测试推送"
```

预期响应：

```json
{
    "success": true,
    "message": "Notification sent successfully"
}
```

### 4. 发送卡片刷新

```bash
# URL编码的JSON数据
curl "http://localhost:8080/api/v1/push/form?device_key=550e8400-e29b-41d4-a716-446655440000&form_id=weather_card&form_data=%7B%22temperature%22%3A%2225%C2%B0C%22%2C%22weather%22%3A%22%E6%99%B4%E5%A4%A9%22%7D"
```

### 5. 批量推送

```bash
curl "http://localhost:8080/api/v1/push/batch?device_keys=key1,key2,key3&title=批量通知&body=这是批量推送消息"
```

### 6. 查询统计

```bash
curl "http://localhost:8080/api/v1/push/statistics?date=2026-01-13"
```

## 第四步：集成到客户端

### HarmonyOS客户端示例

```typescript
// 1. 注册设备
import { pushService } from '@kit.PushKit';

async function registerPushDevice() {
  try {
    // 获取Push Token
    const pushToken = await pushService.getToken();
    
    // 向服务器注册
    const response = await fetch('http://your-server.com/api/v1/device/register', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        push_token: pushToken,
        device_type: 'phone',
        os_version: 'HarmonyOS 5.0',
        app_version: '1.0.0'
      })
    });
    
    const result = await response.json();
    
    // 保存device_key
    preferences.put('device_key', result.device_key);
    
    console.log('Device registered:', result.device_key);
  } catch (error) {
    console.error('Registration failed:', error);
  }
}

// 2. 应用启动时调用
@Entry
@Component
struct Index {
  aboutToAppear() {
    registerPushDevice();
  }
}
```

## 常见问题

### 1. 编译失败

```bash
# 清理并重新下载依赖
go clean -modcache
go mod download
go mod tidy
```

### 2. 数据库连接失败

检查：
- PostgreSQL是否运行：`pg_isready`
- 数据库是否存在：`psql -l`
- 连接信息是否正确：`.env`文件

### 3. JWT认证失败

检查：
- `config/agconnect-services.json` 文件是否正确
- 项目ID是否匹配（从client.project_id字段获取）
- Push Kit服务是否已开启

### 4. 推送失败

常见原因：
- Push Token无效或过期
- Device Key不存在
- 华为Push服务配额已用尽

查看日志获取详细错误信息。

## 生产部署建议

1. **使用HTTPS**: 部署反向代理（Nginx/Caddy）
2. **数据库备份**: 定期备份PostgreSQL
3. **日志监控**: 使用ELK或Prometheus
4. **限流保护**: 配置API限流
5. **密钥安全**: 
   - 不要提交 `.env` 到Git
   - 使用密钥管理服务（如HashiCorp Vault）
   - 定期轮换服务账号密钥

## 下一步

- 阅读完整 [API文档](README_NEW.md)
- 查看 [安全最佳实践](README_NEW.md#安全说明)
- 加入开发者社区

## 获取帮助

- 🐛 [提交Bug](https://github.com/yourusername/dangdangdang-push-server/issues)
- 💬 [讨论区](https://github.com/yourusername/dangdangdang-push-server/discussions)
- 📧 Email: your@email.com
