# 🚀 快速命令参考

## 启动和管理

```bash
# 重启服务（推荐）
./restart.sh

# 手动启动
docker compose -f docker-compose.single.yml up -d

# 停止服务
docker compose -f docker-compose.single.yml down

# 重新构建并启动
docker compose -f docker-compose.single.yml up --build -d
```

## 📋 查看日志

### 方法1: 使用交互式脚本（推荐）
```bash
./view-logs.sh
```

### 方法2: 直接命令

```bash
# 实时日志（跟踪模式）
docker compose -f docker-compose.single.yml logs -f

# 最近100行
docker compose -f docker-compose.single.yml logs --tail=100

# 只看错误
docker compose -f docker-compose.single.yml logs | grep ERROR

# 只看推送相关
docker compose -f docker-compose.single.yml logs | grep -i "push\|notification"

# 查看OAuth认证
docker compose -f docker-compose.single.yml logs | grep -i "oauth\|token"
```

## 🔍 日志类型说明

- `[INFO]` - 服务启动、配置信息等
- `[ERROR]` - 错误日志，需要关注
- `[DEBUG]` - 调试信息，详细的请求/响应数据
- `[ACCESS]` - HTTP请求日志
- `→` - 收到请求
- `←` - 返回响应
- `✓` - 成功
- `✗` - 失败

## 🧪 测试推送

```bash
# 健康检查
curl http://localhost:8081/health

# 发送推送通知
curl "http://localhost:8081/api/v1/push/notification?device_key=YOUR_DEVICE_KEY&title=测试&body=测试消息"
```

## 📖 详细文档

查看完整日志指南: [LOGGING_GUIDE.md](LOGGING_GUIDE.md)
