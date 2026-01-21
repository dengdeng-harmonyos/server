# 快速部署指南

## 🚀 一键部署

### 使用部署脚本（推荐）

```bash
# 1. 下载部署脚本
curl -O https://raw.githubusercontent.com/你的用户名/server/main/deploy.sh
chmod +x deploy.sh

# 2. 配置镜像地址（修改脚本中的镜像名称）
# 或使用环境变量
export DOCKER_IMAGE="你的用户名/dengdeng-server:latest"

# 3. 执行部署
./deploy.sh
```

脚本会自动完成：
- ✅ 检查 Docker 环境
- ✅ 生成加密密钥
- ✅ 创建配置文件
- ✅ 拉取最新镜像
- ✅ 启动服务
- ✅ 验证健康状态

---

## 📋 手动部署步骤

### 1. 创建部署目录

```bash
mkdir -p ~/dengdeng-server && cd ~/dengdeng-server
```

### 2. 创建 docker-compose.yml

```bash
cat > docker-compose.yml <<'EOF'
services:
  push-server:
    image: 你的用户名/dengdeng-server:latest
    container_name: push-server
    environment:
      - SERVER_NAME=${SERVER_NAME}
      - PUSH_TOKEN_ENCRYPTION_KEY=${PUSH_TOKEN_ENCRYPTION_KEY}
    ports:
      - "8080:8080"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: unless-stopped
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres && wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1"]
      interval: 30s
      timeout: 10s
      start_period: 40s
      retries: 3

volumes:
  postgres_data:
EOF
```

### 3. 创建 .env 文件

```bash
# 生成随机加密密钥
ENCRYPTION_KEY=$(openssl rand -base64 24)

cat > .env <<EOF
# Push Token 加密密钥
PUSH_TOKEN_ENCRYPTION_KEY=$ENCRYPTION_KEY

# 服务器名称（可选）
SERVER_NAME=噔噔推送服务
EOF

# 保护配置文件
chmod 600 .env

# 显示生成的密钥（请保存）
echo "加密密钥: $ENCRYPTION_KEY"
```

### 4. 启动服务

```bash
# 拉取镜像
docker-compose pull

# 启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f
```

### 5. 验证服务

```bash
# 检查容器状态
docker-compose ps

# 测试健康检查
curl http://localhost:8080/health

# 查看启动日志
docker-compose logs --tail=100
```

---

## 🔧 配置说明

### 环境变量

| 变量名 | 说明 | 是否必需 | 默认值 |
|--------|------|----------|--------|
| `PUSH_TOKEN_ENCRYPTION_KEY` | Token加密密钥（32字节） | ✅ **必需** | - |
| `SERVER_NAME` | 服务器显示名称 | 可选 | `噔噔推送服务` |

### 端口配置

默认使用 8080 端口，如需修改：

```yaml
ports:
  - "自定义端口:8080"
```

例如使用 9000 端口：
```yaml
ports:
  - "9000:8080"
```

---

## 📦 使用指定版本

### 方式1：使用时间戳 tag

```yaml
image: 你的用户名/dengdeng-server:20260121
```

### 方式2：使用 latest tag

```yaml
image: 你的用户名/dengdeng-server:latest
```

### 方式3：使用 SHA tag

```yaml
image: 你的用户名/dengdeng-server:sha-61dd4df
```

---

## 🔄 日常维护

### 更新服务

```bash
cd ~/dengdeng-server

# 拉取最新镜像
docker-compose pull

# 重启服务（会自动执行数据库迁移）
docker-compose up -d

# 查看更新日志
docker-compose logs -f
```

### 查看日志

```bash
# 实时日志
docker-compose logs -f

# 最近100行
docker-compose logs --tail=100

# 特定时间段
docker-compose logs --since="2026-01-21T10:00:00"
```

### 重启服务

```bash
# 优雅重启
docker-compose restart

# 完全重新创建
docker-compose down && docker-compose up -d
```

### 停止服务

```bash
# 停止但保留容器
docker-compose stop

# 停止并删除容器（数据保留）
docker-compose down

# 停止并删除所有数据（危险！）
docker-compose down -v
```

---

## 💾 数据备份与恢复

### 备份数据库

```bash
# 导出 SQL 文件
docker exec push-server pg_dump -U postgres push_server > backup_$(date +%Y%m%d).sql

# 或使用数据卷备份
docker run --rm \
  -v postgres_data:/data \
  -v $(pwd):/backup \
  alpine tar czf /backup/postgres_$(date +%Y%m%d).tar.gz -C /data .
```

### 恢复数据库

```bash
# 从 SQL 文件恢复
cat backup_20260121.sql | docker exec -i push-server psql -U postgres push_server

# 从数据卷备份恢复
docker run --rm \
  -v postgres_data:/data \
  -v $(pwd):/backup \
  alpine tar xzf /backup/postgres_20260121.tar.gz -C /data
```

