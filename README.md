# 噔噔推送服务 (Dengdeng Push Server)

[![GitHub release](https://img.shields.io/github/v/release/dengdeng-harmonyos/server)](https://github.com/dengdeng-harmonyos/server/releases)
[![GitHub stars](https://img.shields.io/github/stars/dengdeng-harmonyos/server?style=social)](https://github.com/dengdeng-harmonyos/server)
[![License](https://img.shields.io/github/license/dengdeng-harmonyos/server)](LICENSE)
[![Docker Pulls](https://img.shields.io/docker/pulls/ricwang/dengdeng-server)](https://hub.docker.com/r/ricwang/dengdeng-server)
[![Go Version](https://img.shields.io/github/go-mod/go-version/dengdeng-harmonyos/server)](go.mod)

## 📖 项目简介

噔噔推送服务是一个专为 **HarmonyOS Next** 设计的**安全、隐私友好**的推送服务解决方案。本项目完全开源，致力于为开发者提供一个可信赖、易部署的推送服务基础设施。

> 🎯 **v1.1 正式发布**：强化推送可用性、设备诊断和消息同步可靠性

### ✨ 主要亮点

- **🚀 一键部署**：单个 Docker 容器即可运行，内置 PostgreSQL 数据库
- **🔐 安全优先**：配置编译时嵌入，支持 AES-256-GCM 加密
- **📦 零依赖**：无需外部配置文件，开箱即用
- **🤖 CI/CD 自动化**：GitHub Actions 自动构建和部署

### 🔒 安全与隐私承诺

- **🔑 RSA 公钥支持**：消息使用端到端消息加密，app端同步完即删除服务端加密数据
- **🔐 端到端加密**：Push Token 使用 AES-256-GCM 加密存储
- **🎭 匿名化设计**：使用随机生成的 Device Id，与真实设备无关联
- **🔑 配置编译时嵌入**：敏感配置在构建时嵌入二进制文件，无需配置文件
- **🛡️ 开源透明**：所有代码公开，接受社区审查

## ✨ 核心特性

### 🚀 部署与运维

- **📦 单容器部署**：包含 PostgreSQL + 推送服务，开箱即用
- **🔧 配置嵌入**：华为推送配置编译时嵌入，无需外部文件
- **🤖 自动化 CI/CD**：GitHub Actions 自动构建、测试和部署
- **🏥 健康检查**：内置健康检查接口，支持监控
- **🧭 设备诊断**：提供非敏感设备状态检查，便于排查注册与同步问题
- **🐳 Docker 支持**：官方镜像托管在 Docker Hub
- **🔄 自动重启**：容器崩溃自动恢复

### 📡 功能特性

- **📬 通知推送**：支持通知栏消息（带标题、内容、自定义数据）
- **🔄 加密消息同步**：消息加密暂存 30 天，App 同步确认后删除
- **🏥 健康监控**：内置健康检查和服务状态接口

## 🚀 快速开始

### 方式一：使用一键运行脚本（推荐）

```bash
curl -sSL https://raw.githubusercontent.com/dengdeng-harmonyos/server/refs/heads/release/deploy.sh | bash
```

### 方式二：使用 Docker Hub 镜像

这是最简单快速的部署方式：

#### 1. 生成加密密钥

```bash
# 生成 32 字节随机密钥（Base64 编码）
openssl rand -base64 32
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
  -e SERVER_NAME=你的自定义服务名称 \
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

## 📡 API 接口

### 快速概览

所有接口使用简单的 HTTP GET 请求，无需复杂的认证流程。

| 功能 | 接口路径 | 说明 |
|------|---------|------|
| 健康检查 | `GET /health` | 检查服务状态 |
| 设备注册 | `POST /api/v1/device/register` | 注册设备获取 Device Id |
| 通知推送 | `GET /api/v1/push/notification` | 发送通知栏消息 |
| 设备诊断 | `GET /api/v1/diagnostics/device` | 查询非敏感设备状态 |

### 示例：发送通知

```bash
curl "http://your-server:8080/api/v1/push/notification?device_id=YOUR_DEVICE_KEY&title=测试消息&content=这是一条测试推送"
```

### 示例：点击通知后打开链接或 App

`data` 中 key 为 `__url` 的项会在用户点击通知后打开，支持网页链接和合法 App URL Scheme。手写请求时请对参数做 URL 编码；`file`、`javascript`、`data`、`content`、`tel`、`sms`、`mailto` 等高风险 scheme 会被拒绝。

```bash
curl --get "http://your-server:8080/api/v1/push/notification" \
  --data-urlencode "device_id=YOUR_DEVICE_KEY" \
  --data-urlencode "title=测试消息" \
  --data-urlencode "content=点击通知后跳转" \
  --data-urlencode 'data=[{"key":"__url","value":"myapp://page/detail?id=1"}]'
```

### 示例：设备诊断

```bash
curl "http://your-server:8080/api/v1/diagnostics/device?device_id=YOUR_DEVICE_KEY"
```

诊断接口只返回设备是否存在、是否有公钥、是否活跃、最近活跃时间和待同步消息数；不会返回 Push Token、公钥内容、消息内容，也不提供聚合统计数据。

> 批量推送、后台消息和表单刷新接口当前未开放；请使用单条通知接口。

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
| `SERVER_VERSION` | 服务端版本号，用于 App 兼容性检查 | ❌ | `1.1.2` |
| `SERVER_API_VERSION` | 服务端 API 兼容版本 | ❌ | `3` |
| `SERVER_CAPABILITIES` | 服务端能力列表，逗号分隔 | ❌ | `message_crypto_v1,push_url_data,push_deep_link_scheme,background_push_wake,app_update_policy,device_diagnostics` |
| `SERVER_UPGRADE_URL` | App 提示用户升级服务端时展示的地址 | ❌ | `https://github.com/dengdeng-harmonyos/server` |
| `PORT` | HTTP 服务端口 | ❌ | `8080` |

`GET /health` 会返回 `version`、`apiVersion`、`capabilities` 和 `upgradeUrl`。App 会用这些字段判断自部署服务端是否支持当前 App 功能；如果版本过低，用户需要更新服务端镜像或源码后再继续使用该服务端。

### App 更新策略

App 强制更新策略存储在数据库表 `app_update_policies` 中。发布 App 新版本时，更新仓库内的 `config/app_update_policy.json`，新服务端镜像启动后会自动将该版本写入 `app_update_releases` 历史表；当清单中 `enabled=true` 时，还会把该版本同步为 `app_update_policies` 当前生效策略。

同一个镜像重复启动不会产生重复版本记录：`app_update_releases` 使用 `(platform, version_code)` 作为主键，`app_update_policies` 使用 `platform` 作为主键。旧镜像重启也不会把当前策略降级到更低版本。

发布清单示例：

```json
{
  "platform": "harmonyos",
  "versionCode": 1000102,
  "versionName": "1.1.2",
  "minVersionCode": 1000102,
  "forceUpdate": true,
  "storeUrl": "store://appgallery.huawei.com/app/detail?id=top.yidingyaojizhu.dengdeng",
  "releaseNotes": "新增后台 Push 唤醒同步，有待收消息时可更早同步到本机；旧版本 App 需要升级后继续使用。",
  "enabled": true
}
```

也可以直接修改当前策略表：

```sql
UPDATE app_update_policies
SET latest_version_code = 1000102,
    latest_version_name = '1.1.2',
    min_version_code = 1000102,
    force_update = true,
    release_notes = '新增后台 Push 唤醒同步，有待收消息时可更早同步到本机；旧版本 App 需要升级后继续使用。',
    updated_at = CURRENT_TIMESTAMP
WHERE platform = 'harmonyos';
```

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
2. **待同步消息**（加密暂存）
   - 仅保存 RSA/AES 加密后的消息内容
   - 默认 30 天过期，App 同步确认后即从服务端删除
   - 服务启动后立即清理过期消息，并每 6 小时重复清理

## 🏗️ 架构设计

### 系统架构

```
┌─────────────────────────────────────────────────────────┐
│                    Docker 容器                           │
│  ┌─────────────────┐         ┌──────────────────┐      │
│  │  PostgreSQL 15  │ ←────→  │  推送服务 (Go)    │      │
│  │  - 设备信息      │         │  - Gin Web框架    │      │
│  │  - 加密Token    │         │  - AES-256加密    │      │
│  │                 │         │  - 华为推送API    │      │
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

### 3. 数据库安全

**定期备份**
```bash
# 创建备份
docker exec push-server pg_dump -U postgres push_server > backup-$(date +%Y%m%d).sql

# 自动备份脚本（添加到 crontab）
0 2 * * * docker exec push-server pg_dump -U postgres push_server | gzip > /backup/push-$(date +\%Y\%m\%d).sql.gz
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

## 📦 Docker 镜像

### 官方镜像

🐳 **Docker Hub**: [ricwang/dengdeng-server](https://hub.docker.com/r/ricwang/dengdeng-server)

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
├── docker-compose.yml          # Docker Compose 配置
└── README.md
```

## 📄 开源协议

本项目采用 **MIT 协议**开源，详见 [LICENSE](LICENSE) 文件。

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

### Star 历史

[![Star History Chart](https://api.star-history.com/svg?repos=dengdeng-harmonyos/server&type=Date)](https://star-history.com/#dengdeng-harmonyos/server&Date)

---

## ⚠️ 免责声明

**本服务提供推送基础设施，不存储任何用户数据。**

- 🔒 请确保你的加密密钥安全，不要与他人共享
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
