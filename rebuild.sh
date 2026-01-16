#!/bin/bash

echo "🔧 强制重新构建并部署..."
echo ""

# 停止并删除所有相关容器和镜像
echo "1️⃣ 停止旧容器..."
docker compose -f docker-compose.single.yml down

echo ""
echo "2️⃣ 删除旧镜像..."
docker rmi dangdangdang-server-push-server-all-in-one 2>/dev/null || true
docker rmi $(docker images | grep dangdangdang | awk '{print $3}') 2>/dev/null || true

echo ""
echo "3️⃣ 清理构建缓存..."
docker builder prune -f

echo ""
echo "4️⃣ 重新构建镜像（不使用缓存）..."
docker compose -f docker-compose.single.yml build --no-cache

echo ""
echo "5️⃣ 启动新容器..."
docker compose -f docker-compose.single.yml up -d

echo ""
echo "⏳ 等待服务启动（10秒）..."
sleep 10

echo ""
echo "✅ 部署完成！"
echo ""
echo "📋 查看详细日志："
echo "   docker compose -f docker-compose.single.yml logs -f"
echo ""
echo "🔍 测试服务："
echo "   curl http://localhost:8081/health"
echo ""

# 自动显示最近的日志
echo "📄 最近的日志："
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
docker compose -f docker-compose.single.yml logs --tail=50
