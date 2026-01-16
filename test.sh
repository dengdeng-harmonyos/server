#!/bin/bash

# 消息推送测试脚本

SERVER_URL="${SERVER_URL:-http://localhost:8080}"
PUSH_TOKEN="${PUSH_TOKEN:-test_token_12345}"

echo "=========================================="
echo "消息推送服务器测试脚本"
echo "=========================================="
echo "服务器地址: $SERVER_URL"
echo "测试 Token: $PUSH_TOKEN"
echo ""

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试函数
test_api() {
    local name=$1
    local method=$2
    local endpoint=$3
    local data=$4
    
    echo -e "${YELLOW}测试: $name${NC}"
    
    if [ "$method" = "GET" ]; then
        response=$(curl -s -w "\n%{http_code}" -X GET "$SERVER_URL$endpoint")
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" "$SERVER_URL$endpoint" \
            -H "Content-Type: application/json" \
            -d "$data")
    fi
    
    http_code=$(echo "$response" | tail -n 1)
    body=$(echo "$response" | head -n -1)
    
    if [ "$http_code" -eq 200 ] || [ "$http_code" -eq 201 ]; then
        echo -e "${GREEN}✓ 成功${NC} (HTTP $http_code)"
        echo "响应: $body"
    else
        echo -e "${RED}✗ 失败${NC} (HTTP $http_code)"
        echo "响应: $body"
    fi
    echo ""
}

# 1. 健康检查
test_api "健康检查" "GET" "/health" ""

# 2. 注册设备
test_api "注册设备" "POST" "/api/device/register" \
'{
  "push_token": "'"$PUSH_TOKEN"'",
  "device_id": "device_test_001",
  "device_type": "phone",
  "os_version": "HarmonyOS 4.0",
  "app_version": "1.0.0"
}'

# 3. 单播推送
test_api "单播推送" "POST" "/api/push/single" \
'{
  "push_token": "'"$PUSH_TOKEN"'",
  "message": {
    "title": "测试消息",
    "content": "这是一条测试推送消息",
    "data": {
      "type": "test",
      "timestamp": "'$(date +%s)'"
    }
  }
}'

# 4. 查询推送记录
test_api "查询推送记录" "GET" "/api/query/records?limit=10&offset=0" ""

# 5. 查询推送统计
test_api "查询推送统计" "GET" "/api/query/statistics" ""

echo "=========================================="
echo "测试完成"
echo "=========================================="
