# Docker 部署指南

## 单容器部署（PostgreSQL + 应用）

### 快速启动

```bash
# 1. 确保配置文件存在
cp .env.example .env
# 编辑.env填入实际配置

# 2. 确保agconnect-services.json存在
# 将从AppGallery Connect下载的文件放到 config/agconnect-services.json

# 3. 构建并启动
docker-compose -f docker-compose.single.yml up -d

# 4. 查看日志
docker-compose -f docker-compose.single.yml logs -f

# 5. 停止服务
docker-compose -f docker-compose.single.yml down
```

### 直接使用Docker命令

```bash
# 构建镜像
docker build -t dangdangdang-push-server .

# 运行容器
docker run -d \
  --name push-server \
  -p 8080:8080 \
  -v $(pwd)/data:/var/lib/postgresql/data \
  -v $(pwd)/config/agconnect-services.json:/app/config/agconnect-services.json:ro \
  -e POSTGRES_PASSWORD=postgres123 \
  -e HUAWEI_PROJECT_ID=101653523863440882 \
  -e PUSH_TOKEN_ENCRYPTION_KEY=your_encryption_key \
  dangdangdang-push-server

# 查看日志
docker logs -f push-server

# 进入容器
docker exec -it push-server sh

# 停止容器
docker stop push-server

# 删除容器
docker rm push-server
```

## 多容器部署（推荐生产环境）

如果需要独立管理数据库和应用：

```bash
# 使用原始的docker-compose.yml
docker-compose up -d
```

## 容器说明

### 单容器模式特点

✅ **优点：**
- 一个容器提供完整服务
- 部署简单，无需管理多个容器
- 适合开发、测试环境
- 资源占用相对较少

❌ **缺点：**
- 不符合Docker最佳实践（一个容器一个进程）
- 数据库和应用耦合，难以独立扩展
- 升级维护相对复杂

### 多容器模式特点

✅ **优点：**
- 符合Docker最佳实践
- 数据库和应用独立管理
- 易于扩展和维护
- 适合生产环境

❌ **缺点：**
- 需要管理多个容器
- 配置相对复杂

## 服务验证

```bash
# 检查容器状态
docker ps

# 检查健康状态
docker inspect --format='{{.State.Health.Status}}' push-server

# 测试Web服务
curl http://localhost:8080/health

# 测试PostgreSQL
docker exec -it push-server psql -U postgres -d push_server -c "SELECT version();"
```

## 数据持久化

数据保存在Docker volume中：

```bash
# 查看volumes
docker volume ls

# 备份数据
docker exec push-server pg_dump -U postgres push_server > backup.sql

# 恢复数据
docker exec -i push-server psql -U postgres push_server < backup.sql
```

## 环境变量

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `POSTGRES_DB` | 数据库名 | `push_server` |
| `POSTGRES_USER` | 数据库用户 | `postgres` |
| `POSTGRES_PASSWORD` | 数据库密码 | `postgres123` |
| `PORT` | Web服务端口 | `8080` |
| `HUAWEI_PROJECT_ID` | 华为项目ID | 必填 |
| `HUAWEI_SERVICE_ACCOUNT_FILE` | 配置文件路径 | `/app/config/agconnect-services.json` |
| `PUSH_TOKEN_ENCRYPTION_KEY` | 加密密钥 | 必填 |

## 故障排查

### 容器无法启动

```bash
# 查看详细日志
docker logs push-server

# 检查配置文件
docker exec push-server ls -la /app/config/
```

### PostgreSQL连接失败

```bash
# 检查PostgreSQL是否运行
docker exec push-server pg_isready -U postgres

# 查看PostgreSQL日志
docker exec push-server cat /var/lib/postgresql/data/logfile
```

### 应用服务无法访问

```bash
# 检查端口映射
docker port push-server

# 检查进程
docker exec push-server ps aux
```

## 生产环境建议

1. **使用外部PostgreSQL数据库**
   - 使用云数据库服务（如AWS RDS、Azure Database）
   - 或单独部署PostgreSQL容器

2. **配置资源限制**
   ```yaml
   deploy:
     resources:
       limits:
         cpus: '2'
         memory: 2G
       reservations:
         cpus: '1'
         memory: 1G
   ```

3. **配置日志**
   ```yaml
   logging:
     driver: "json-file"
     options:
       max-size: "10m"
       max-file: "3"
   ```

4. **使用Secrets管理敏感信息**
   ```bash
   docker secret create postgres_password password.txt
   docker secret create encryption_key key.txt
   ```

5. **配置反向代理（Nginx/Caddy）**
   - 启用HTTPS
   - 负载均衡
   - 限流保护
