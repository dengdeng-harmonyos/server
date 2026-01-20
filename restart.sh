#!/bin/bash

echo "🔄 重新构建并启动服务器..."

# 停止旧容器
docker compose down

# 重新构建并启动
docker compose up --build -d

# 等待服务启动
echo "⏳ 等待服务启动..."
sleep 3

# 显示状态
echo ""
echo "✅ 服务状态:"
docker compose ps

echo ""
echo "📋 查看日志选项:"
echo "  1. 实时日志: docker compose logs -f"
echo "  2. 使用脚本: ./view-logs.sh"
echo ""
echo "🚀 服务已启动!"
echo "   访问: http://localhost:8081/health"
