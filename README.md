# 噔噔推送服务 (Dengdeng Push Server)

[![GitHub stars](https://img.shields.io/github/stars/dengdeng-harmenyos/server?style=social)](https://github.com/dengdeng-harmenyos/server)
[![License](https://img.shields.io/github/license/dengdeng-harmenyos/server)](LICENSE)
[![Docker Pulls](https://img.shields.io/docker/pulls/ricwang/dengdeng-server)](https://hub.docker.com/r/ricwang/dengdeng-server)
[![Go Version](https://img.shields.io/github/go-mod/go-version/dengdeng-harmenyos/server)](go.mod)

[English](README_EN.md) | 简体中文

## 📖 项目简介

噔噔推送服务是一个专为 HarmonyOS Next 设计的**安全、隐私友好**的推送服务解决方案。本项目是完全开源的，致力于为开发者提供一个可信赖的推送服务基础设施。

### 🔒 安全与隐私承诺

- **🚫 零用户数据保存**：不存储任何推送消息内容，仅保存匿名统计数据
- **🔐 端到端加密**：Push Token 使用 AES-256-GCM 加密存储
- **🎭 匿名化设计**：使用随机生成的 Device Key，与设备无任何关联
- **📊 统计数据脱敏**：仅保存推送成功/失败次数，不记录具体内容
- **🔑 客户端密钥管理**：加密密钥由用户自行生成和保管
- **🛡️ 开源透明**：所有代码公开，接受社区审查

## ✨ 核心特性

### 安全性
- 🔒 AES-256-GCM 加密算法保护敏感数据
- 🎲 使用加密安全的随机数生成器 (crypto/rand)
- 🔐 RSA 公钥支持（可选，用于客户端消息加密）
- ⏱️ Device Key 时效性管理，自动失效机制
- 🚦 速率限制，防止滥用（每设备每日最大推送数限制）

### 隐私保护
- 📝 不存储推送消息内容
- 🎭 设备标识完全匿名化
- 📊 统计数据聚合，无法追溯到具体设备
- 🗑️ 定期清理过期设备记录
- 🔍 数据库字段最小化原则

### 功能特性
- 📬 支持通知消息推送 (Notification)
- 🃏 支持卡片刷新 (Form Update)
- 🔄 支持后台推送 (Background Push)
- 📦 支持批量推送
- ⚡ 单容器部署，包含 PostgreSQL 和推送服务
- 🏥 健康检查接口
- 📊 推送统计（仅统计数据）

## 🚀 快速开始

### 使用 Docker（推荐）

这是最简单的部署方式，只需要两个步骤：

#### 1. 生成加密密钥

```bash
# 生成32字节随机密钥（Base64编码）
openssl rand -base64 32
```

将生成的密钥保存到 `.env` 文件：

```bash
echo "PUSH_TOKEN_ENCRYPTION_KEY=你生成的密钥" > .env
```

#### 2. 启动服务

使用 Docker Hub 的镜像直接启动：

```bash
# 拉取最新镜像
docker pull ricwang/dengdeng-server:latest

# 启动服务
docker run -d \
  --name push-server \
  -p 8080:8080 \
  -e PUSH_TOKEN_ENCRYPTION_KEY=你的密钥 \
  -e SERVER_NAME=噔噔推送服务 \
  -v push-data:/var/lib/postgresql/data \
  --restart unless-stopped \
  ricwang/dengdeng-server:latest
```

或者使用 docker-compose：

```yaml
services:
  push-server:
    image: ricwang/dengdeng-server:latest
    container_name: push-server
    environment:
      - PUSH_TOKEN_ENCRYPTION_KEY=${PUSH_TOKEN_ENCRYPTION_KEY}
      - SERVER_NAME=噔噔推送服务
    ports:
      - "8080:8080"
    volumes:
      - push-data:/var/lib/postgresql/data
    restart: unless-stopped

volumes:
  push-data:
```

启动命令：

```bash
docker-compose up -d
```

#### 3. 验证服务

```bash
# 检查健康状态
curl http://localhost:8080/health

# 查看日志
docker logs -f push-server
```

### 从源码构建

如果你想从源码构建：

```bash
# 克隆仓库
git clone https://github.com/dengdeng-harmenyos/server.git
cd server

# 生成密钥
echo "PUSH_TOKEN_ENCRYPTION_KEY=$(openssl rand -base64 32)" > .env

# 使用本地 Dockerfile 构建并启动
docker-compose up -d --build
```

## 📡 API 接口

完整的 API 接口文档请查看：[API.md](API.md)

主要接口包括：
- **设备注册** - 注册设备并获取 Device Key
- **推送通知** - 发送通知栏消息
- **卡片刷新** - 更新 HarmonyOS 卡片
- **后台推送** - 发送后台数据
- **批量推送** - 同时向多个设备推送
- **健康检查** - 检查服务状态

## 🔧 配置说明

### 必需配置

| 环境变量 | 说明 | 示例 |
|---------|------|------|
| `PUSH_TOKEN_ENCRYPTION_KEY` | Push Token加密密钥（32字节） | `openssl rand -base64 32` 生成 |

### 可选配置

| 环境变量 | 说明 | 默认值 |
|---------|------|--------|
| `SERVER_NAME` | 服务器名称 | `噔噔推送服务` |
| `PORT` | 服务端口 | `8080` |
| `GIN_MODE` | 运行模式 | `release` |
| `DEVICE_KEY_TTL` | Device Key有效期（秒） | `31536000` (1年) |
| `MAX_DAILY_PUSH_PER_DEVICE` | 每设备每日最大推送数 | `1000` |

### 内部数据库配置（无需修改）

容器内部自动配置 PostgreSQL，用户无需关心数据库设置。

## 📊 数据存储说明

### 存储的数据

1. **设备信息**（匿名化）
   - Device Key（随机生成）
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

### 安全架构

```
客户端应用
    ↓ (注册请求)
设备注册
    ↓ (生成 Device Key)
Push Token 加密存储 (AES-256-GCM)
    ↓
推送请求 → Device Key 验证
    ↓
Push Token 解密
    ↓
华为推送服务
```

### 数据流

1. **设备注册**：生成随机 Device Key，加密存储 Push Token
2. **推送请求**：使用 Device Key 查找设备，解密 Push Token
3. **推送执行**：调用华为推送 API，**不保存消息内容**
4. **统计记录**：仅记录成功/失败次数

## 🔐 安全最佳实践

1. **密钥管理**
   - 使用强随机密钥（32字节）
   - 定期轮换加密密钥
   - 不要在代码中硬编码密钥
   - 使用环境变量或密钥管理服务

2. **网络安全**
   - 使用 HTTPS/TLS 加密传输
   - 配置防火墙规则
   - 启用速率限制

3. **数据库安全**
   - 定期备份数据
   - 使用持久化卷存储数据
   - 定期清理过期设备

4. **监控与审计**
   - 监控异常推送行为
   - 记录访问日志
   - 定期审查统计数据

## 📦 Docker 镜像

官方镜像托管在 Docker Hub：

🐳 **镜像地址**：[ricwang/dengdeng-server](https://hub.docker.com/r/ricwang/dengdeng-server)

### 可用标签

- `latest` - 最新稳定版本
- `v1.x.x` - 特定版本号

### 镜像说明

- 基础镜像：`postgres:15-alpine`
- 包含组件：PostgreSQL 15 + 推送服务
- 镜像大小：约 300MB
- 架构支持：`linux/amd64`

## 🛠️ 开发指南

### 本地开发

```bash
# 安装依赖
go mod download

# 运行数据库迁移
cd database
./migrate.sh

# 运行开发服务器
go run cmd/server/main.go
```

### 构建

```bash
# 编译二进制文件
go build -o bin/server cmd/server/main.go

# 构建 Docker 镜像
docker build -t dengdeng-server .
```

### 数据库迁移

```bash
# 创建新迁移
cd database
./create-migration.sh <migration_name>

# 执行迁移
./migrate.sh
```

## 🤝 贡献指南

我们欢迎各种形式的贡献！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

### 贡献重点

- 🔒 安全性改进
- 🔐 隐私保护增强
- 📝 文档完善
- 🐛 Bug 修复
- ✨ 新功能开发
- 🧪 测试覆盖

## 📄 开源协议

本项目采用 MIT 协议开源，详见 [LICENSE](LICENSE) 文件。

## 🌟 致谢

- HarmonyOS Next 开发团队
- 所有贡献者和支持者

## 📞 联系方式

- 项目主页：[https://github.com/dengdeng-harmenyos/server](https://github.com/dengdeng-harmenyos/server)
- 问题反馈：[GitHub Issues](https://github.com/dengdeng-harmenyos/server/issues)
- Docker 镜像：[Docker Hub](https://hub.docker.com/r/ricwang/dengdeng-server)

## 🎯 路线图

- [ ] 支持多种推送服务商
- [ ] Web 管理界面
- [ ] 更详细的推送统计
- [ ] 消息优先级队列
- [ ] 推送模板管理
- [ ] API 密钥管理系统
- [ ] 多语言 SDK

---

## 📊 项目趋势

[![Star History Chart](https://api.star-history.com/svg?repos=dengdeng-harmenyos/server&type=Date)](https://star-history.com/#dengdeng-harmenyos/server&Date)

**⚠️ 安全提示**：本服务仅提供推送基础设施，不存储任何用户数据。请确保你的加密密钥安全，不要与他人共享。

**💡 提示**：如果你觉得这个项目有帮助，欢迎给我们一个 ⭐️ Star！
