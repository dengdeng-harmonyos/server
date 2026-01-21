#!/bin/bash

# Docker 部署启动脚本
# 用于在服务器上快速部署噔噔推送服务

set -e

echo "=========================================="
echo "噔噔推送服务 Docker 部署脚本"
echo "=========================================="
echo ""

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 配置变量
DOCKER_IMAGE="ricwang/dengdeng-server:latest"
CONTAINER_NAME="push-server"
PORT="${PORT:-8080}"
SERVER_NAME="${SERVER_NAME:-噔噔推送服务}"

# 检查 Docker 是否安装
echo -e "${YELLOW}[1/6] 检查 Docker 环境...${NC}"
if ! command -v docker &> /dev/null; then
    echo -e "${RED}错误: Docker 未安装${NC}"
    echo "请先安装 Docker: curl -fsSL https://get.docker.com | bash"
    exit 1
fi

if ! command -v docker compose &> /dev/null; then
    echo -e "${RED}错误: Docker Compose 未安装${NC}"
    echo "Docker Compose 是 Docker Desktop 的一部分，或者单独安装"
    exit 1
fi

echo -e "${GREEN}✓ Docker 环境检查通过${NC}"
echo ""

# 检查或生成加密密钥
echo -e "${YELLOW}[2/6] 配置加密密钥...${NC}"
if [ -f .env ] && grep -q "PUSH_TOKEN_ENCRYPTION_KEY=" .env; then
    EXISTING_KEY=$(grep "^PUSH_TOKEN_ENCRYPTION_KEY=" .env | cut -d'=' -f2)
    if [ -n "$EXISTING_KEY" ] && [ "$EXISTING_KEY" != "your-32-character-encryption-key-here" ]; then
        echo -e "${GREEN}✓ 使用已存在的加密密钥${NC}"
    else
        echo -e "${YELLOW}⚠ 检测到默认密钥，将重新生成...${NC}"
        ENCRYPTION_KEY=$(openssl rand -base64 24)
        sed -i.bak "s|PUSH_TOKEN_ENCRYPTION_KEY=.*|PUSH_TOKEN_ENCRYPTION_KEY=$ENCRYPTION_KEY|" .env
        rm -f .env.bak
        echo -e "${GREEN}✓ 已更新加密密钥${NC}"
        echo -e "${YELLOW}新密钥: $ENCRYPTION_KEY${NC}"
    fi
else
    echo "生成新的加密密钥..."
    ENCRYPTION_KEY=$(openssl rand -base64 24)
    cat > .env <<EOF
# Push Token 加密密钥（自动生成）
PUSH_TOKEN_ENCRYPTION_KEY=$ENCRYPTION_KEY

# 服务器名称（可选）
SERVER_NAME=$SERVER_NAME
EOF
    chmod 600 .env
    echo -e "${GREEN}✓ 已创建 .env 文件${NC}"
    echo -e "${YELLOW}加密密钥: $ENCRYPTION_KEY${NC}"
    echo -e "${RED}⚠️  请妥善保存此密钥！丢失后将无法解密已存储的设备Token${NC}"
fi
echo ""

# 创建  .yml
echo -e "${YELLOW}[3/6] 创建 Docker Compose 配置...${NC}"
cat > docker-compose.yml <<EOF
services:
  push-server:
    image: $DOCKER_IMAGE
    container_name: $CONTAINER_NAME
    environment:
      - SERVER_NAME=\${SERVER_NAME}
      - PUSH_TOKEN_ENCRYPTION_KEY=\${PUSH_TOKEN_ENCRYPTION_KEY}
    ports:
      - "$PORT:8080"
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
echo -e "${GREEN}✓ Docker Compose 配置已创建${NC}"
echo ""

# 拉取最新镜像
echo -e "${YELLOW}[4/6] 拉取 Docker 镜像...${NC}"
if docker compose pull; then
    echo -e "${GREEN}✓ 镜像拉取成功${NC}"
else
    echo -e "${RED}警告: 镜像拉取失败，将尝试使用本地镜像${NC}"
fi
echo ""

# 停止旧容器（如果存在）
echo -e "${YELLOW}[5/6] 停止旧容器（如果存在）...${NC}"
if docker ps -a --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
    echo "检测到已存在的容器，正在停止..."
    docker compose down
    echo -e "${GREEN}✓ 旧容器已停止${NC}"
else
    echo -e "${GREEN}✓ 无需停止旧容器${NC}"
fi
echo ""

# 启动服务
echo -e "${YELLOW}[6/6] 启动服务...${NC}"
docker compose up -d
echo -e "${GREEN}✓ 服务已启动${NC}"
echo ""

# 等待服务就绪
echo -e "${YELLOW}等待服务启动完成...${NC}"
for i in {1..30}; do
    if docker exec $CONTAINER_NAME wget --spider -q http://localhost:8080/health 2>/dev/null; then
        echo -e "${GREEN}✓ 服务已就绪！${NC}"
        break
    fi
    if [ $i -eq 30 ]; then
        echo -e "${RED}⚠️  服务启动超时，请检查日志${NC}"
        docker compose logs --tail=50
        exit 1
    fi
    echo -n "."
    sleep 2
done
echo ""

# 显示服务信息
echo ""
echo "=========================================="
echo -e "${GREEN}✓ 部署完成！${NC}"
echo "=========================================="
echo ""
echo "服务信息："
echo "  • 容器名称: $CONTAINER_NAME"
echo "  • 访问地址: http://localhost:$PORT"
echo "  • 健康检查: http://localhost:$PORT/health"
echo "  • 镜像版本: $DOCKER_IMAGE"
echo ""
echo "常用命令："
echo "  • 查看日志:   docker compose logs -f"
echo "  • 重启服务:   docker compose restart"
echo "  • 停止服务:   docker compose down"
echo "  • 更新镜像:   docker compose pull && docker compose up -d"
echo "  • 查看状态:   docker compose ps"
echo ""
echo "数据备份："
echo "  • 备份数据库: docker exec $CONTAINER_NAME pg_dump -U postgres push_server > backup.sql"
echo "  • 恢复数据库: cat backup.sql | docker exec -i $CONTAINER_NAME psql -U postgres push_server"
echo ""
echo -e "${YELLOW}⚠️  重要提示:${NC}"
echo "  • 加密密钥已保存在 .env 文件中，请妥善保管"
echo "  • 不要将 .env 文件提交到版本控制系统"
echo "  • 定期备份 PostgreSQL 数据"
echo ""

# 显示容器状态
echo "当前容器状态："
docker compose ps
echo ""

# 测试健康检查
echo -e "${YELLOW}测试健康检查...${NC}"
if curl -f http://localhost:$PORT/health 2>/dev/null; then
    echo -e "${GREEN}✓ 健康检查通过${NC}"
else
    echo -e "${RED}⚠️  健康检查失败，请查看日志${NC}"
    echo "运行以下命令查看详细日志："
    echo "  docker compose logs -f"
fi
echo ""

echo "=========================================="
echo -e "${GREEN}部署脚本执行完毕！${NC}"
echo "=========================================="
