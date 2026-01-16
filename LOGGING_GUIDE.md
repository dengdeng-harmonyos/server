# 日志查看指南

## 🔍 查看日志的多种方式

### 方式1: 使用便捷脚本 (推荐)

```bash
./view-logs.sh
```

这个交互式脚本提供了6种查看方式：
1. **实时日志** - 持续跟踪最新日志
2. **最近100行** - 查看最近的日志
3. **只看错误** - 过滤出错误信息
4. **只看访问** - 查看HTTP请求记录
5. **搜索关键词** - 查找特定内容
6. **导出日志** - 保存到文件

### 方式2: Docker Compose 命令

#### 实时查看日志（跟踪模式）
```bash
docker compose -f docker-compose.single.yml logs -f
```

#### 查看最近N行日志
```bash
docker compose -f docker-compose.single.yml logs --tail=100
```

#### 查看特定时间的日志
```bash
docker compose -f docker-compose.single.yml logs --since 30m   # 最近30分钟
docker compose -f docker-compose.single.yml logs --since 2h    # 最近2小时
```

#### 只看某个服务的日志
```bash
docker compose -f docker-compose.single.yml logs app
```

### 方式3: 直接用 Docker 命令

```bash
# 找到容器ID
docker ps | grep dangdangdang

# 查看日志
docker logs -f <container_id>

# 查看最近100行
docker logs --tail 100 <container_id>
```

### 方式4: 使用 grep 过滤

#### 只看错误日志
```bash
docker compose -f docker-compose.single.yml logs | grep -i "ERROR\|error"
```

#### 只看特定关键词
```bash
docker compose -f docker-compose.single.yml logs | grep "access token"
docker compose -f docker-compose.single.yml logs | grep "push"
```

#### 查看HTTP请求
```bash
docker compose -f docker-compose.single.yml logs | grep "ACCESS"
```

### 方式5: 保存日志到文件

```bash
# 导出所有日志
docker compose -f docker-compose.single.yml logs > full_logs.txt

# 导出最近1000行
docker compose -f docker-compose.single.yml logs --tail=1000 > recent_logs.txt

# 导出并实时追踪
docker compose -f docker-compose.single.yml logs -f | tee live_logs.txt
```

## 📊 日志级别说明

新的日志系统包含以下级别：

- `[INFO]` - 一般信息，如服务启动、配置加载
- `[ERROR]` - 错误信息，需要关注
- `[DEBUG]` - 调试信息，包含详细的请求/响应数据
- `[ACCESS]` - HTTP访问日志，记录所有API请求

## 🎯 常用查看场景

### 场景1: 服务启动问题
```bash
docker compose -f docker-compose.single.yml logs --tail=50 | grep -i "starting\|error\|failed"
```

### 场景2: 追踪推送请求
```bash
docker compose -f docker-compose.single.yml logs -f | grep -i "push\|notification"
```

### 场景3: 查看OAuth认证过程
```bash
docker compose -f docker-compose.single.yml logs | grep -i "oauth\|token\|access"
```

### 场景4: 监控HTTP请求
```bash
docker compose -f docker-compose.single.yml logs -f | grep "ACCESS"
```

### 场景5: 排查错误
```bash
docker compose -f docker-compose.single.yml logs | grep -B 5 -A 5 "ERROR"
# -B 5: 显示错误前5行
# -A 5: 显示错误后5行
```

## 🔧 日志配置

当前日志配置位于代码中，包括：

1. **启动日志** - 显示所有配置信息
2. **请求日志** - 记录每个HTTP请求的详情
3. **推送日志** - 记录推送过程的每一步
4. **错误日志** - 详细的错误堆栈和上下文

## 💡 提示

- 使用 `Ctrl+C` 退出实时日志模式
- 日志中的 `✓` 表示成功操作
- 日志中的 `✗` 表示失败操作
- 使用 `--timestamps` 参数可以显示准确时间戳

```bash
docker compose -f docker-compose.single.yml logs --timestamps
```

## 🐛 调试模式

如需更详细的调试信息，可以在 docker-compose.single.yml 中设置：

```yaml
environment:
  - GIN_MODE=debug
```

然后重启服务：
```bash
docker compose -f docker-compose.single.yml restart
```
