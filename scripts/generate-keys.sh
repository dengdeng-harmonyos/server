#!/bin/bash

# 生成配置密钥脚本

echo "========================================="
echo "配置密钥生成工具"
echo "========================================="
echo ""

# 生成32字节加密密钥
echo "1. 生成Push Token加密密钥（32字节）:"
ENCRYPTION_KEY=$(openssl rand -base64 32 | tr -d '\n')
echo "   PUSH_TOKEN_ENCRYPTION_KEY=${ENCRYPTION_KEY}"
echo ""

# 检查agconnect-services.json是否存在
if [ -f "config/agconnect-services.json" ]; then
    echo "2. ✓ AGConnect配置文件已存在: config/agconnect-services.json"
    
    # 提取项目ID
    PROJECT_ID=$(cat config/agconnect-services.json | grep -o '"project_id":"[^"]*"' | head -1 | cut -d'"' -f4)
    if [ -n "$PROJECT_ID" ]; then
        echo "   项目ID: ${PROJECT_ID}"
    fi
else
    echo "2. ✗ AGConnect配置文件不存在: config/agconnect-services.json"
    echo ""
    echo "请从AppGallery Connect下载配置文件："
    echo "   1) 登录 https://developer.huawei.com/consumer/cn/service/josp/agc/index.html"
    echo "   2) 选择项目"
    echo "   3) 项目设置 → 常规 → 我的应用"
    echo "   4) 下载 agconnect-services.json 文件"
    echo "   5) 将文件保存到 config/agconnect-services.json"
    echo "   3) 开启 Push Kit 服务"
    echo "   4) 创建服务账号并下载JSON文件"
    echo "   5) 将文件保存为 config/private.json"
fi

echo ""
echo "========================================="
echo "配置到.env文件"
echo "========================================="
echo ""

# 检查.env文件
if [ -f ".env" ]; then
    echo ".env文件已存在"
    
    # 更新加密密钥
    if grep -q "PUSH_TOKEN_ENCRYPTION_KEY=" .env; then
        sed -i.bak "s|PUSH_TOKEN_ENCRYPTION_KEY=.*|PUSH_TOKEN_ENCRYPTION_KEY=${ENCRYPTION_KEY}|" .env
        echo "✓ 已更新加密密钥"
    else
        echo "PUSH_TOKEN_ENCRYPTION_KEY=${ENCRYPTION_KEY}" >> .env
        echo "✓ 已添加加密密钥"
    fi
    
    # 更新项目ID
    if [ -n "$PROJECT_ID" ]; then
        if grep -q "HUAWEI_PROJECT_ID=" .env; then
            sed -i.bak "s|HUAWEI_PROJECT_ID=.*|HUAWEI_PROJECT_ID=${PROJECT_ID}|" .env
            echo "✓ 已更新项目ID"
        else
            echo "HUAWEI_PROJECT_ID=${PROJECT_ID}" >> .env
            echo "✓ 已添加项目ID"
        fi
    fi
    
    rm -f .env.bak
else
    echo "创建新的.env文件..."
    cp .env.example .env
    
    # 设置加密密钥
    sed -i.bak "s|PUSH_TOKEN_ENCRYPTION_KEY=.*|PUSH_TOKEN_ENCRYPTION_KEY=${ENCRYPTION_KEY}|" .env
    
    # 设置项目ID
    if [ -n "$PROJECT_ID" ]; then
        sed -i.bak "s|HUAWEI_PROJECT_ID=.*|HUAWEI_PROJECT_ID=${PROJECT_ID}|" .env
    fi
    
    rm -f .env.bak
    echo "✓ 已创建.env文件"
fi

echo ""
echo "========================================="
echo "下一步"
echo "========================================="
echo "1. 编辑 .env 文件，配置数据库连接"
echo "2. 确保 config/agconnect-services.json 文件已正确配置"
echo "3. 运行: go run cmd/server/main.go"
echo ""
