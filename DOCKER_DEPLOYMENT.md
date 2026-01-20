# Docker部署配置说明

## 概述

本项目使用单容器方案，将PostgreSQL和应用服务打包在一起，简化部署和维护。

## 关键设计

### 1. 数据库迁移在容器启动时执行

数据库迁移**不是在镜像构建时执行**，而是在**容器启动时自动执行**。

#### 为什么这样设计？

✅ **用户升级镜像后自动迁移** - 拉取新镜像后，重启容器会自动应用新的数据库变更  
✅ **幂等性保证** - 已执行的迁移不会重复执行  
✅ **数据持久化** - 数据库数据存储在volume中，升级镜像不影响数据  
✅ **灵活性** - 支持手动触发迁移，不依赖镜像重建  

#### 执行流程

```
docker run/restart
    ↓
/start.sh 执行
    ↓
启动 PostgreSQL
    ↓
执行 /app/database/migrate.sh  ← 这里执行迁移
    ↓
启动应用服务
```

### 2. 内部数据库配置已写死

PostgreSQL作为内部依赖，配置已在Dockerfile中写死，用户无需关注。

#### 内部配置（在Dockerfile中）

```dockerfile
ENV POSTGRES_DB=push_server
ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=postgres123
ENV DB_HOST=localhost
ENV DB_PORT=5432
```

这些配置：
- ✅ 不需要暴露给用户
- ✅ 不建议修改（内部使用）
- ✅ 保证一致性

### 3. 用户可配置项（在docker-compose.yml中）

```yaml
environment:
  # 应用服务端口
  - PORT=8081
  
  # 运行模式
  - GIN_MODE=release
  
  # 服务器名称
  - SERVER_NAME=噔噔推送服务
  
  # 安全密钥（必需）
  - PUSH_TOKEN_ENCRYPTION_KEY=your-32-character-key
```

## 使用方法

### 首次部署

```bash
# 1. 设置加密密钥
export PUSH_TOKEN_ENCRYPTION_KEY="your-32-character-encryption-key-here"

# 2. 启动容器（自动初始化数据库）
docker-compose up -d

# 3. 查看日志
docker logs -f push-server
```

### 升级镜像版本

```bash
# 1. 拉取新镜像
docker-compose pull

# 2. 重启容器（自动执行数据库迁移）
docker-compose up -d

# 3. 验证迁移
docker logs push-server | grep MIGRATE
```

预期输出：
```
[MIGRATE] Running database migrations...
[INFO] Database already initialized
[INFO] Checking for pending migrations...
[INFO] Applying migration: 20260120100000 - add_device_metadata
[INFO] Migration 20260120100000 applied successfully
[MIGRATE] Database migration completed successfully
```

### 查看迁移状态

```bash
docker exec -it push-server psql -U postgres -d push_server -c \
  "SELECT version, description, applied_at FROM schema_migrations ORDER BY version;"
```

## 数据持久化

数据库数据存储在Docker volume中：

```yaml
volumes:
  - postgres_data:/var/lib/postgresql/data
```

### 特点
- ✅ 升级镜像不丢失数据
- ✅ 容器删除后数据保留
- ✅ 支持备份和恢复

### 备份

```bash
# 备份数据库
docker exec push-server pg_dump -U postgres push_server > backup.sql

# 或备份整个volume
docker run --rm -v push_server_postgres_data:/data -v $(pwd):/backup alpine \
  tar czf /backup/postgres_backup_$(date +%Y%m%d).tar.gz /data
```

### 恢复

```bash
# 从SQL恢复
cat backup.sql | docker exec -i push-server psql -U postgres -d push_server

# 从volume备份恢复
docker run --rm -v push_server_postgres_data:/data -v $(pwd):/backup alpine \
  tar xzf /backup/postgres_backup_20260120.tar.gz -C /
```

## 完全重置（开发环境）

```bash
# 停止并删除容器和数据
docker-compose down -v

# 重新启动（重新初始化）
docker-compose up -d
```

## 配置说明

### 必需配置

| 变量 | 说明 | 默认值 |
|------|------|--------|
| PUSH_TOKEN_ENCRYPTION_KEY | 推送令牌加密密钥（32字符） | 无，必须设置 |

