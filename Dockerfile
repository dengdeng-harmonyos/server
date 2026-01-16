# 构建阶段
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 安装依赖
RUN apk add --no-cache git

# 复制 go mod 文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# 运行阶段 - 包含PostgreSQL和应用服务
FROM postgres:15-alpine

# 安装必要的工具
RUN apk --no-cache add ca-certificates tzdata bash supervisor

# 创建应用目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .

# 复制配置文件和初始化脚本
COPY database/init.sql /docker-entrypoint-initdb.d/
COPY config /app/config

# 创建supervisor配置目录
RUN mkdir -p /etc/supervisor.d

# 创建supervisor配置文件
RUN cat > /etc/supervisor.d/supervisord.conf <<'EOF'
[supervisord]
nodaemon=true

[program:postgres]
command=/bin/bash -c "pg_ctl -D $PGDATA -w start && tail -f /dev/null"
autorestart=true
stdout_logfile=/dev/stdout
stdout_logfile_maxbytes=0
stderr_logfile=/dev/stderr
stderr_logfile_maxbytes=0

[program:push-server]
command=/app/main
autostart=true
autorestart=true
priority=10
startsecs=30
stdout_logfile=/dev/stdout
stdout_logfile_maxbytes=0
stderr_logfile=/dev/stderr
stderr_logfile_maxbytes=0
EOF

# 创建启动脚本
RUN cat > /start.sh <<'EOF'
#!/bin/bash
set -e

# 初始化PostgreSQL数据目录
if [ ! -s "$PGDATA/PG_VERSION" ]; then
  echo "Initializing PostgreSQL database..."
  su-exec postgres initdb
  echo "host all all 0.0.0.0/0 md5" >> $PGDATA/pg_hba.conf
  echo "listen_addresses='*'" >> $PGDATA/postgresql.conf
fi

# 启动PostgreSQL
su-exec postgres pg_ctl -D "$PGDATA" -w start

# 等待PostgreSQL启动
until su-exec postgres pg_isready; do
  echo "Waiting for PostgreSQL to be ready..."
  sleep 2
done

# 创建数据库和用户
su-exec postgres psql -c "CREATE DATABASE ${POSTGRES_DB:-push_server};" 2>/dev/null || true

# 初始化数据库表结构
if [ -f /docker-entrypoint-initdb.d/init.sql ]; then
  su-exec postgres psql -d ${POSTGRES_DB:-push_server} -f /docker-entrypoint-initdb.d/init.sql
fi

# 启动应用服务
echo "Starting push server application..."
exec /app/main
EOF

RUN chmod +x /start.sh

# 设置环境变量
ENV POSTGRES_DB=push_server \
    POSTGRES_USER=postgres \
    POSTGRES_PASSWORD=postgres123 \
    PGDATA=/var/lib/postgresql/data \
    DB_HOST=localhost \
    DB_PORT=5432 \
    PORT=8080 \
    TZ=Asia/Shanghai \
    PUSH_TOKEN_ENCRYPTION_KEY=12345678901234567890123456789012

# 暴露端口
EXPOSE 8080 5432

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
    CMD pg_isready -U postgres && wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 使用启动脚本
CMD ["/start.sh"]
