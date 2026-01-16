# Dangdangdang Push Server

基于华为 Push Kit v3 REST API 的推送服务器，支持 HarmonyOS Next/5.x 及更高版本。

## 功能特性

✅ **华为 Push Kit v3 支持** - 使用最新的JWT认证方式  
✅ **安全设计** - Device Key替代Push Token，加密存储敏感信息  
✅ **隐私友好** - 简化数据记录，不存储推送内容  
✅ **GET接口** - 推送接口支持GET方式，方便直接调用  
✅ **多种消息类型**:
  - 通知消息（Alert）
  - 卡片刷新消息（Form Update）
  - 后台消息（Background）
  - 批量推送（最多1000个设备）

## 快速开始

### 1. 前置要求

- Go 1.21+
- PostgreSQL 12+
- 华为开发者账号
- AppGallery Connect 项目（已开启 Push Kit）

### 2. 获取华为服务账号密钥

1. 登录 [AppGallery Connect](https://developer.huawei.com/consumer/cn/service/josp/agc/index.html)
2. 选择项目 → **项目设置** → **API管理**
3. 开启 **Push Kit** 服务
4. 点击 **服务账号** → **创建服务账号**
5. 下载 JSON 格式的密钥文件

### 3. 安装配置

```bash
# 克隆项目
git clone https://github.com/yourusername/dangdangdang-push-server.git
cd dangdangdang-push-server

# 安装依赖
go mod download

# 复制.env配置文件
cp .env.example .env

# 从AppGallery Connect下载agconnect-services.json到config/目录
```
```

### 4. 配置环境变量

编辑 `.env` 文件：

```bash
# 数据库配置
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=push_server

# 华为Push Kit配置
HUAWEI_PROJECT_ID=101653523863440882  # 从agconnect-services.json的client.project_id获取
HUAWEI_SERVICE_ACCOUNT_FILE=./config/agconnect-services.json

# 生成32字节加密密钥
# openssl rand -base64 32
PUSH_TOKEN_ENCRYPTION_KEY=abcdefghijklmnopqrstuvwxyz123456
```

### 5. 初始化数据库

```bash
# 创建数据库
createdb push_server

# 运行服务（自动创建表）
go run cmd/server/main.go
```

### 6. 启动服务

```bash
# 开发模式
go run cmd/server/main.go

# 编译运行
go build -o bin/push-server cmd/server/main.go
./bin/push-server
```

## API接口文档

### 设备管理

#### 1. 注册设备

**接口**: `POST /api/v1/device/register`

**请求体**:
```json
{
    "push_token": "APA91bHun4MxP5egoKMwt2KZFBaFUH...",
    "device_type": "phone",
    "os_version": "HarmonyOS 5.0",
    "app_version": "1.0.0"
}
```

**响应**:
```json
{
    "success": true,
    "device_key": "550e8400-e29b-41d4-a716-446655440000",
    "message": "Device registered successfully"
}
```

**重要**: 客户端需要保存返回的 `device_key`，后续所有推送接口都使用此Key。

#### 2. 更新Push Token

**接口**: `PUT /api/v1/device/update-token`

```json
{
    "device_key": "550e8400-e29b-41d4-a716-446655440000",
    "new_push_token": "NEW_TOKEN_FROM_HUAWEI..."
}
```

#### 3. 停用设备

**接口**: `GET /api/v1/device/deactivate?device_key=xxx`

### 推送消息（GET方式）

#### 1. 发送通知消息

**接口**: `GET /api/v1/push/notification`

**参数**:
- `device_key` (必需): 设备Key
- `title` (必需): 通知标题
- `body` (必需): 通知内容
- `data` (可选): 额外数据，JSON字符串

**示例**:
```
GET /api/v1/push/notification?device_key=550e8400-e29b-41d4-a716-446655440000&title=新消息&body=您有一条新通知&data={"order_id":"12345"}
```

**响应**:
```json
{
    "success": true,
    "message": "Notification sent successfully"
}
```

#### 2. 发送卡片刷新消息

**接口**: `GET /api/v1/push/form`

**参数**:
- `device_key` (必需): 设备Key
- `form_id` (必需): 卡片ID
- `form_data` (必需): 卡片数据，JSON字符串

**示例**:
```
GET /api/v1/push/form?device_key=xxx&form_id=weather_card&form_data={"temperature":"25°C","weather":"晴天"}
```

#### 3. 发送后台消息

**接口**: `GET /api/v1/push/background`

**参数**:
- `device_key` (必需): 设备Key
- `data` (必需): 消息数据，JSON字符串

**示例**:
```
GET /api/v1/push/background?device_key=xxx&data={"action":"sync","timestamp":"2026-01-13T10:00:00Z"}
```

#### 4. 批量推送

**接口**: `GET /api/v1/push/batch`

**参数**:
- `device_keys` (必需): 设备Key列表，逗号分隔，最多1000个
- `title` (必需): 通知标题
- `body` (必需): 通知内容
- `data` (可选): 额外数据，JSON字符串

**示例**:
```
GET /api/v1/push/batch?device_keys=key1,key2,key3&title=系统通知&body=重要更新
```

**响应**:
```json
{
    "success": true,
    "message": "Batch notification sent successfully",
    "total_sent": 3,
    "total_failed": 0
}
```

### 统计查询

#### 查询推送统计

**接口**: `GET /api/v1/push/statistics?date=2026-01-13`

**响应**:
```json
{
    "success": true,
    "date": "2026-01-13",
    "stats": [
        {
            "push_type": "notification",
            "total_count": 100,
            "success_count": 95,
            "failed_count": 5
        },
        {
            "push_type": "form",
            "total_count": 50,
            "success_count": 48,
            "failed_count": 2
        }
    ]
}
```

## 安全说明

### Device Key 设计

- **客户端**：只存储和传输 `device_key`（UUID格式）
- **服务端**：内部将 `device_key` 映射到加密的 `push_token`
- **优势**：即使接口被抓包，也无法获取真实的Push Token

### Push Token 加密

- 使用 AES-256-GCM 加密算法
- 密钥通过环境变量配置，不写入代码
- 数据库中存储的是加密后的Token

### 数据隐私

- **不记录推送内容**：仅统计推送数量，不存储消息标题和内容
- **开源友好**：适合开源项目使用

## 客户端集成示例（HarmonyOS）

```typescript
// 1. 注册设备，获取device_key
async function registerDevice(pushToken: string): Promise<string> {
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
  
  const data = await response.json();
  // 保存device_key到本地持久化存储
  preferences.put('device_key', data.device_key);
  return data.device_key;
}

// 2. 发送推送消息（服务端或其他客户端）
async function sendNotification(deviceKey: string, title: string, body: string) {
  const url = `http://your-server.com/api/v1/push/notification?device_key=${deviceKey}&title=${encodeURIComponent(title)}&body=${encodeURIComponent(body)}`;
  const response = await fetch(url);
  return await response.json();
}
```

## 错误码说明

| 错误码 | 说明 | 处理方案 |
|--------|------|----------|
| 80000000 | 成功 | - |
| 80100000 | 部分Token失败 | 检查返回的illegal_tokens |
| 80200001 | JWT认证失败 | 检查服务账号配置 |
| 80200005 | JWT过期 | 自动刷新（已内置） |
| 80300007 | Token无效 | 客户端重新注册设备 |
| 80300008 | 消息体过大 | 减小消息内容（<4KB） |

## 开发调试

### 测试健康检查

```bash
curl http://localhost:8080/health
```

### 测试设备注册

```bash
curl -X POST http://localhost:8080/api/v1/device/register \
  -H "Content-Type: application/json" \
  -d '{
    "push_token": "test_token_123",
    "device_type": "phone",
    "os_version": "HarmonyOS 5.0"
  }'
```

### 测试推送（使用device_key）

```bash
# 注意：需要先注册设备获取device_key
curl "http://localhost:8080/api/v1/push/notification?device_key=YOUR_DEVICE_KEY&title=测试&body=这是一条测试消息"
```

## 生产部署

### Docker部署

```bash
# 构建镜像
docker build -t dangdangdang-push-server .

# 运行容器
docker run -d \
  --name push-server \
  -p 8080:8080 \
  -e DB_HOST=your_db_host \
  -e HUAWEI_PROJECT_ID=your_project_id \
  -v ./config:/app/config \
  dangdangdang-push-server
```

### Systemd服务

创建 `/etc/systemd/system/push-server.service`:

```ini
[Unit]
Description=Dangdangdang Push Server
After=network.target postgresql.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/push-server
ExecStart=/opt/push-server/bin/push-server
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！

## 联系方式

- Issues: [GitHub Issues](https://github.com/yourusername/dangdangdang-push-server/issues)
- Email: your@email.com