### 可选配置

| 变量 | 说明 | 默认值 |
|------|------|--------|
| PORT | 应用服务端口 | 8080 |
| GIN_MODE | 运行模式（debug/release） | release |
| SERVER_NAME | 服务器标识名称 | 噔噔推送服务 |

### 内部配置（无需修改）

| 变量 | 说明 | 值 |
|------|------|-----|
| POSTGRES_DB | 数据库名 | push_server |
| POSTGRES_USER | 数据库用户 | postgres |
| POSTGRES_PASSWORD | 数据库密码 | postgres123 |
| DB_HOST | 数据库主机 | localhost |
| DB_PORT | 数据库端口 | 5432 |

## 端口说明

| 端口 | 说明 | 暴露 |
|------|------|------|
| 8080/8081 | HTTP API服务 | ✅ 是 |
| 5432 | PostgreSQL | ❌ 否（内部使用） |

## 健康检查

容器自动进行健康检查：

```yaml
healthcheck:
  test: ["CMD-SHELL", "pg_isready -U postgres && wget --spider http://localhost:8081/health"]
  interval: 30s
  timeout: 10s
  start_period: 40s
  retries: 3
```

检查内容：
1. PostgreSQL运行状态
2. 应用服务健康接口

## 故障排除

### 容器无法启动

```bash
# 查看详细日志
docker logs push-server

# 检查PostgreSQL状态
docker exec push-server pg_isready -U postgres

# 检查应用服务
docker exec push-server wget --spider http://localhost:8080/health
```

### 迁移失败

```bash
# 查看迁移日志
docker logs push-server | grep -A 20 MIGRATE

# 手动执行迁移
docker exec push-server /app/database/migrate.sh

# 查看迁移状态
docker exec -it push-server psql -U postgres -d push_server -c \
  "SELECT * FROM schema_migrations ORDER BY version;"
```

### 数据丢失

如果数据丢失，检查：

```bash
# 列出所有volume
docker volume ls

# 检查volume是否存在
docker volume inspect push_server_postgres_data

# 如果volume丢失，可能需要从备份恢复
```

## 生产环境建议

1. **设置强密钥**
   ```bash
   # 生成32字符随机密钥
   openssl rand -base64 24
   ```

2. **使用环境变量文件**
   ```bash
   # 创建 .env 文件
   cat > .env <<EOF
   PUSH_TOKEN_ENCRYPTION_KEY=your-generated-key-here
   EOF
   
   # docker-compose会自动读取
   docker-compose up -d
   ```

3. **定期备份**
   ```bash
   # 添加到crontab
   0 2 * * * docker exec push-server pg_dump -U postgres push_server > /backup/push_server_$(date +\%Y\%m\%d).sql
   ```

4. **监控日志**
   ```bash
   # 使用日志聚合工具
   docker logs -f push-server | tee /var/log/push-server.log
   ```

5. **限制资源**
   ```yaml
   services:
     push-server:
       deploy:
         resources:
           limits:
             cpus: '1.0'
             memory: 1G
           reservations:
             cpus: '0.5'
             memory: 512M
   ```

## 相关文档

- [database/QUICKSTART.md](database/QUICKSTART.md) - 数据库迁移快速指南
- [database/MIGRATIONS.md](database/MIGRATIONS.md) - 迁移系统详细文档
- [DATABASE_MIGRATION_COMPLETE.md](DATABASE_MIGRATION_COMPLETE.md) - 迁移系统实施说明

## 技术架构

```
┌─────────────────────────────────────┐
│      Docker Container               │
│                                     │
│  ┌──────────────┐  ┌─────────────┐ │
│  │ PostgreSQL   │  │   Go App    │ │
│  │   :5432      │  │   :8080     │ │
│  │              │  │             │ │
│  │  数据库      │←─│  应用服务   │ │
│  └──────────────┘  └─────────────┘ │
│          ↓                          │
│  ┌──────────────┐                  │
│  │   Volume     │                  │
│  │  postgres_   │                  │
│  │    data      │                  │
│  └──────────────┘                  │
└─────────────────────────────────────┘
         ↑
    迁移脚本在启动时执行
    /app/database/migrate.sh
```