### 自动备份脚本

```bash
cat > ~/dengdeng-server/backup.sh <<'EOF'
#!/bin/bash
BACKUP_DIR=~/dengdeng-backups
mkdir -p $BACKUP_DIR
docker exec push-server pg_dump -U postgres push_server > $BACKUP_DIR/backup_$(date +%Y%m%d_%H%M%S).sql
# 保留最近30天的备份
find $BACKUP_DIR -name "backup_*.sql" -mtime +30 -delete
EOF

chmod +x ~/dengdeng-server/backup.sh

# 添加到 crontab（每天凌晨2点备份）
(crontab -l 2>/dev/null; echo "0 2 * * * ~/dengdeng-server/backup.sh") | crontab -
```

---

## 🔒 安全建议

### 1. 保护配置文件

```bash
# 设置正确的权限
chmod 600 .env

# 不要提交到版本控制
echo ".env" >> .gitignore
```

### 2. 使用防火墙

```bash
# UFW 示例
sudo ufw allow from 你的IP地址 to any port 8080
sudo ufw enable

# iptables 示例
sudo iptables -A INPUT -p tcp -s 你的IP地址 --dport 8080 -j ACCEPT
sudo iptables -A INPUT -p tcp --dport 8080 -j DROP
```

### 3. 使用反向代理（推荐）

安装 Nginx：
```bash
sudo apt install nginx certbot python3-certbot-nginx
```

配置 HTTPS：
```bash
cat > /etc/nginx/sites-available/push-server <<'EOF'
server {
    listen 443 ssl http2;
    server_name your-domain.com;

    ssl_certificate /etc/letsencrypt/live/your-domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/your-domain.com/privkey.pem;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

server {
    listen 80;
    server_name your-domain.com;
    return 301 https://$server_name$request_uri;
}
EOF

sudo ln -s /etc/nginx/sites-available/push-server /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx

# 获取 SSL 证书
sudo certbot --nginx -d your-domain.com
```

### 4. 定期更新

```bash
# 设置自动更新脚本
cat > ~/dengdeng-server/auto-update.sh <<'EOF'
#!/bin/bash
cd ~/dengdeng-server
docker-compose pull
docker-compose up -d
docker image prune -f
EOF

chmod +x ~/dengdeng-server/auto-update.sh

# 每周日凌晨3点自动更新
(crontab -l 2>/dev/null; echo "0 3 * * 0 ~/dengdeng-server/auto-update.sh") | crontab -
```

---

## ❓ 故障排查

### 容器无法启动

```bash
# 查看详细日志
docker-compose logs push-server

# 检查配置
docker-compose config

# 进入容器调试
docker exec -it push-server sh
```

### 端口冲突

```bash
# 查看端口占用
sudo lsof -i :8080

# 修改端口映射
# 编辑 docker-compose.yml
ports:
  - "8081:8080"
```

### 数据库连接失败

```bash
# 检查 PostgreSQL 状态
docker exec push-server ps aux | grep postgres

# 手动连接数据库
docker exec -it push-server psql -U postgres push_server
```

### 健康检查失败

```bash
# 手动测试健康接口
curl -v http://localhost:8080/health

# 查看应用日志
docker-compose logs -f push-server

# 检查容器资源
docker stats push-server
```

---

## 📊 监控建议

### 简单监控脚本

```bash
cat > ~/dengdeng-server/monitor.sh <<'EOF'
#!/bin/bash
if ! curl -f http://localhost:8080/health > /dev/null 2>&1; then
    echo "服务异常: $(date)" >> ~/dengdeng-server/alerts.log
    # 可选：发送邮件或webhook通知
fi
EOF

chmod +x ~/dengdeng-server/monitor.sh

# 每5分钟检查一次
(crontab -l 2>/dev/null; echo "*/5 * * * * ~/dengdeng-server/monitor.sh") | crontab -
```

### 使用 Docker stats

```bash
# 实时监控资源使用
docker stats push-server

# 查看容器详情
docker inspect push-server
```

---

## 🎯 生产环境检查清单

部署前确认：

- [ ] Docker 和 Docker Compose 已安装
- [ ] 端口 8080 未被占用
- [ ] 磁盘空间充足（至少 5GB）
- [ ] `.env` 文件已创建并配置加密密钥
- [ ] 防火墙规则已配置
- [ ] SSL 证书已配置（如使用 HTTPS）
- [ ] 备份策略已制定
- [ ] 监控脚本已设置

部署后验证：

- [ ] 容器正常运行：`docker-compose ps`
- [ ] 健康检查通过：`curl http://localhost:8080/health`
- [ ] 日志无错误：`docker-compose logs`
- [ ] 数据持久化：重启后数据仍存在
- [ ] 推送功能正常：发送测试消息

---

**完成部署后，你的噔噔推送服务已准备就绪！** 🎉
