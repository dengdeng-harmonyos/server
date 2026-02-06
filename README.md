# 噔噔推送服务 (Dengdeng Push Server)

[![GitHub release](https://img.shields.io/github/v/release/dengdeng-harmonyos/server)](https://github.com/dengdeng-harmonyos/server/releases)
[![GitHub stars](https://img.shields.io/github/stars/dengdeng-harmonyos/server?style=social)](https://github.com/dengdeng-harmonyos/server)
[![License](https://img.shields.io/github/license/dengdeng-harmonyos/server)](LICENSE)
[![Docker Pulls](https://img.shields.io/docker/pulls/ricwang/dengdeng-server)](https://hub.docker.com/r/ricwang/dengdeng-server)
[![Go Version](https://img.shields.io/github/go-mod/go-version/dengdeng-harmonyos/server)](go.mod)

[English](README_EN.md) | 简体中文

## 📖 项目简介

噔噔推送服务是一个专为 **HarmonyOS Next** 设计的**安全、隐私友好**的推送服务解决方案。本项目完全开源，致力于为开发者提供一个可信赖、易部署的推送服务基础设施。

> 🎯 **v1.0 正式发布**：生产就绪，支持完整的推送功能和自动化部署

### ✨ 主要亮点

- **🚀 一键部署**：单个 Docker 容器即可运行，内置 PostgreSQL 数据库
- **🔐 安全优先**：配置编译时嵌入，支持 AES-256-GCM 加密
- **📦 零依赖**：无需外部配置文件，开箱即用
- **🤖 CI/CD 自动化**：GitHub Actions 自动构建和部署
- **🌍 生产就绪**：支持测试和生产环境分离部署

### 🔒 安全与隐私承诺

- **🚫 零消息存储**：不存储任何推送消息内容，仅保存匿名统计数据
- **🔐 端到端加密**：Push Token 使用 AES-256-GCM 加密存储
- **🎭 匿名化设计**：使用随机生成的 Device Id，与真实设备无关联
- **📊 统计数据脱敏**：仅保存推送成功/失败次数，不记录具体内容
- **🔑 配置编译时嵌入**：敏感配置在构建时嵌入二进制文件，无需配置文件
- **🛡️ 开源透明**：所有代码公开，接受社区审查

## ✨ 核心特性

### 🚀 部署与运维

- **📦 单容器部署**：包含 PostgreSQL + 推送服务，开箱即用
- **🔧 配置嵌入**：华为推送配置编译时嵌入，无需外部文件
- **🤖 自动化 CI/CD**：GitHub Actions 自动构建、测试和部署
- **🏥 健康检查**：内置健康检查接口，支持监控
- **🐳 Docker 支持**：官方镜像托管在 Docker Hub
- **🔄 自动重启**：容器崩溃自动恢复

### 🔐 安全性

- **🔒 AES-256-GCM 加密**：保护 Push Token 存储安全
- **🎲 加密安全随机数**：使用 crypto/rand 生成 Device Id
- **🔑 RSA 公钥支持**：可选的端到端消息加密
- **⏱️ 自动过期机制**：Device Id 时效性管理
- **🚦 速率限制**：防止推送滥用（每设备每日限额）
- **🛡️ 编译时密钥注入**：通过 ldflags 嵌入敏感配置

### 🎯 隐私保护

- **📝 零消息存储**：不保存任何推送消息内容
- **🎭 完全匿名**：设备标识无法追溯到真实设备
- **📊 聚合统计**：仅记录统计数据，无法追溯具体设备
- **🗑️ 自动清理**：定期清理过期设备记录
- **🔍 最小化原则**：数据库字段遵循最小必要原则

### 🎯 隐私保护

- **📝 零消息存储**：不保存任何推送消息内容
- **🎭 完全匿名**：设备标识无法追溯到真实设备
- **📊 聚合统计**：仅记录统计数据，无法追溯具体设备
- **🗑️ 自动清理**：定期清理过期设备记录
- **🔍 最小化原则**：数据库字段遵循最小必要原则

### 📡 功能特性

- **📬 通知推送**：支持通知栏消息（带标题、内容、自定义数据）
- **🃏 卡片刷新**：支持 HarmonyOS 卡片更新
- **🔄 后台推送**：支持后台数据推送
- **📦 批量推送**：一次性向多个设备发送消息
- **📊 推送统计**：查看推送成功率和历史数据
- **🏥 健康监控**：内置健康检查和服务状态接口
- **🌐 RESTful API**：简洁的 HTTP GET 接口，易于集成

## 🚀 快速开始

### 前提条件

在开始之前，你需要：

1. **华为开发者账号**：[华为开发者联盟](https://developer.huawei.com/)
2. **HarmonyOS 应用**：已创建的 HarmonyOS Next 应用
3. **推送服务配置**：
   - `agconnect-services.json` - 从 AppGallery Connect 下载
   - `private.json` - 华为推送服务账号私钥

### 方式一：使用 Docker Hub 镜像（推荐）

这是最简单快速的部署方式：

#### 1. 生成加密密钥

```bash
# 生成 32 字节随机密钥（Base64 编码）
openssl rand -base64 32
```

将生成的密钥保存到 `.env` 文件：

```bash
echo "PUSH_TOKEN_ENCRYPTION_KEY=你生成的密钥" > .env
```

#### 2. 启动服务

```bash
# 拉取最新镜像
docker pull ricwang/dengdeng-server:latest

# 启动服务
docker run -d \
  --name push-server \
  -p 8080:8080 \
  -e PUSH_TOKEN_ENCRYPTION_KEY=你的加密密钥 \
  -e SERVER_NAME=噔噔推送服务 \
  -v push-data:/var/lib/postgresql/data \
  --restart unless-stopped \
  ricwang/dengdeng-server:latest
```

> ⚠️ **注意**：Docker Hub 镜像使用编译时嵌入的华为推送配置，仅适用于公共演示。生产环境请使用方式二自行构建。

#### 2. 启动服务

```bash
# 拉取最新镜像
docker pull ricwang/dengdeng-server:latest

# 启动服务
docker run -d \
  --name push-server \
  -p 8080:8080 \
  -e PUSH_TOKEN_ENCRYPTION_KEY=你的加密密钥 \
  -e SERVER_NAME=噔噔推送服务 \
  -v push-data:/var/lib/postgresql/data \
  --restart unless-stopped \
  ricwang/dengdeng-server:latest
```

> ⚠️ **注意**：Docker Hub 镜像使用编译时嵌入的华为推送配置，仅适用于公共演示。生产环境请使用方式二自行构建。

#### 3. 验证服务

```bash
# 检查健康状态
curl http://localhost:8080/health

# 查看日志
docker logs -f push-server
```

### 方式二：使用自己的配置构建（生产推荐）

如果你要部署到生产环境，建议使用自己的华为推送配置：

#### 1. 准备配置文件

将从华为开发者后台下载的配置文件保存到 GitHub Secrets：

- `AGCONNECT_JSON` - `agconnect-services.json` 的完整内容
- `PRIVATE_JSON` - `private.json` 的完整内容
- `PUSH_TOKEN_ENCRYPTION_KEY` - 使用 `openssl rand -base64 32` 生成

#### 2. Fork 仓库并配置 Secrets

1. Fork 本仓库到你的 GitHub 账号
2. 在仓库设置中添加上述 Secrets
3. 推送代码到 `main` 分支（测试环境）或 `release` 分支（生产环境）

#### 3. 自动构建和部署

GitHub Actions 会自动：
- ✅ 编译时嵌入你的华为推送配置
- ✅ 构建优化的静态链接二进制文件
- ✅ 构建 Docker 镜像
- ✅ 部署到你配置的服务器

### 方式三：本地开发构建

```bash
# 克隆仓库
git clone https://github.com/dengdeng-harmonyos/server.git
cd server

# 准备配置文件（放在项目根目录）
# - agconnect-services.json
# - private.json

# 生成加密密钥
echo "PUSH_TOKEN_ENCRYPTION_KEY=$(openssl rand -base64 32)" > .env

# 方式 A：使用 Docker Compose
docker-compose up -d --build

# 方式 B：本地编译运行
go mod download
go build -o bin/server cmd/server/main.go

# 启动数据库
docker run -d --name postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=push_server \
  -p 5432:5432 \
  postgres:15-alpine

# 运行服务器
export AGCONNECT_SERVICES_FILE=agconnect-services.json
export PRIVATE_KEY_FILE=private.json
./bin/server
```

## 📡 API 接口

### 快速概览

所有接口使用简单的 HTTP GET 请求，无需复杂的认证流程。

| 功能 | 接口路径 | 说明 |
|------|---------|------|
| 健康检查 | `GET /health` | 检查服务状态 |
| 设备注册 | `GET /api/v1/device/register` | 注册设备获取 Device Id |
| 通知推送 | `GET /api/v1/push/notification` | 发送通知栏消息 |

### 示例：发送通知

```bash
curl "http://your-server:8080/api/v1/push/notification?device_id=YOUR_DEVICE_KEY&title=测试消息&content=这是一条测试推送"
```

### 示例：批量推送

```bash
curl "http://your-server:8080/api/v1/push/batch?device_ids=key1,key2,key3&title=批量通知&body=发送给多个设备"
```

### 完整文档

详细的 API 文档和参数说明，请参考：

- 📚 **API 文档**：查看仓库中的 API 使用示例
- 🔍 **源码参考**：[internal/handler](internal/handler) 目录
- 💡 **集成示例**：查看 HarmonyOS 客户端项目

## 🔧 配置说明

### 环境变量配置

| 环境变量 | 说明 | 必需 | 默认值 |
|---------|------|:----:|--------|
| `PUSH_TOKEN_ENCRYPTION_KEY` | Push Token 加密密钥（32字节，Base64） | ✅ | - |
| `SERVER_NAME` | 服务器标识名称 | ❌ | `噔噔推送服务` |
| `PORT` | HTTP 服务端口 | ❌ | `8080` |
| `GIN_MODE` | 运行模式（debug/release） | ❌ | `release` |
| `MAX_DAILY_PUSH_PER_DEVICE` | 每设备每日推送限额 | ❌ | `1000` |

### 数据持久化

Docker 容器使用命名卷存储 PostgreSQL 数据：

```bash
# 查看数据卷
docker volume ls | grep push-data

# 备份数据
docker run --rm -v push-data:/data -v $(pwd):/backup alpine \
  tar czf /backup/push-data-backup.tar.gz /data

# 恢复数据
docker run --rm -v push-data:/data -v $(pwd):/backup alpine \
  tar xzf /backup/push-data-backup.tar.gz -C /
```

## 📊 数据存储说明

### 存储的数据

1. **设备信息**（匿名化）
   - Device Id（随机生成）
   - Push Token（AES-256-GCM 加密）
   - 设备元数据（类型、版本等）
   - RSA 公钥（可选）

2. **统计数据**（聚合）
   - 每日推送次数
   - 成功/失败次数
   - 推送类型分布

### 不存储的数据

- ❌ 推送消息内容
- ❌ 用户身份信息
- ❌ 设备硬件标识
- ❌ 地理位置信息
- ❌ IP 地址
- ❌ 任何可追溯到用户的信息

## 🏗️ 架构设计

### 系统架构

```
┌─────────────────────────────────────────────────────────┐
│                    Docker 容器                           │
│  ┌─────────────────┐         ┌──────────────────┐      │
│  │  PostgreSQL 15  │ ←────→ │  推送服务 (Go)    │      │
│  │  - 设备信息      │         │  - Gin Web框架    │      │
│  │  - 加密Token    │         │  - AES-256加密    │      │
│  │  - 推送统计      │         │  - 华为推送API    │      │
│  └─────────────────┘         └──────────────────┘      │
│         ↓                             ↑                 │
│    数据持久化卷                    端口8080              │
└───────────────────────────────────────┬─────────────────┘
                                        │
                                   HTTP API
                                        │
                    ┌───────────────────┼───────────────────┐
                    ↓                   ↓                   ↓
            HarmonyOS 应用 1     HarmonyOS 应用 2    其他客户端
```

### 数据流程

#### 1. 设备注册流程

```
客户端                推送服务                数据库
  │                      │                      │
  │─ 注册请求 ─────────→ │                      │
  │                      │                      │
  │                      │─ 生成 Device Id ──→ │
  │                      │  (crypto/rand)       │
  │                      │                      │
  │                      │─ 加密 Push Token ──→ │
  │                      │  (AES-256-GCM)       │
  │                      │                      │
  │← 返回 Device Id ──  │                      │
```

#### 2. 推送消息流程

```
应用后端              推送服务                华为推送
  │                      │                      │
  │─ 推送请求 ─────────→ │                      │
  │  (Device Id)        │                      │
  │                      │─ 解密 Token ───────→ │
  │                      │                      │
  │                      │                      │─ 发送推送 ──→ 用户设备
  │                      │← 返回结果 ──────────  │
  │                      │                      │
  │← 推送成功 ──────────  │                      │
  │                      │                      │
  │                      │─ 记录统计 ──────────→ 数据库
  │                      │  (不含消息内容)
```

### 安全机制

1. **编译时配置嵌入**
   ```
   源码 + Secrets → GitHub Actions
          ↓
   ldflags 编译时注入 (Base64)
          ↓
   静态链接二进制文件 (无外部依赖)
          ↓
   Docker 镜像 (配置已嵌入)
   ```

2. **Push Token 加密存储**
   ```
   明文 Token → AES-256-GCM 加密 → 数据库
   (随机 Nonce)     (32字节密钥)
   ```

3. **Device Id 生成**
   ```
   crypto/rand → Base64 URL Safe → 存储
   (32字节随机)    (无特殊字符)
   ```

## 🔐 安全最佳实践

### 1. 密钥管理

**生成强密钥**
```bash
# 推荐：使用 OpenSSL 生成 32 字节随机密钥
openssl rand -base64 32

# 或使用 /dev/urandom (Linux/macOS)
head -c 32 /dev/urandom | base64
```

**密钥轮换**
```bash
# 1. 生成新密钥
NEW_KEY=$(openssl rand -base64 32)

# 2. 更新 GitHub Secrets 或环境变量

# 3. 重新构建和部署
# GitHub Actions 会自动使用新密钥编译

# 4. 旧设备需要重新注册
```

**存储安全**
- ✅ 使用环境变量或密钥管理服务
- ✅ 使用 GitHub Secrets 存储敏感配置
- ✅ 编译时嵌入，避免配置文件暴露
- ❌ 不要在代码中硬编码
- ❌ 不要提交到 Git 仓库
- ❌ 不要通过日志输出

### 2. 网络安全

**使用 HTTPS**
```nginx
# Nginx 反向代理配置示例
server {
    listen 443 ssl http2;
    server_name push.yourdomain.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

**防火墙规则**
```bash
# 仅允许特定 IP 访问（可选）
sudo ufw allow from YOUR_IP to any port 8080

# 或仅允许内网访问
sudo ufw allow from 10.0.0.0/8 to any port 8080
```

**速率限制**
```bash
# 应用级别已内置速率限制：
# - 每设备每日最大 1000 条推送
# - 可通过 MAX_DAILY_PUSH_PER_DEVICE 环境变量调整
```

### 3. 数据库安全

**定期备份**
```bash
# 创建备份
docker exec push-server pg_dump -U postgres push_server > backup-$(date +%Y%m%d).sql

# 自动备份脚本（添加到 crontab）
0 2 * * * docker exec push-server pg_dump -U postgres push_server | gzip > /backup/push-$(date +\%Y\%m\%d).sql.gz
```

**清理过期数据**
```sql
-- 清理 30 天前的过期设备
DELETE FROM devices WHERE expired_at < NOW() - INTERVAL '30 days';

-- 清理 90 天前的推送统计
DELETE FROM push_statistics WHERE date < NOW() - INTERVAL '90 days';
```

### 4. 监控与审计

**健康监控**
```bash
# 基础健康检查
curl http://localhost:8080/health

# 配合监控系统（如 Prometheus）
# 可以定期检查健康状态并告警
```

**日志审计**
```bash
# 查看推送日志
docker logs push-server | grep "Push"

# 查看错误日志
docker logs push-server | grep "ERROR"

# 实时监控
docker logs -f push-server
```

**异常检测**
```bash
# 检查异常高频推送
# 查看推送统计 API
curl "http://localhost:8080/api/v1/push/statistics?date=$(date +%Y-%m-%d)"
```

### 5. 部署安全检查清单

部署前确认：

- [ ] ✅ 已生成强随机加密密钥
- [ ] ✅ 已配置 HTTPS/TLS
- [ ] ✅ 已设置防火墙规则
- [ ] ✅ 已配置数据备份策略
- [ ] ✅ 已启用健康检查监控
- [ ] ✅ 已审查日志输出（无敏感信息）
- [ ] ✅ 已限制服务器访问权限
- [ ] ✅ 已更新所有依赖到最新版本
- [ ] ✅ 已配置自动重启策略
- [ ] ✅ 已测试推送功能正常

## 📦 Docker 镜像

### 官方镜像

🐳 **Docker Hub**: [ricwang/dengdeng-server](https://hub.docker.com/r/ricwang/dengdeng-server)

### 可用标签

| 标签 | 说明 | 更新频率 |
|------|------|---------|
| `latest` | 最新稳定版本（main 分支） | 每次提交到 main |
| `v1.0.0`, `v1.0.x` | 特定版本号 | 发布时创建 |
| `release` | 生产发布版本 | 提交到 release 分支 |

### 镜像说明

- **基础镜像**: `postgres:15-alpine`
- **包含组件**: PostgreSQL 15 + Go 推送服务
- **镜像大小**: ~300MB
- **支持架构**: `linux/amd64`
- **配置方式**: 编译时嵌入（Docker Hub 镜像使用演示配置）

### 镜像构建

所有镜像通过 GitHub Actions 自动构建，确保：

- ✅ **可重现构建**：相同代码生成相同镜像
- ✅ **安全扫描**：构建过程无敏感信息泄露
- ✅ **静态链接**：无外部依赖，直接运行
- ✅ **最小化体积**：使用 Alpine 基础镜像

### 自建镜像

```bash
# 克隆仓库
git clone https://github.com/dengdeng-harmonyos/server.git
cd server

# 准备配置文件
# - agconnect-services.json
# - private.json

# 构建镜像
docker build -t my-dengdeng-server .

# 运行
docker run -d \
  --name push-server \
  -p 8080:8080 \
  -e PUSH_TOKEN_ENCRYPTION_KEY=$(openssl rand -base64 32) \
  -v push-data:/var/lib/postgresql/data \
  my-dengdeng-server
```

## 🛠️ 开发指南

### 环境要求

- **Go**: 1.21 或更高版本
- **PostgreSQL**: 15 或更高版本
- **Docker**: 20.10 或更高版本（可选）
- **Git**: 2.x

### 本地开发环境搭建

#### 1. 克隆代码

```bash
git clone https://github.com/dengdeng-harmonyos/server.git
cd server
```

#### 2. 安装依赖

```bash
go mod download
```

#### 3. 准备配置文件

将从华为开发者后台下载的配置文件放在项目根目录：
- `agconnect-services.json`
- `private.json`

#### 4. 启动数据库

```bash
# 使用 Docker 启动 PostgreSQL
docker run -d \
  --name postgres-dev \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=push_server \
  -p 5432:5432 \
  postgres:15-alpine

# 等待数据库启动
sleep 5

# 执行数据库迁移
cd database
./migrate.sh
cd ..
```

#### 5. 运行开发服务器

```bash
# 设置环境变量
export PUSH_TOKEN_ENCRYPTION_KEY=$(openssl rand -base64 32)
export GIN_MODE=debug
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=push_server

# 运行服务器
go run cmd/server/main.go
```

服务器将在 `http://localhost:8080` 启动。

### 编译构建

#### 本地编译

```bash
# 编译二进制文件
go build -o bin/server cmd/server/main.go

# 运行
./bin/server
```

#### 带配置嵌入的编译

```bash
# Base64 编码配置文件
AGCONNECT_BASE64=$(cat agconnect-services.json | base64)
PRIVATE_BASE64=$(cat private.json | base64)
ENCRYPTION_KEY_BASE64=$(echo "your-encryption-key" | base64)

# 编译时注入配置
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -ldflags "\
    -X 'github.com/dengdeng-harmonyos/server/internal/config.embeddedAgConnectJSON=$AGCONNECT_BASE64' \
    -X 'github.com/dengdeng-harmonyos/server/internal/config.embeddedPrivateJSON=$PRIVATE_BASE64' \
    -X 'github.com/dengdeng-harmonyos/server/internal/config.embeddedEncryptionKey=$ENCRYPTION_KEY_BASE64' \
    -s -w" \
  -o bin/dengdeng-server \
  cmd/server/main.go
```

#### 构建 Docker 镜像

```bash
# 本地构建
docker build -t dengdeng-server:dev .

# 使用 CI Dockerfile（需要预编译的二进制文件）
docker build -f Dockerfile.ci -t dengdeng-server:ci .
```

### 数据库管理

#### 创建迁移文件

```bash
cd database
./create-migration.sh add_new_feature
```

这会创建两个文件：
- `migrations/YYYYMMDDHHMMSS_add_new_feature.up.sql` - 正向迁移
- `migrations/YYYYMMDDHHMMSS_add_new_feature.down.sql` - 回滚迁移

#### 执行迁移

```bash
cd database
./migrate.sh
```

#### 回滚迁移

```bash
cd database
migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/push_server?sslmode=disable" down 1
```

### 项目结构

```
server/
├── cmd/
│   └── server/
│       └── main.go              # 应用入口
├── internal/
│   ├── config/
│   │   ├── config.go            # 配置加载
│   │   └── embedded_secrets.go  # 嵌入式配置
│   ├── database/
│   │   └── database.go          # 数据库操作
│   ├── handler/
│   │   ├── device.go            # 设备管理
│   │   ├── message.go           # 消息处理
│   │   ├── push.go              # 推送逻辑
│   │   └── response.go          # 响应封装
│   ├── logger/
│   │   └── logger.go            # 日志系统
│   ├── middleware/
│   │   └── middleware.go        # HTTP 中间件
│   ├── models/
│   │   └── models.go            # 数据模型
│   └── service/
│       ├── crypto.go            # 加密服务
│       ├── encryption.go        # Token 加密
│       └── huawei_push.go       # 华为推送 API
├── database/
│   ├── migrations/              # 数据库迁移文件
│   ├── migrate.sh              # 迁移脚本
│   └── 001_initial_schema.sql  # 初始数据库结构
├── .github/
│   └── workflows/
│       └── build.yml           # CI/CD 配置
├── Dockerfile                  # 标准 Dockerfile
├── Dockerfile.ci               # CI 专用 Dockerfile
├── docker-compose.yml          # Docker Compose 配置
└── README.md
```

### 测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./internal/service/...

# 运行测试并显示覆盖率
go test -cover ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 代码规范

```bash
# 格式化代码
go fmt ./...

# 静态检查
go vet ./...

# 使用 golangci-lint（推荐）
golangci-lint run
```

## 🤝 贡献指南

我们非常欢迎各种形式的贡献！无论是报告 bug、提出新功能建议，还是提交代码，都能帮助这个项目变得更好。

### 如何贡献

1. **Fork 本仓库**
   ```bash
   # 在 GitHub 上点击 Fork 按钮
   ```

2. **克隆你的 Fork**
   ```bash
   git clone https://github.com/你的用户名/server.git
   cd server
   ```

3. **创建特性分支**
   ```bash
   git checkout -b feature/amazing-feature
   ```

4. **进行修改并提交**
   ```bash
   git add .
   git commit -m "Add some amazing feature"
   ```

5. **推送到你的 Fork**
   ```bash
   git push origin feature/amazing-feature
   ```

6. **创建 Pull Request**
   - 在 GitHub 上打开你的 Fork
   - 点击 "New Pull Request"
   - 描述你的更改

### 贡献重点领域

我们特别欢迎以下方面的贡献：

- 🔒 **安全性改进**：加密算法优化、安全漏洞修复
- 🔐 **隐私保护增强**：更好的数据匿名化方案
- 📝 **文档完善**：API 文档、使用教程、最佳实践
- 🐛 **Bug 修复**：发现和修复问题
- ✨ **新功能开发**：新的推送类型、管理功能等
- 🧪 **测试覆盖**：单元测试、集成测试
- 🌍 **国际化**：多语言支持
- 🎨 **UI/UX**：管理界面改进

### 代码提交规范

请遵循以下提交信息格式：

```
<类型>: <简短描述>

<详细描述>

<相关 Issue>
```

**类型**：
- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 代码重构
- `test`: 测试相关
- `chore`: 构建/工具相关

**示例**：
```
feat: 添加批量推送API

实现了同时向多个设备发送推送的功能，
支持最多100个设备的批量操作。

Closes #123
```

### 开发流程

1. **确保代码通过测试**
   ```bash
   go test ./...
   ```

2. **格式化代码**
   ```bash
   go fmt ./...
   go vet ./...
   ```

3. **更新文档**
   - 如果添加新功能，更新 README.md
   - 如果修改 API，更新 API 文档

4. **提交前检查**
   - [ ] 代码已格式化
   - [ ] 测试已通过
   - [ ] 文档已更新
   - [ ] 提交信息清晰

### 报告问题

发现 bug？请[创建 Issue](https://github.com/dengdeng-harmonyos/server/issues/new) 并包含：

- 🔍 **问题描述**：清晰描述遇到的问题
- 📋 **复现步骤**：如何触发这个问题
- 💻 **环境信息**：OS、Go 版本、Docker 版本等
- 📸 **截图/日志**：如果适用

### 功能建议

有新想法？请[创建 Feature Request](https://github.com/dengdeng-harmonyos/server/issues/new) 并说明：

- 💡 **功能描述**：你想要什么功能
- 🎯 **使用场景**：为什么需要这个功能
- 📝 **期望行为**：功能应该如何工作
- 🔄 **替代方案**：是否有其他解决方案

### 行为准则

- ✅ 尊重所有贡献者
- ✅ 保持友好和专业
- ✅ 接受建设性批评
- ✅ 关注项目的整体利益
- ❌ 不允许骚扰或歧视性言论

### 获得帮助

遇到问题？可以通过以下方式获得帮助：

- 📖 查看[文档](README.md)
- 💬 在 [Issues](https://github.com/dengdeng-harmonyos/server/issues) 中提问
- 🔍 搜索已有的 Issues 和 Pull Requests

## 📄 开源协议

本项目采用 **MIT 协议**开源，详见 [LICENSE](LICENSE) 文件。

### 许可说明

- ✅ 可以商业使用
- ✅ 可以修改源代码
- ✅ 可以分发
- ✅ 可以私用
- ⚠️ 需要包含许可证和版权声明
- ⚠️ 不提供责任担保

---

## 🌟 致谢

感谢所有为这个项目做出贡献的人！

### 技术支持

- [HarmonyOS Next](https://developer.harmonyos.com/) - 鸿蒙操作系统
- [Huawei Push Kit](https://developer.huawei.com/consumer/cn/hms/huawei-pushkit/) - 华为推送服务
- [Gin Web Framework](https://gin-gonic.com/) - Go Web 框架
- [PostgreSQL](https://www.postgresql.org/) - 开源数据库

### 贡献者

<a href="https://github.com/dengdeng-harmonyos/server/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=dengdeng-harmonyos/server" />
</a>

---

## 📞 联系与支持

### 项目链接

- 🏠 **项目主页**: [https://github.com/dengdeng-harmonyos/server](https://github.com/dengdeng-harmonyos/server)
- 🐛 **问题反馈**: [GitHub Issues](https://github.com/dengdeng-harmonyos/server/issues)
- 🐳 **Docker 镜像**: [Docker Hub](https://hub.docker.com/r/ricwang/dengdeng-server)
- 📖 **文档**: [README](README.md) | [English](README_EN.md)

### 获取帮助

- 💬 通过 [GitHub Issues](https://github.com/dengdeng-harmonyos/server/issues) 提问
- 📧 发送邮件到项目维护者
- ⭐ 给项目一个 Star，关注最新动态

---

## 📊 项目状态

### 统计数据

[![GitHub stars](https://img.shields.io/github/stars/dengdeng-harmonyos/server?style=social)](https://github.com/dengdeng-harmonyos/server)
[![GitHub forks](https://img.shields.io/github/forks/dengdeng-harmonyos/server?style=social)](https://github.com/dengdeng-harmonyos/server/fork)
[![GitHub issues](https://img.shields.io/github/issues/dengdeng-harmonyos/server)](https://github.com/dengdeng-harmonyos/server/issues)
[![GitHub pull requests](https://img.shields.io/github/issues-pr/dengdeng-harmonyos/server)](https://github.com/dengdeng-harmonyos/server/pulls)
[![GitHub license](https://img.shields.io/github/license/dengdeng-harmonyos/server)](LICENSE)

### Star 历史

[![Star History Chart](https://api.star-history.com/svg?repos=dengdeng-harmonyos/server&type=Date)](https://star-history.com/#dengdeng-harmonyos/server&Date)

---

## ⚠️ 免责声明

**本服务提供推送基础设施，不存储任何用户数据。**

- 🔒 请确保你的加密密钥安全，不要与他人共享
- 🔐 请妥善保管华为推送服务配置文件
- 📝 请遵守当地法律法规和隐私保护政策
- ⚖️ 本项目不对使用本服务造成的任何后果负责
- 🛡️ 请定期更新依赖和安全补丁

---

## 💡 最后的话

如果这个项目对你有帮助，欢迎：

- ⭐ 给项目一个 Star
- 🔄 Fork 并参与贡献
- 📢 分享给更多的开发者
- 💬 反馈问题和建议

**让我们一起构建一个安全、可靠的推送服务！** 🚀
