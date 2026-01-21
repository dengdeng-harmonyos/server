#!/bin/bash

# 容器诊断脚本

echo "=========================================="
echo "容器诊断工具"
echo "=========================================="
echo ""

CONTAINER_NAME="push-server"

# 检查容器是否存在
if ! docker ps -a --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
    echo "❌ 容器 $CONTAINER_NAME 不存在"
    echo ""
    echo "运行中的容器："
    docker ps
    exit 1
fi

# 显示容器状态
echo "📊 容器状态："
docker ps -a --filter "name=$CONTAINER_NAME" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
echo ""

# 显示容器重启次数
echo "🔄 重启信息："
RESTARTS=$(docker inspect --format='{{.RestartCount}}' $CONTAINER_NAME 2>/dev/null)
echo "重启次数: $RESTARTS"
echo ""

# 显示完整日志（最近200行）
echo "📝 完整容器日志（最近200行）："
echo "=========================================="
docker logs --tail=200 $CONTAINER_NAME
echo "=========================================="
echo ""

# 检查容器内部进程
echo "🔍 容器内进程："
docker exec $CONTAINER_NAME ps aux 2>/dev/null || echo "无法连接到容器（可能正在重启）"
echo ""

# 检查端口监听
echo "🌐 端口监听状态："
docker exec $CONTAINER_NAME netstat -tlnp 2>/dev/null || echo "无法检查端口（可能正在重启）"
echo ""

# 检查 PostgreSQL 连接
echo "🗄️  PostgreSQL 状态："
docker exec $CONTAINER_NAME su-exec postgres pg_isready 2>/dev/null && echo "✅ PostgreSQL 正常" || echo "❌ PostgreSQL 异常"
echo ""

# 检查应用健康接口
echo "🏥 应用健康检查："
if docker exec $CONTAINER_NAME wget --spider -q http://localhost:8080/health 2>/dev/null; then
    echo "✅ 应用健康检查通过"
else
    echo "❌ 应用健康检查失败"
fi
echo ""

# 检查环境变量
echo "⚙️  环境变量："
docker exec $CONTAINER_NAME env | grep -E "PORT|GIN_MODE|SERVER_NAME|DB_|POSTGRES_|ENCRYPTION" | sort
echo ""

# 检查数据库连接
echo "🔗 测试数据库连接："
docker exec $CONTAINER_NAME psql -U postgres -d push_server -c "SELECT version();" 2>/dev/null && echo "✅ 数据库连接正常" || echo "❌ 数据库连接失败"
echo ""

# 显示容器资源使用
echo "💻 资源使用情况："
docker stats --no-stream $CONTAINER_NAME
echo ""

# 显示最近的错误日志
echo "❌ 错误日志（最近50行）："
echo "=========================================="
docker logs --tail=50 $CONTAINER_NAME 2>&1 | grep -iE "error|fatal|panic|fail" || echo "未发现错误日志"
echo "=========================================="
echo ""

# 建议
echo "💡 诊断建议："
if [ "$RESTARTS" -gt 5 ]; then
    echo "  ⚠️  容器重启次数过多 ($RESTARTS 次)"
    echo "  建议："
    echo "    1. 检查完整日志查找错误原因"
    echo "    2. 检查加密密钥是否正确配置"
    echo "    3. 检查端口是否被占用"
    echo "    4. 尝试完全停止并重新创建容器："
    echo "       docker compose down && docker compose up -d"
fi

echo ""
echo "🔧 有用的调试命令："
echo "  • 实时查看日志:    docker logs -f $CONTAINER_NAME"
echo "  • 进入容器调试:    docker exec -it $CONTAINER_NAME sh"
echo "  • 查看完整日志:    docker logs $CONTAINER_NAME"
echo "  • 重启容器:        docker compose restart"
echo "  • 完全重建:        docker compose down && docker compose up -d"
echo ""
