# 当当当消息推送服务器

<div align="center">

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)
![Docker](https://img.shields.io/badge/Docker-ready-2496ED?logo=docker)

一个基于 Go 语言开发的开源消息推送服务器，支持对接华为推送服务，可自部署。

[快速开始](#快速开始) •
[部署指南](#部署指南) •
[API 文档](#api-文档) •
[贡献指南](#贡献)

</div>

---

## ✨ 特性

- 🚀 **高性能**: 基于 Go 语言和 Gin 框架，支持高并发请求
- 📱 **华为推送**: 完整对接华为推送服务 API
- 🐳 **Docker 支持**: 提供 Docker 和 Docker Compose 一键部署
- 💾 **PostgreSQL**: 使用 PostgreSQL 数据库，稳定可靠
- 📊 **统计分析**: 内置推送记录和统计分析功能
- 🔧 **易于配置**: 通过环境变量灵活配置
- 📖 **开源免费**: MIT 协议，完全开源

## 📋 功能列表

- ✅ 设备 Token 注册与管理
- ✅ 单播推送（向单个设备推送）
- ✅ 群播推送（向多个设备推送）
- ✅ 广播推送（向所有设备推送）
- ✅ 推送记录查询
- ✅ 推送统计分析
- ✅ RESTful API 接口

## 🏗️ 架构

```
┌─────────────┐      ┌──────────────────┐      ┌─────────────────┐
│  鸿蒙应用    │ ←──→ │   推送服务器      │ ←──→ │  华为推送服务   │
│  (客户端)    │      │  (Go Backend)    │      │  (HMS Push)     │
└─────────────┘      └──────────────────┘      └─────────────────┘
                              ↓
                     ┌──────────────────┐
                     │   PostgreSQL     │
                     │   (数据库)        │
                     └──────────────────┘
```

## 🚀 快速开始

### 前置要求

- Go 1.21 或更高版本
- PostgreSQL 15 或更高版本
- Docker 和 Docker Compose (可选)
- 华为开发者账号和推送服务凭证

### 方式一：Docker Compose 部署（推荐）

1. **克隆仓库**
   ```bash
   git clone https://github.com/yourusername/dangdangdang-push-server.git
   cd dangdangdang-push-server
   ```

2. **配置环境变量**
   ```bash
   cp .env.example .env
   # 编辑 .env 文件，填入你的华为推送配置
   nano .env
   ```

3. **启动服务**
   ```bash
   docker-compose up -d
   ```

4. **检查服务状态**
   ```bash
   docker-compose ps
   curl http://localhost:8080/health
   ```

### 方式二：本地开发部署

1. **安装 PostgreSQL**
   ```bash
   # macOS
   brew install postgresql@15
   brew services start postgresql@15

   # Linux (Ubuntu/Debian)
   sudo apt-get install postgresql-15
   ```

2. **创建数据库**
   ```bash
   createdb push_server
   ```

3. **配置环境变量**
   ```bash
   cp .env.example .env
   # 编辑 .env 文件
   ```

4. **安装依赖**
   ```bash
   go mod download
   ```

5. **运行服务**
   ```bash
   go run cmd/server/main.go
   ```

服务将在 `http://localhost:8080` 启动。

## 📝 配置说明

### 华为推送配置

1. 登录 [AppGallery Connect](https://developer.huawei.com/consumer/cn/service/josp/agc/index.html)
2. 创建应用并开启推送服务
3. 下载 `agconnect-services.json` 文件并保存到 `config/` 目录
4. 在 `.env` 文件中配置：
   ```env
   HUAWEI_PROJECT_ID=your_project_id  # 从agconnect-services.json的client.project_id获取
   HUAWEI_SERVICE_ACCOUNT_FILE=./config/agconnect-services.json
   PUSH_TOKEN_ENCRYPTION_KEY=your_32_byte_key
   ```

### 环境变量

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `PORT` | 服务端口 | `8080` |
| `GIN_MODE` | Gin 运行模式 | `debug` |
| `DB_HOST` | 数据库主机 | `localhost` |
| `DB_PORT` | 数据库端口 | `5432` |
| `DB_USER` | 数据库用户 | `postgres` |
| `DB_PASSWORD` | 数据库密码 | - |
| `DB_NAME` | 数据库名称 | `push_server` |
| `HUAWEI_PROJECT_ID` | 华为项目ID（从agconnect-services.json获取） | - |
| `HUAWEI_SERVICE_ACCOUNT_FILE` | AGConnect配置文件路径 | `./config/agconnect-services.json` |
| `PUSH_TOKEN_ENCRYPTION_KEY` | Push Token加密密钥（32字节） | - |

## 📡 API 文档

### 设备管理

#### 注册设备
```http
POST /api/device/register
Content-Type: application/json

{
  "push_token": "设备推送Token",
  "device_id": "设备唯一标识",
  "device_type": "phone",
  "os_version": "HarmonyOS 4.0",
  "app_version": "1.0.0"
}
```

#### 更新设备信息
```http
PUT /api/device/update
Content-Type: application/json

{
  "push_token": "设备推送Token",
  "device_id": "设备唯一标识",
  "os_version": "HarmonyOS 4.1"
}
```

#### 注销设备
```http
DELETE /api/device/unregister?push_token=xxxxx
```

### 推送消息

#### 单播推送
```http
POST /api/push/single
Content-Type: application/json

{
  "push_token": "目标设备Token",
  "message": {
    "title": "消息标题",
    "content": "消息内容",
    "data": {
      "key": "value"
    }
  }
}
```

#### 群播推送
```http
POST /api/push/multiple
Content-Type: application/json

{
  "push_tokens": ["token1", "token2", "token3"],
  "message": {
    "title": "消息标题",
    "content": "消息内容"
  }
}
```

#### 广播推送
```http
POST /api/push/all
Content-Type: application/json

{
  "message": {
    "title": "消息标题",
    "content": "消息内容"
  }
}
```

### 查询接口

#### 获取推送记录
```http
GET /api/query/records?limit=50&offset=0
```

#### 获取推送统计
```http
GET /api/query/statistics
```

## 🗄️ 数据库表结构

### users（用户表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | SERIAL | 主键 |
| username | VARCHAR(100) | 用户名 |
| phone | VARCHAR(20) | 手机号 |
| created_at | TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | 更新时间 |

### devices（设备表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | SERIAL | 主键 |
| user_id | INTEGER | 用户ID |
| push_token | VARCHAR(500) | 推送Token |
| device_id | VARCHAR(200) | 设备ID |
| device_type | VARCHAR(50) | 设备类型 |
| os_version | VARCHAR(50) | 系统版本 |
| app_version | VARCHAR(50) | 应用版本 |
| is_active | BOOLEAN | 是否激活 |
| created_at | TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | 更新时间 |

### push_records（推送记录表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | SERIAL | 主键 |
| user_id | INTEGER | 用户ID |
| device_id | INTEGER | 设备ID |
| title | VARCHAR(200) | 消息标题 |
| content | TEXT | 消息内容 |
| data | JSONB | 附加数据 |
| status | VARCHAR(50) | 推送状态 |
| error_message | TEXT | 错误信息 |
| sent_at | TIMESTAMP | 发送时间 |
| clicked_at | TIMESTAMP | 点击时间 |

## 🔧 开发

### 项目结构
```
server/
├── cmd/
│   └── server/
│       └── main.go           # 应用入口
├── internal/
│   ├── config/
│   │   └── config.go         # 配置管理
│   ├── database/
│   │   └── database.go       # 数据库连接
│   ├── handler/
│   │   ├── device.go         # 设备管理处理器
│   │   └── push.go           # 推送处理器
│   ├── middleware/
│   │   └── middleware.go     # 中间件
│   ├── models/
│   │   └── models.go         # 数据模型
│   └── service/
│       └── huawei_push.go    # 华为推送服务
├── .env.example              # 环境变量示例
├── .gitignore
├── docker-compose.yml        # Docker Compose 配置
├── Dockerfile                # Docker 镜像配置
├── go.mod
├── go.sum
└── README.md
```

### 运行测试
```bash
go test ./...
```

### 构建
```bash
go build -o push-server cmd/server/main.go
```

## 🚢 生产部署

### 使用 Docker

```bash
# 构建镜像
docker build -t dangdangdang-push-server .

# 运行容器
docker run -d \
  --name push-server \
  -p 8080:8080 \
  -v $(pwd)/config:/app/config \
  -e DB_HOST=your-db-host \
  -e DB_PASSWORD=your-db-password \
  -e HUAWEI_PROJECT_ID=your-project-id \
  -e PUSH_TOKEN_ENCRYPTION_KEY=your-32-byte-key \
  dangdangdang-push-server
```

### 使用 Systemd

创建 systemd 服务文件 `/etc/systemd/system/push-server.service`:

```ini
[Unit]
Description=Push Server
After=network.target postgresql.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/push-server
ExecStart=/opt/push-server/push-server
Restart=on-failure
EnvironmentFile=/opt/push-server/.env

[Install]
WantedBy=multi-user.target
```

启动服务：
```bash
sudo systemctl daemon-reload
sudo systemctl enable push-server
sudo systemctl start push-server
```

## 🤝 贡献

欢迎贡献代码、报告问题或提出建议！

1. Fork 本仓库
2. 创建你的特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交你的更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启一个 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🙏 致谢

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [PostgreSQL](https://www.postgresql.org/)
- [Huawei Push Kit](https://developer.huawei.com/consumer/cn/hms/huawei-pushkit)

## 📧 联系方式

- 问题反馈: [GitHub Issues](https://github.com/yourusername/dangdangdang-push-server/issues)
- 邮箱: your.email@example.com

## 🗺️ 路线图

- [ ] 支持更多推送服务商（小米、OPPO、VIVO 等）
- [ ] 添加 Web 管理后台
- [ ] 支持定时推送任务
- [ ] 添加消息模板管理
- [ ] 支持 A/B 测试推送
- [ ] 增强统计和分析功能

---

<div align="center">
  Made with ❤️ by the Community
</div>
