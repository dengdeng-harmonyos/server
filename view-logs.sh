#!/bin/bash

# 日志查看脚本

echo "==========================================="
echo "  Dangdangdang Push Server - 日志查看工具"
echo "==========================================="
echo ""

# 检查容器是否运行
if ! docker ps | grep -q "dangdangdang"; then
    echo "❌ 容器未运行，尝试启动..."
    docker compose -f docker-compose.single.yml up -d
    sleep 3
fi

echo "📋 选择查看方式:"
echo ""
echo "  1) 实时日志 (跟踪模式)"
echo "  2) 最近100行日志"
echo "  3) 只看错误日志"
echo "  4) 只看访问日志"
echo "  5) 搜索关键词"
echo "  6) 导出日志到文件"
echo ""
read -p "请选择 (1-6): " choice

case $choice in
    1)
        echo ""
        echo "🔴 实时日志 (按 Ctrl+C 退出)"
        echo "==========================================="
        docker compose -f docker-compose.single.yml logs -f
        ;;
    2)
        echo ""
        echo "📄 最近100行日志"
        echo "==========================================="
        docker compose -f docker-compose.single.yml logs --tail=100
        ;;
    3)
        echo ""
        echo "❌ 错误日志"
        echo "==========================================="
        docker compose -f docker-compose.single.yml logs | grep -i "error\|ERROR\|failed\|FAILED\|✗"
        ;;
    4)
        echo ""
        echo "🌐 访问日志"
        echo "==========================================="
        docker compose -f docker-compose.single.yml logs | grep -i "ACCESS\|→\|←"
        ;;
    5)
        echo ""
        read -p "🔍 输入搜索关键词: " keyword
        echo "搜索结果:"
        echo "==========================================="
        docker compose -f docker-compose.single.yml logs | grep -i "$keyword"
        ;;
    6)
        filename="logs_$(date +%Y%m%d_%H%M%S).txt"
        echo ""
        echo "💾 正在导出日志到: $filename"
        docker compose -f docker-compose.single.yml logs > "$filename"
        echo "✅ 导出完成: $filename"
        ;;
    *)
        echo "❌ 无效选择"
        exit 1
        ;;
esac
